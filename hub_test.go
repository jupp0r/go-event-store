package main

import (
	"fmt"
	"testing"

	log "gopkg.in/inconshreveable/log15.v2"
)

func TestPubSubHub(t *testing.T) {
	// setup
	h := NewHub()
	logger := log.New()
	nsubscribers := 20
	nmessages := 20
	messages := make([]string, nmessages)
	subscriberConnections := make([]int, nsubscribers)
	subscriberChannels := make([]chan []byte, nsubscribers)

	topic := "foo"

	for i, _ := range subscriberConnections {
		conn := connection(fmt.Sprintf("%d", i))
		subscriberChannels[i] = h.AddSubscriber(
			topic,
			conn,
			logger.New(
				log.Ctx{
					"connection": conn,
				},
			),
		)
	}

	for i, _ := range messages {
		messages[i] = fmt.Sprintf("message %d", i)
	}

	// test
	for _, m := range messages {
		h.Publish(topic, m, logger)
	}

	// verify
	for i, c := range subscriberChannels {
		for _, m := range messages {
			s := string(<-c)
			if m != s {
				t.Fatalf("Expected %s, got %s on connection %d", m, s, i)
			}
		}
	}
}
