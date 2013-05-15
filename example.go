package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"github.com/trevex/golem"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	go golem.StartHub()
	http.Handle("/ws", websocket.Handler(golem.WebSocketHandler))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
