package main

import (
	"flag"
	"fmt"
	"github.com/trevex/golem"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

/*
 * Custom Protocol
 */

type BSONProtocol struct{}

// RawIncomingBSONMessage represents the first step of unmarshalling of incoming data
// and partially unmarshals the incoming data to get the event name.
type RawIncomingBSONMessage struct {
	Event string   `bson:"e"`
	Data  bson.Raw `bson:"d,omitempty"`
}

// General structure to represent outgoing data.
type OutgoingBSONMessage struct {
	Event string      `bson:"e"`
	Data  interface{} `bson:"d,omitempty"`
}

// Partially unmarshal incoming data to unpack event name.
func (_ *BSONProtocol) Unpack(data []byte) (string, []byte, error) {
	rawMsg := &RawIncomingBSONMessage{}
	err := bson.Unmarshal(data, rawMsg)
	if err != nil {
		return "", nil, err
	}
	return rawMsg.Event, rawMsg.Data.Data, nil
}

// Unmarshal the leftover data into the desired type of the callback.
func (_ *BSONProtocol) Unmarshal(data []byte, structPtr interface{}) error {
	raw := bson.Raw{
		Kind: 3, // Embedded document
		Data: data,
	}
	return raw.Unmarshal(structPtr)
}

// Marshal and pack data into array of bytes for sending.
func (_ *BSONProtocol) MarshalAndPack(name string, structPtr interface{}) ([]byte, error) {
	outMsg := &OutgoingBSONMessage{
		Event: name,
		Data:  structPtr,
	}
	return bson.Marshal(outMsg)
}

// Read mode should be binary for BSON
func (_ *BSONProtocol) GetReadMode() int {
	return golem.BinaryMode
}

// Write mode should be binary for BSON
func (_ *BSONProtocol) GetWriteMode() int {
	return golem.BinaryMode
}

/*
 * Message types
 * See example_simple for more documentation.
 */

type Hello struct {
	To   string `bson:"to"`
	From string `bson:"from"`
}

type Answer struct {
	Msg string `bson:"msg"`
}

/*
 * Event handlers
 * See example_simple for more documentation on these handlers.
 * This program is the same except the different protocol in use.
 */

func hello(conn *golem.Connection, data *Hello) {
	fmt.Println("Hello from", data.From, "to", data.To)
	conn.Emit("answer", &Answer{"Thanks, client!"})
}

func poke(conn *golem.Connection) {
	fmt.Println("Poke-Event triggered!")
	conn.Emit("answer", &Answer{"Ouch I am sensible!"})
}

func main() {
	flag.Parse()

	// Create a router
	myrouter := golem.NewRouter()

	// Create protocol instance representing the interface
	protocol := &BSONProtocol{}
	// Use the new protocol, if all future handlers should use the protocol
	// golem.SetInitialProtocol can be used.
	myrouter.SetProtocol(protocol)

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
