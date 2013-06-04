package main

import (
	"flag"
	"github.com/trevex/golem"
	"log"
	"net"
	"net/http"
	"os"
)

var addr = flag.String("addr", ":8080", "http service address")

type EchoMessage struct {
	Msg string `json:"msg"`
}

func echo(conn *golem.Connection, data *EchoMessage) {
	log.Print("Echo message received.")
	conn.Emit("echo", &data.Msg)
}

func loadPolicyFile(filename string) ([]byte, error) {
	if file, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		defer file.Close()
		stat, _ := file.Stat()
		data := make([]byte, stat.Size()+1)
		file.Read(data)
		data[stat.Size()] = 0
		return data, nil
	}
}

func servePolicy(policy []byte, port string) {
	if addr, err := net.ResolveTCPAddr("tcp4", port); err != nil {
		log.Fatal("ResolveTCPAddr: ", err)
	} else {
		log.Print("Serving on flash policy at: ", addr.String())
		if l, err := net.ListenTCP("tcp4", addr); err != nil {
			log.Fatal("ListenTCP:", err)
			return
		} else {
			for {
				if s, err := l.AcceptTCP(); err != nil {
					log.Print("Accept: ", err)
				} else {
					go func(s *net.TCPConn) {
						log.Print("Policy requested.")
						s.Write(policy)
						s.Close()
					}(s)
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	// Create a router
	myrouter := golem.NewRouter()
	// Add the events to the router
	myrouter.On("echo", echo)

	// Serve the public files
	http.Handle("/", http.FileServer(http.Dir("./public")))
	// Handle websockets using golems handler
	http.HandleFunc("/ws", myrouter.Handler())

	if policy, err := loadPolicyFile("example_flash_policy.xml"); err != nil {
		log.Fatal("Flash policy: ", err)
	} else {
		go servePolicy(policy, ":843")
	}

	// Listen
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
