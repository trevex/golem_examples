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

type RawIncomingBSONMessage struct {
	Event string   `bson:"e"`
	Data  bson.Raw `bson:"d,omitempty"`
}

type OutgoingBSONMessage struct {
	Event string      `bson:"e"`
	Data  interface{} `bson:"d,omitempty"`
}

func (_ *BSONProtocol) Unpack(data []byte) (string, []byte, error) {
	rawMsg := &RawIncomingBSONMessage{}
	err := bson.Unmarshal(data, rawMsg)
	if err != nil {
		return "", nil, err
	}
	return rawMsg.Event, rawMsg.Data.Data, nil
}

func (_ *BSONProtocol) Unmarshal(data []byte, structPtr interface{}) error {
	raw := bson.Raw{
		Kind: 3, // Embedded document
		Data: data,
	}
	return raw.Unmarshal(structPtr)
}

func (_ *BSONProtocol) MarshalAndPack(name string, structPtr interface{}) ([]byte, error) {
	outMsg := &OutgoingBSONMessage{
		Event: name,
		Data:  structPtr,
	}
	return bson.Marshal(outMsg)
}

func (_ *BSONProtocol) GetReadMode() int {
	return golem.BinaryMode
}

func (_ *BSONProtocol) GetWriteMode() int {
	return golem.BinaryMode
}

/*
 * Message types
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
	protocol := &BSONProtocol{}
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
