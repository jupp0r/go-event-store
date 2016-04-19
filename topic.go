package main

import log "gopkg.in/inconshreveable/log15.v2"

type Topic interface {
	AddSubscriber(connection, log.Logger) chan []byte
	RemoveSubscriber(connection, log.Logger)
	Publish(string, log.Logger)
}

type pubsubTopic struct {
	subscribers map[connection]chan []byte
	Persister
}

func NewTopic(p Persister, log log.Logger) Topic {
	log.Info("Create topic")
	return &pubsubTopic{
		make(map[connection]chan []byte),
		p,
	}
}

func (t *pubsubTopic) AddSubscriber(c connection, log log.Logger) chan []byte {
	log.Info("Add subscriber")
	subscriberChannel := make(chan []byte)

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
	t.Persister.Persist(message)
	for _, subscriber := range t.subscribers {
		go publishToSubscriber(subscriber, message)
	}
}

func publishToSubscriber(c chan []byte, message string) {
	c <- []byte(message)
}
