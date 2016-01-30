package main

import "fmt"

type Topic interface {
	AddSubscriber(connection) chan []byte
	RemoveSubscriber(connection)
	Publish(string)
}

type pubsubTopic struct {
	subscribers map[connection]chan []byte
	Persister
}

func NewTopic(p Persister) Topic {
	return &pubsubTopic{
		make(map[connection]chan []byte),
		p,
	}
}

func (t *pubsubTopic) AddSubscriber(c connection) chan []byte {
	subscriberChannel := make(chan []byte)

	go func() {
		for _, message := range t.Persister.Read() {
			subscriberChannel <- []byte(message)
		}
	}()

	t.subscribers[c] = subscriberChannel
	return t.subscribers[c]
}

func (t *pubsubTopic) RemoveSubscriber(c connection) {
	delete(t.subscribers, c)
}

func (t *pubsubTopic) Publish(message string) {
	fmt.Printf("publish %s\n", message)
	t.Persister.Persist(message)
	for _, subscriber := range t.subscribers {
		go publishToSubscriber(subscriber, message)
	}
}

func publishToSubscriber(c chan []byte, message string) {
	c <- []byte(message)
}
