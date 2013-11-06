package main

import (
	"flag"
	"fmt"
	"github.com/trevex/golem"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

// Create room manager
var myroommanager = golem.NewRoomManager()

// Every join and leave request will specify the name of the room
// the connection wants to join or leave.
type RoomRequest struct {
	Name string `json:"name"`
}

// On join, join the room with the specified name.
func join(conn *golem.Connection, data *RoomRequest) {
	fmt.Println("Joining", data.Name)
	myroommanager.Join(data.Name, conn)
}

// On leave, leave the specified room.
func leave(conn *golem.Connection, data *RoomRequest) {
	fmt.Println("Leaving", data.Name)
	myroommanager.Leave(data.Name, conn)
}

// Every message has a room it broadcast to and the actual message content.
type RoomMessage struct {
	To  string `json:"to"`
	Msg string `json:"msg"`
}

// Emit the msg event to every member of the To-Room with the provided message content.
func msg(conn *golem.Connection, data *RoomMessage) {
	fmt.Println("Sending to", data.To)
	myroommanager.Emit(data.To, "msg", &data.Msg)
}

// Make sure the connection leaves all lobbies.
// When roommanager is used, this is a necessary step! Otherwise the member counting of the manager won't
// be accurate anymore. If the connection didn't join or already left all lobbies this will result in a single if
// check and therefore is not costly.
func connClose(conn *golem.Connection) {
	myroommanager.LeaveAll(conn)
	fmt.Println("Connection was closed.")
}

// On connection
func connConnect(conn *golem.Connection) {
	myroommanager.SetConnectionOptions(conn, golem.CloseConnectionOnLastRoomLeft, true)
	fmt.Println("User connected and setup.")
}

// If a room is created or removed because of insufficient users
// print the name!
// ( The functions need to be of the type func(string) and receive the rooms name as argument )
func roomCreated(name string) {
	fmt.Println("Room created:", name)
}
func roomRemoved(name string) {
	fmt.Println("Room removed:", name)
}

func main() {
	flag.Parse()

	// Create a router
	myrouter := golem.NewRouter()
	// Add the events to the router
	myrouter.On("join", join)
	myrouter.On("leave", leave)
	myrouter.On("msg", msg)
	// Connection events
	myrouter.OnClose(connClose)
	myrouter.OnConnect(connConnect)

	// React on room manager events
	myroommanager.On("create", roomCreated)
	myroommanager.On("remove", roomRemoved)

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
