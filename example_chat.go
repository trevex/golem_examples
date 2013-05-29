package main

import (
	"flag"
	"fmt"
	"github.com/trevex/golem"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

// Create lobby manager
var mylobbymanager = golem.NewLobbyManager()

// Every join and leave request will specify the name of the lobby
// the connection wants to join or leave.
type LobbyRequest struct {
	Name string `json:"name"`
}

// On join, join the lobby with the specified name.
func join(conn *golem.Connection, data *LobbyRequest) {
	fmt.Println("Joining", data.Name)
	mylobbymanager.Join(data.Name, conn)
}

// On leave, leave the specified lobby.
func leave(conn *golem.Connection, data *LobbyRequest) {
	fmt.Println("Leaving", data.Name)
	mylobbymanager.Leave(data.Name, conn)
}

// Every message has a lobby it broadcast to and the actual message content.
type LobbyMessage struct {
	To  string `json:"to"`
	Msg string `json:"msg"`
}

// Emit the msg event to every member of the To-Lobby with the provided message content.
func msg(conn *golem.Connection, data *LobbyMessage) {
	fmt.Println("Sending to", data.To)
	mylobbymanager.Emit(data.To, "msg", data.Msg)
}

// Make sure the connection leaves all lobbies.
// When lobbymanager is used, this is a necessary step! Otherwise the member counting of the manager won't
// be accurate anymore. If the connection didn't join or already left all lobbies this will result in a single if
// check and therefore is not costly.
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
