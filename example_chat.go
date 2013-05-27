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

type LobbyRequest struct {
	Name string `json:"name"`
}

func join(conn *golem.Connection, data *LobbyRequest) {
	fmt.Println("Joining", data.Name)
	mylobbymanager.Join(data.Name, conn)
}

func leave(conn *golem.Connection, data *LobbyRequest) {
	fmt.Println("Leaving", data.Name)
	mylobbymanager.Leave(data.Name, conn)
}

type LobbyMessage struct {
	To  string `json:"to"`
	Msg string `json:"msg"`
}

func msg(conn *golem.Connection, data *LobbyMessage) {
	fmt.Println("Sending to", data.To)
	mylobbymanager.Emit(data.To, "msg", data.Msg)
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
	myrouter.On("msg", msg)
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
