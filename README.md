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
To run the session example gorilla/sessions needs to be available in your GOPATH, so go get it first.

Note
--------------------------
Some examples require the user to open the browser's console to see the output of the example.





example_data.go
-------------------------
This basic example shows several ways to interact with incoming data.

example_room.go
-------------------------
Another basic example illustrating the use of a single room instance.

example_chat.go
-------------------------
Simple example illustrating how to use the room manager. After accessing the page the client can
join chat rooms, leave them and send messages to rooms.

example_session.go
-------------------------
Advanced example showing how handshake verfication can be used to only allow upgrade connection with
a valid session and authentication. The user has a login and logout button to change the session state
and see the result.