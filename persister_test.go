package main

import "testing"

func TestInMemoryPersister(t *testing.T) {
	p := NewInMemoryPersister()

	testData := []string{
		"foo",
		"bar",
		"baz",
	}

	for _, s := range testData {
		p.Persist(s)
	}

	result := p.Read()

	for i, _ := range result {
		if testData[i] != result[i] {
			t.Fatalf("Read wrong data. Expected %s, got %s", testData[i], result[i])
		}
	}
}
