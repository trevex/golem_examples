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

func json(conn *golem.Connection, data *ChatMessage) {
	fmt.Println("JSON:  ", data.Msg)
}

func raw(conn *golem.Connection, data []byte) {
	fmt.Println("Raw:   ", string(data))
}

func custom(conn *golem.Connection, data string) {
	fmt.Println("Custom:", data)
}

func customParser(data []byte) (string, bool) {
	return string(data), true
}

func main() {
	flag.Parse()

	golem.AddParser(customParser)

	myrouter := golem.NewRouter()
	myrouter.On("json", json)
	myrouter.On("raw", raw)
	myrouter.On("custom", custom)

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", myrouter.Handler())

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
