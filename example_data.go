package main

import (
	"flag"
	"fmt"
	"github.com/trevex/golem"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

// This struct represent the message the json function
// should receive. If a function takes a special data
// type that is not an byte array, golem automatically
// tries to unmarshal the data in to this specific type.
// Since the json package is used tags work as well.
type ChatMessage struct {
	Msg string `json:"msg"`
}

// Function taken special data type and utilizing golem's
// inbuilt unmarshalling
func json(conn *golem.Connection, data *ChatMessage) {
	fmt.Println("JSON:  ", data.Msg)
	conn.Send <- []byte("Json received.")
}

// If a function accepts a byte array the data is directly
// forwarded to the function without any parsing involved.
// Hence it is the fastest way.
func raw(conn *golem.Connection, data []byte) {
	fmt.Println("Raw:   ", string(data))
	conn.Send <- []byte("Raw byte array received.")
}

// If a parser is known for the specific data type it is
// automatically used.
func custom(conn *golem.Connection, data string) {
	fmt.Println("Custom:", data)
	conn.Send <- []byte("Custom handler use to receive data.")
}

// Custom parsers take a byte array as argument and return
// the data type they parse to and a boolean (to validate if
// parsing was successful).
func customParser(data []byte) (string, bool) {
	return string(data), true
}

func main() {
	flag.Parse()

	// Add the custom parser that returns strings
	golem.AddParser(customParser)

	// Create a router
	myrouter := golem.NewRouter()
	// Add the events to the router
	myrouter.On("json", json)
	myrouter.On("raw", raw)
	myrouter.On("custom", custom)

	//
	myrouter.OnClose(func(conn *golem.Connection) {
		fmt.Println("Client disconnected!")
	})

	// Serve the public files
	http.Handle("/", http.FileServer(http.Dir("./public")))
	// Handle websockets using golems handler
	http.HandleFunc("/ws", myrouter.Handler())

	// Listen
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
