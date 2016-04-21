package main

import (
	"strconv"
	"testing"
)

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

func TestDumping(t *testing.T) {
	p := NewInMemoryPersister()

	for k := 0; k < 100; k++ {
		testData := []string{}
		for i := 0; i < 10; i++ {
			testData = append(testData, strconv.Itoa(i))
		}

		for _, s := range testData {
			p.Persist(s)
		}

		result := p.Read()

		for i, _ := range testData {
			if testData[i] != result[i] {
				t.Fatalf("Read wrong data. Expected %s, got %s", testData[i], result[i])
			}
		}
	}
}
