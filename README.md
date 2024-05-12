# Live CLI chat made with GO.

## Usage:

You can change the host and the port of the program by setting the --h and --p flags when building.
These are optional, the default values are "localhost" and "3090"
 
1. Build the chat.go file. This is the chat server, responsible of handling al incoming users and messages
`go build -h localhost -p 3090  chat/chat.go`
2. Build the netcat.go file. This is the client wich send the messages to the chat server.
`go build -h localhost -p 3090 netcat/netcat.go`

Done :)

Now execute first the chat.exe and then how many netcat.exe you want, so you can see the live chat working.
