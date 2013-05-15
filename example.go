package main

import (
	"flag"
	"fmt"
	"github.com/trevex/golem"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func chat(conn *golem.Connection, data *golem.DataType) {
	fmt.Println(data)
	conn.Send(data)
}

func main() {

	flag.Parse()

	myrouter := golem.NewRouter()
	myrouter.On("chat", chat)

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/ws", myrouter.Handler())

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
