package main

import (
	"flag"
	"fmt"
	"github.com/trevex/golem"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

var mylobbymanager = golem.NewLobbyManager()

type EmptyMessage struct{}

func join(conn *golem.Connection, data *EmptyMessage) {
	mylobbymanager.Join("test", conn)
}

func leave(conn *golem.Connection, data *EmptyMessage) {
	mylobbymanager.Leave("test", conn)
}

type LobbyMessage struct {
	Msg string `json:"msg"`
}

func lobby(conn *golem.Connection, data *LobbyMessage) {
	mylobbymanager.Send("test", []byte(data.Msg))
}

func connClose(conn *golem.Connection) {
	mylobbymanager.LeaveAll(conn)
}

func main() {
	flag.Parse()

	// Create a router
	myrouter := golem.NewRouter()
	// Add the events to the router
	myrouter.On("join", join)
	myrouter.On("leave", leave)
	myrouter.On("lobby", lobby)
	myrouter.OnClose(connClose)

	// Serve the public files
	http.Handle("/", http.FileServer(http.Dir("./public")))
	// Handle websockets using golems handler
	http.HandleFunc("/ws", myrouter.Handler())

	// Listen
	fmt.Println("Listening on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
