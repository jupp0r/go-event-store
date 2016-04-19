package main

import log "gopkg.in/inconshreveable/log15.v2"

type Hub interface {
	AddSubscriber(topic string, connection connection, logger log.Logger) chan []byte
	RemoveSubscriber(topic string, connection connection, logger log.Logger)
	Publish(topic, message string, logger log.Logger)
	Delete(topic string, logger log.Logger)
}

type pubsubHub struct {
	topics      map[string]Topic
	connections map[connection]struct{}
}

func NewHub() Hub {
	return &pubsubHub{
		make(map[string]Topic),
		make(map[connection]struct{}),
	}
}

func (h *pubsubHub) AddSubscriber(topic string, connection connection, logger log.Logger) chan []byte {
	t := h.fetchOrCreateTopic(topic, logger)
	return t.AddSubscriber(connection, logger.New(log.Ctx{"topic": topic}))
}

func (h *pubsubHub) RemoveSubscriber(topic string, connection connection, logger log.Logger) {
	t := h.fetchOrCreateTopic(topic, logger.New(log.Ctx{"topic": topic}))
	t.RemoveSubscriber(connection, logger.New(log.Ctx{"topic": topic}))
}

func (h *pubsubHub) Publish(topic, message string, logger log.Logger) {
	t, ok := h.topics[topic]
	if !ok {
		t = NewTopic(NewInMemoryPersister(), logger.New(log.Ctx{"topic": topic}))
		h.topics[topic] = t
	}

	t.Publish(message, logger.New(log.Ctx{"topic": topic}))
}

func (h *pubsubHub) Delete(topic string, logger log.Logger) {
	logger.New(log.Ctx{"topic": topic}).Info("Deleted")
	delete(h.topics, topic)
}

func (h *pubsubHub) fetchOrCreateTopic(topic string, logger log.Logger) Topic {
	_, ok := h.topics[topic]
	if !ok {
		t := NewTopic(NewInMemoryPersister(), logger.New(log.Ctx{"topic": topic}))
		h.topics[topic] = t
	}

	return h.topics[topic]
}
