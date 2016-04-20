package main

import log "gopkg.in/inconshreveable/log15.v2"

type Topic interface {
	AddSubscriber(connection, log.Logger) chan []byte
	RemoveSubscriber(connection, log.Logger)
	Publish(string, log.Logger)
}

type addMessage struct {
	Conn    connection
	Channel chan []byte
	Ready   chan struct{}
}

type publishMessage struct {
	Message string
	Ready   chan struct{}
}

type pubsubTopic struct {
	subscribers       map[connection]chan []byte
	publishChannel    chan publishMessage
	addChannel        chan addMessage
	removeChannel     chan connection
	enableLiveChannel chan addMessage
	Persister
}

func NewTopic(p Persister, log log.Logger) Topic {
	log.Info("Create topic")

	topic := &pubsubTopic{
		subscribers:       make(map[connection]chan []byte),
		publishChannel:    make(chan publishMessage),
		addChannel:        make(chan addMessage),
		removeChannel:     make(chan connection),
		enableLiveChannel: make(chan addMessage),
		Persister:         p,
	}

	go topic.run()

	return topic
}

func (t *pubsubTopic) AddSubscriber(c connection, log log.Logger) chan []byte {
	log.Info("Add subscriber")
	subscriberChannel := make(chan []byte, 1000)
	ready := make(chan struct{})

	t.addChannel <- addMessage{c, subscriberChannel, ready}

	<-ready
	return subscriberChannel
}

func (t *pubsubTopic) RemoveSubscriber(c connection, log log.Logger) {
	log.Info("Remove subscriber")
	t.removeChannel <- c
}

func (t *pubsubTopic) Publish(message string, log log.Logger) {
	log.Info("Publish", "message", message)
	ready := make(chan struct{})
	t.publishChannel <- publishMessage{message, ready}
	<-ready
}

func (t *pubsubTopic) run() {
	for {
		select {
		case m := <-t.publishChannel:
			t.publishToSubscribers(m.Message)
			t.Persist(m.Message)
			m.Ready <- struct{}{}
		case add := <-t.addChannel:
			go t.sendSubscriberHistory(add)
		case add := <-t.enableLiveChannel:
			t.enableLiveUpdates(add)
		case conn := <-t.removeChannel:
			t.removeSubscriber(conn)
		}
	}
}

func (t *pubsubTopic) sendSubscriberHistory(a addMessage) {
	a.Ready <- struct{}{}
	for _, message := range t.Persister.Read() {
		a.Channel <- []byte(message)
	}
	t.enableLiveChannel <- a
}

func (t *pubsubTopic) enableLiveUpdates(a addMessage) {
	t.subscribers[a.Conn] = a.Channel
}

func (t *pubsubTopic) removeSubscriber(conn connection) {
	delete(t.subscribers, conn)
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
