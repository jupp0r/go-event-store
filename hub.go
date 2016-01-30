package main

type Hub interface {
	AddSubscriber(topic string, connection connection) chan []byte
	RemoveSubscriber(topic string, connection connection)
	Publish(topic, message string)
	Delete(topic string)
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

func (h *pubsubHub) AddSubscriber(topic string, connection connection) chan []byte {
	t := h.fetchOrCreateTopic(topic)
	return t.AddSubscriber(connection)
}

func (h *pubsubHub) RemoveSubscriber(topic string, connection connection) {
	t := h.fetchOrCreateTopic(topic)
	t.RemoveSubscriber(connection)
}

func (h *pubsubHub) Publish(topic, message string) {
	t, ok := h.topics[topic]
	if !ok {
		t = NewTopic(NewInMemoryPersister())
		h.topics[topic] = t
	}

	t.Publish(message)
}

func (h *pubsubHub) Delete(topic string) {
	delete(h.topics, topic)
}

func (h *pubsubHub) fetchOrCreateTopic(topic string) Topic {
	_, ok := h.topics[topic]
	if !ok {
		t := NewTopic(NewInMemoryPersister())
		h.topics[topic] = t
	}

	return h.topics[topic]
}
