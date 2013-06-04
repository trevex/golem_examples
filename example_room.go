package main

import (
	"flag"
	"fmt"
	"github.com/trevex/golem"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

// Create single room.
var myroom = golem.NewRoom()

// No input is necessary for joining or leaving.
type EmptyMessage struct{}

// Join myroom.
func join(conn *golem.Connection, data *EmptyMessage) {
	myroom.Join(conn)
	fmt.Println("Someone joined myroom.")
}

// Leave myroom.
func leave(conn *golem.Connection, data *EmptyMessage) {
	myroom.Leave(conn)
	fmt.Println("Someone left myroom.")
}

// Simple string will be received as message.
type RoomMessage struct {
	Msg string `json:"msg"`
}

// Emits the received message to all members of room.
func msg(conn *golem.Connection, data *RoomMessage) {
	myroom.Emit("msg", &data)
	fmt.Println("\"" + data.Msg + "\" sent to members of myroom.")
}

func connClose(conn *golem.Connection) {
	// Make sure to get rid of player, not necessary!
	// If room is used often, leaving on disconnects
	// can be left out, because when sending to lobbies
	// unavailable connection are automatically sorted out.
	myroom.Leave(conn)
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
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
