package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var tls = flag.String("tls", "", "tls certificate and private key, example: -tls cert.pem:key.pem")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool { return true },
}

var hub = NewHub()

type connection interface{}

func subscribe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topic := vars["topic"]
	c, err := upgrader.Upgrade(w, r, nil)
	defer c.Close()

	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	messageChannel := hub.AddSubscriber(topic, c)

	for message := range messageChannel {
		c.WriteMessage(websocket.TextMessage, message)
	}
}

func publish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topic := vars["topic"]
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return
		}
		hub.Publish(topic, string(message))
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	r := mux.NewRouter()
	r.HandleFunc("/subscribe/{topic}", subscribe)
	r.HandleFunc("/publish/{topic}", publish)

	if *tls != "" {
		parsedTls := strings.Split(*tls, ":")
		log.Fatal(http.ListenAndServeTLS(*addr, parsedTls[0], parsedTls[1], r))
	} else {
		log.Fatal(http.ListenAndServe(*addr, r))
	}
}
