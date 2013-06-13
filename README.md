golem examples
================================
These examples should provide insight on how to use the golem framework.

Instructions
-------------------------
After cloning the repository, simply run an example:
```
go run example_name.go
```
And access the example through your browser of choice, that supports WebSocket:
```
http://127.0.0.1:8080/example_name.html
```
To run the session example gorilla/sessions needs to be available in your GOPATH, so `go get` it first.

To run the custom protocol example v2/mgo/bson needs to be available in your GOPATH, so `go get` it first.

Note
--------------------------
Some examples require the user to open the browser's console to see the output of the example.

Tutorials
-------------------------
Most of the provided examples are part of a tutorial from golem's wiki:
* [Getting started](https://github.com/trevex/golem/wiki/Getting-started)
* [Using rooms](https://github.com/trevex/golem/wiki/Using-rooms)
* [Building a Chat application](https://github.com/trevex/golem/wiki/Building-a-chat-application)
* [Handshake authorisation using Sessions](https://github.com/trevex/golem/wiki/Handshake-authorisation-using-Sessions)
* [Using flash as WebSocket fallback](https://github.com/trevex/golem/wiki/Using-flash-as-WebSocket-fallback)
* [Custom protocol using BSON](https://github.com/trevex/golem/wiki/Custom-protocol-using-BSON)
* [Using an extended connection type](https://github.com/trevex/golem/wiki/Using-an-extended-connection-type)




example_simple.go
-------------------------
The simple example presents the most common and useful patterns provided by golem by default for simple communication.

example_data.go
-------------------------
This is a very simple example illustrating several ways to accept data. It shows how using an interface{} the raw interstage product of
a protocol can be directly forwarded to the callback function and how the protocol can be extended for certain types (i.e. strings) along with
simple communication patterns.

example_room.go
-------------------------
A single room is used to show off the general features of rooms. For a better demonstration multiple browser-tabs should be opened.

example_chat.go
-------------------------
A very basic chat application allowing the user to join or leave rooms and broadcast messages to them.

example_session.go
-------------------------
Advanced example showing how handshake verfication can be used to only allow upgrading of connections with
a valid session and authentication. The user has a login and logout button to change the session state
and see the result.

example_echo_flash.go
-------------------------
This example illustrates how to use flash sockets as a fallback for WebSockets and how to handle the flash policy file. Needs to be run as root, because flash policy socket is <= 1024.

example_protocol.go
-------------------------
The protocol example presents how to implement custom protocols. In this case a custom BSON-based protocol is used to communicate between
client and server.

example_connection_extension.go
-------------------------
Based on the simple example, but enhancing it to use a connection extension. Demonstrating additional methods and members of an extended connection.