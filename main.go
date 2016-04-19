package main

import (
	"flag"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "gopkg.in/inconshreveable/log15.v2"
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
		log.Error("upgrade:", err)
		return
	}

	messageChannel := hub.AddSubscriber(
		topic,
		c,
		log.New(
			log.Ctx{
				"remote": c.RemoteAddr().String(),
			},
		),
	)

	for message := range messageChannel {
		c.WriteMessage(websocket.TextMessage, message)
	}
}

func publish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topic := vars["topic"]
	c, err := upgrader.Upgrade(w, r, nil)

	logger := log.New(
		log.Ctx{
			"remote": c.RemoteAddr().String(),
		},
	)

	if err != nil {
		logger.Error("upgrade:", err)
		return
	}

	logger.Info("New publisher connected")
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return
		}
		hub.Publish(
			topic,
			string(message),
			logger)
	}
}

func deleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topic := vars["topic"]

	hub.Delete(
		topic,
		log.New(
			log.Ctx{
				"remote": r.RemoteAddr,
			},
		),
	)
}

func main() {
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/subscribe/{topic}", subscribe)
	r.HandleFunc("/publish/{topic}", publish)
	r.HandleFunc("/topics/{topic}", deleteTopic).Methods("DELETE")

	logger := log.New(log.Ctx{"addr": *addr})
	logger.Info("Start listening")

	if *tls != "" {
		parsedTls := strings.Split(*tls, ":")
		http.ListenAndServeTLS(*addr, parsedTls[0], parsedTls[1], r)
	} else {
		http.ListenAndServe(*addr, r)
	}
}
