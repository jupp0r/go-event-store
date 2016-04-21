package main

import log "gopkg.in/inconshreveable/log15.v2"

type Hub interface {
	AddSubscriber(topic string, conn connection, logger log.Logger) chan []byte
	RemoveSubscriber(topic string, conn connection, logger log.Logger)
	Publish(topic, message string, logger log.Logger)
	Delete(topic string, logger log.Logger)
	Dump(topic string) []string
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

func (h *pubsubHub) AddSubscriber(topic string, conn connection, logger log.Logger) chan []byte {
	t := h.fetchOrCreateTopic(topic, logger)
	return t.AddSubscriber(conn, logger)
}

func (h *pubsubHub) RemoveSubscriber(topic string, conn connection, logger log.Logger) {
	t := h.fetchOrCreateTopic(topic, logger)
	t.RemoveSubscriber(conn, logger)
}

func (h *pubsubHub) Publish(topic, message string, logger log.Logger) {
	t, ok := h.topics[topic]
	if !ok {
		t = NewTopic(
			NewInMemoryPersister(),
			logger,
		)
		h.topics[topic] = t
	}

	t.Publish(message, logger)
}

func (h *pubsubHub) Delete(topic string, logger log.Logger) {
	logger.Info("Deleted")
	delete(h.topics, topic)
}

func (h *pubsubHub) fetchOrCreateTopic(topic string, logger log.Logger) Topic {
	_, ok := h.topics[topic]
	if !ok {
		t := NewTopic(NewInMemoryPersister(), logger)
		h.topics[topic] = t
	}

	return h.topics[topic]
}

func (h *pubsubHub) Dump(topic string) []string {
	t, ok := h.topics[topic]
	if !ok {
		return []string{}
	}

	res := t.Dump()
	return res
}
