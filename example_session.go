package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/trevex/golem"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

const (
	secret      = "super-secret-key"
	sessionName = "golem.sid"
)

// Create session store.
var store = sessions.NewCookieStore([]byte(secret))

// Handshake callback to validate if request has session and is logged in.
func validateSession(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, sessionName)                       // Get session.
	if v, ok := session.Values["isAuthorized"]; ok && v == true { // Check if session is authorized.
		fmt.Println("Authorized user identified!")
		return true
	} else {
		fmt.Println("Unauthorized user detected!")
		return false
	}
}

// If not available creates session and flags it as authorized.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	session.Values["isAuthorized"] = true
	session.Save(r, w)
	// Redirect back to main page to test websocket connection.
	http.Redirect(w, r, "/example_session.html", http.StatusFound)
}

// Flags session as not authorized.
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	session.Values["isAuthorized"] = false
	session.Save(r, w)
	// Redirect back to main page to test websocket connection.
	http.Redirect(w, r, "/example_session.html", http.StatusFound)
}

func main() {
	flag.Parse()

	// Create a router
	myrouter := golem.NewRouter()
	myrouter.OnHandshake(validateSession)

	// Serve the public files
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Handle login/logout
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	// Handle websockets using golems handler
	http.HandleFunc("/ws", myrouter.Handler())

	// Listen
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
