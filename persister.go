package main

type Persister interface {
	Persist(message string)
	Read() []string
}

type inMemoryPersister struct {
	messages []string
}

func NewInMemoryPersister() Persister {
	return &inMemoryPersister{
		make([]string, 1000),
	}
}

func (p *inMemoryPersister) Persist(message string) {
	p.messages = append(p.messages, message)
}

func (p *inMemoryPersister) Read() []string {
	return p.messages
}
