package main

import (
	"bytes"
	"testing"

	log "gopkg.in/inconshreveable/log15.v2"
)

func TestTopic(test *testing.T) {
	var c1, c2, c3 connection

	logger := log.New()

	t := NewTopic(NewInMemoryPersister(), logger)

	testString := []byte("foobar")

	chan1 := t.AddSubscriber(&c1, logger)
	chan2 := t.AddSubscriber(&c2, logger)
	chan3 := t.AddSubscriber(&c3, logger)

	t.Publish(string(testString), logger)

	res1 := <-chan1
	res2 := <-chan2
	res3 := <-chan3

	if !bytes.Equal(res1, testString) ||
		!bytes.Equal(res2, testString) ||
		!bytes.Equal(res3, testString) {
		test.Fatalf("expected %s, got %s, %s and %s", testString, res1, res2, res3)
	}
}
