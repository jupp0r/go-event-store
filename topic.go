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
}

type publishMessage struct {
	Message string
}

type pubsubTopic struct {
	subscribers    map[connection]chan []byte
	publishChannel chan publishMessage
	addChannel     chan addMessage
	removeChannel  chan connection
	Persister
}

func NewTopic(p Persister, log log.Logger) Topic {
	log.Info("Create topic")

	topic := &pubsubTopic{
		subscribers:    make(map[connection]chan []byte),
		publishChannel: make(chan publishMessage),
		addChannel:     make(chan addMessage),
		removeChannel:  make(chan connection),
		Persister:      p,
	}

	go topic.run()

	return topic
}

func (t *pubsubTopic) AddSubscriber(c connection, log log.Logger) chan []byte {
	log.Info("Add subscriber")
	subscriberChannel := make(chan []byte, 10000)

	t.addChannel <- addMessage{c, subscriberChannel}

	return subscriberChannel
}

func (t *pubsubTopic) RemoveSubscriber(c connection, log log.Logger) {
	log.Info("Remove subscriber")
	t.removeChannel <- c
}

func (t *pubsubTopic) Publish(message string, log log.Logger) {
	log.Info("Publish", "message", message)
	t.publishChannel <- publishMessage{message}
}

func (t *pubsubTopic) run() {
	for {
		select {
		case m := <-t.publishChannel:
			t.publishToSubscribers(m.Message)
			t.Persist(m.Message)
		case add := <-t.addChannel:
			t.addSubscriber(add)
		case conn := <-t.removeChannel:
			t.removeSubscriber(conn)
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
	return t.Read()
}
