package main

import (
	"flag"
	"fmt"
	"github.com/trevex/golem"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

// This struct represents the message which is accepted by
// the hello-function.
// If a handler function takes a special data
// type that is not an byte array, golem automatically
// tries to unmarshal the data into the specific type.
type Hello struct {
	To   string `json:"to"`
	From string `json:"from"`
}

// Type of data being emitted with answer-Event
type Answer struct {
	Msg string `json:"msg"`
}

// Function taken special data type and utilizing golem's
// inbuilt unmarshalling
func hello(conn *golem.Connection, data *Hello) {
	fmt.Println("Hello from", data.From, "to", data.To)
	conn.Emit("answer", Answer{"Thanks, client!"})
}

// Event but no data transmission
func poke(conn *golem.Connection) {
	fmt.Println("Poke-Event triggered!")
	conn.Emit("answer", Answer{"Ouch I am sensible!"})
}

func main() {
	flag.Parse()

	// Create a router
	myrouter := golem.NewRouter()
	// Add the events to the router
	myrouter.On("hello", hello)
	myrouter.On("poke", poke)

	// Serve the public files
	http.Handle("/", http.FileServer(http.Dir("./public")))
	// Handle websockets using golems handler
	http.HandleFunc("/ws", myrouter.Handler())

	// Listen
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
