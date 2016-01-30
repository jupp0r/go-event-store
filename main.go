package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var hub = NewHub()

type connection interface{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// topic := vars["topic"]
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

}

func publish(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// topic := vars["topic"]
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

}

func main() {
	flag.Parse()
	log.SetFlags(0)
	r := mux.NewRouter()
	r.HandleFunc("/echo", echo)
	r.HandleFunc("/subscribe/{topic}", subscribe)
	r.HandleFunc("/publish/{topic}", publish)
	log.Fatal(http.ListenAndServe(*addr, r))
}
