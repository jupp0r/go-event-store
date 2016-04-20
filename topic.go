package main

import log "gopkg.in/inconshreveable/log15.v2"

type Topic interface {
	AddSubscriber(connection, log.Logger) chan []byte
	RemoveSubscriber(connection, log.Logger)
	Publish(string, log.Logger)
}

type pubsubTopic struct {
	subscribers       map[connection]chan []byte
	publishingChannel chan string
	Persister
}

func NewTopic(p Persister, log log.Logger) Topic {
	log.Info("Create topic")

	topic := &pubsubTopic{
		make(map[connection]chan []byte),
		make(chan string),
		p,
	}

	go topic.distribute()

	return topic
}

func (t *pubsubTopic) AddSubscriber(c connection, log log.Logger) chan []byte {
	log.Info("Add subscriber")
	subscriberChannel := make(chan []byte, 1000)

	go func() {
		for _, message := range t.Persister.Read() {
			subscriberChannel <- []byte(message)
		}
	}()

	t.subscribers[c] = subscriberChannel
	return t.subscribers[c]
}

func (t *pubsubTopic) RemoveSubscriber(c connection, log log.Logger) {
	log.Info("Remove subscriber")
	delete(t.subscribers, c)
}

func (t *pubsubTopic) Publish(message string, log log.Logger) {
	log.Info("Publish", "message", message)
	t.publishingChannel <- message
}

func (t *pubsubTopic) distribute() {
	for message := range t.publishingChannel {
		t.publishToSubscribers(message)
	}
}

func (t *pubsubTopic) publishToSubscribers(message string) {
	for c, subscriber := range t.subscribers {
		log.Info("Publishing", "message", message, "connection", c)
		publishToSubscriber(message, subscriber)
	}
}

func publishToSubscriber(message string, c chan []byte) {
	c <- []byte(message)
}
