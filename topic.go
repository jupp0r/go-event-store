package main

import log "gopkg.in/inconshreveable/log15.v2"

type Topic interface {
	AddSubscriber(connection, log.Logger) chan []byte
	RemoveSubscriber(connection, log.Logger)
	Publish(string, log.Logger)
	Dump() []string
}

type addMessage struct {
	Conn    connection
	Channel chan []byte
	ready   chan struct{}
}

type removeMessage struct {
	Conn  connection
	ready chan struct{}
}

type publishMessage struct {
	Message string
	ready   chan struct{}
}

type pubsubTopic struct {
	subscribers    map[connection]chan []byte
	publishChannel chan publishMessage
	addChannel     chan addMessage
	removeChannel  chan removeMessage
	readChannel    chan chan []string
	Persister
}

func NewTopic(p Persister, log log.Logger) Topic {
	log.Info("Create topic")

	topic := &pubsubTopic{
		subscribers:    make(map[connection]chan []byte),
		publishChannel: make(chan publishMessage),
		addChannel:     make(chan addMessage),
		removeChannel:  make(chan removeMessage),
		readChannel:    make(chan chan []string),
		Persister:      p,
	}

	go topic.run()

	return topic
}

func (t *pubsubTopic) AddSubscriber(c connection, log log.Logger) chan []byte {
	log.Info("Add subscriber")
	subscriberChannel := make(chan []byte, 10000)
	ready := make(chan struct{})
	t.addChannel <- addMessage{c, subscriberChannel, ready}
	<-ready
	return subscriberChannel
}

func (t *pubsubTopic) RemoveSubscriber(c connection, log log.Logger) {
	log.Info("Remove subscriber")
	ready := make(chan struct{})
	t.removeChannel <- removeMessage{c, ready}
	<-ready
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
			m.ready <- struct{}{}
		case add := <-t.addChannel:
			t.addSubscriber(add)
			add.ready <- struct{}{}
		case remove := <-t.removeChannel:
			t.removeSubscriber(remove.Conn)
			remove.ready <- struct{}{}
		case read := <-t.readChannel:
			read <- t.Read()
		}
	}
}

func (t *pubsubTopic) addSubscriber(a addMessage) {
	streamInput := make(chan []byte)
	t.subscribers[a.Conn] = streamInput
	go RunStreamer(t.Persister.Read(), streamInput, a.Channel)
}

func (t *pubsubTopic) removeSubscriber(conn connection) {
	close(t.subscribers[conn])
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

func (t *pubsubTopic) Dump() []string {
	result := make(chan []string)
	t.readChannel <- result
	res := <-result
	return res
}
