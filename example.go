package main

import (
	"flag"
	"fmt"
	"github.com/trevex/golem"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

type ChatMessage struct {
	Msg string `json:"msg"`
}

func chat(conn *golem.Connection, data *ChatMessage) {
	fmt.Println(data.Msg)
}

func main() {

	flag.Parse()

	myrouter := golem.NewRouter()
	myrouter.On("chat", chat)

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", myrouter.Handler())

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
