# Go Event Store [![Build Status](https://travis-ci.org/jupp0r/go-event-store.svg?branch=master)](https://travis-ci.org/jupp0r/go-event-store) [![Coverage Status](https://coveralls.io/repos/github/jupp0r/go-event-store/badge.svg?branch=master)](https://coveralls.io/github/jupp0r/go-event-store?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/jupp0r/go-event-store)](https://goreportcard.com/report/github.com/jupp0r/go-event-store)

This service implements a persistent event store with exchangable storage backends and pubsub semantics over websocket.
It's currently in alpha state. The API will change quite frequently.

## Installing

    go get github.com/jupp0r/go-event-store

## Running

    $GOPATH/bin/go-event-store 0.0.0.0:8080
