# Go Event Store [![Build Status](https://travis-ci.org/jupp0r/go-event-store.svg?branch=master)](https://travis-ci.org/jupp0r/go-event-store)

This service implements a persistent event store with exchangable storage backends and pubsub semantics over websocket.
It's currently in alpha state. The API will change quite frequently.

## Installing

    go get github.com/jupp0r/go-event-store

## Running

    $GOPATH/bin/go-event-store 0.0.0.0:8080
