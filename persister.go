package main

import log "gopkg.in/inconshreveable/log15.v2"

type Persister interface {
	Persist(message string)
	Read() []string
	Close()
}

type inMemoryPersister struct {
	messages     []string
	writeChannel chan string
	readChannel  chan chan []string
	closeChannel chan struct{}
}

func NewInMemoryPersister() Persister {
	p := &inMemoryPersister{
		messages:     make([]string, 0),
		writeChannel: make(chan string),
		readChannel:  make(chan chan []string),
		closeChannel: make(chan struct{}),
	}

	go p.run()

	return p
}

func (p *inMemoryPersister) run() {
	for {
		select {
		case message := <-p.writeChannel:
			log.Info("persisting", "message", message)
			p.messages = append(p.messages, message)
		case read := <-p.readChannel:
			log.Info("writing persisted messages", "num", len(p.messages))
			read <- p.messages
		case <-p.closeChannel:
			return
		}
	}
}

func (p *inMemoryPersister) Persist(message string) {
	p.writeChannel <- message
}

func (p *inMemoryPersister) Read() []string {
	returnChannel := make(chan []string)
	p.readChannel <- returnChannel
	return <-returnChannel
}

func (p *inMemoryPersister) Close() {
	p.closeChannel <- struct{}{}
}
