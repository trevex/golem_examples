package main

import (
	"flag"
	"fmt"
	"github.com/trevex/golem"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

var mylobby = golem.NewLobby()

type EmptyMessage struct{}

func join(conn *golem.Connection, data *EmptyMessage) {
	mylobby.Join(conn)
	fmt.Println("Someone joined channel.")
}

func leave(conn *golem.Connection, data *EmptyMessage) {
	mylobby.Leave(conn)
	fmt.Println("Someone left channel.")
}

type LobbyMessage struct {
	Msg string `json:"msg"`
}

func lobby(conn *golem.Connection, data *LobbyMessage) {
	mylobby.Send([]byte("lobbyMessage { \"msg\": \"" + data.Msg + "\" }"))
	fmt.Println("\"", data.Msg, "\" sent to members of channel.")
}

func connClose(conn *golem.Connection) {
	// Make sure to get rid of player, not necessary!
	// If lobby is used often, leaving on disconnects
	// can be left out, because when sending to lobbies
	// unavailable connection are automatically sorted out.
	mylobby.Leave(conn)
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
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
