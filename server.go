package main

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

const MSG_DELIMITER = "<;>"
const PORT = ":8080"

var SocketToUserSyncMap = sync.Map{}
var ActiveSocketCount = 0

type ActiveUser struct {
	Id          string   //unchangable unique id for user USERID_LENGTH random characters a-z, A-Z, 0-9, _
	SocketCount int      //number of websockets user is connected with
	Websockets  sync.Map //websockets user is connected with
}

func main() {
	http.HandleFunc("/", servePublic)      //serve files from public folder on root
	http.HandleFunc("/ws", serveWebsocket) //serve websocket connections

	log.Println("http/websocket server started on " + PORT)
	err := http.ListenAndServe(PORT, nil) //start server on PORT
	if err != nil {
		log.Fatal(err)
	}
}

func servePublic(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./public")).ServeHTTP(w, r) //serve files from public folder
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	if conn, err := upgrader.Upgrade(w, r, nil); err == nil { //if upgrade was successful
		conn.WriteMessage(ConstructMessage("ping", []string{})) //request credentials from client. the application is based on conn being linked to a user
		for {                                                   //loop while connection is open.
			if _, message, err := conn.ReadMessage(); err == nil { //if no error
				messageHandler(conn, message) //handle message
			} else {
				break //break loop if error
			}
		}
		defer func() { //deferred function to run when function ends
			conn.Close() //close connection
		}()
	}
}

func messageHandler(conn *websocket.Conn, message []byte) { //handle messages from client
	segments := strings.Split(string(message), MSG_DELIMITER) //split message into segments
	switch segments[0] {                                      //switch on message type
	case "pong":
		println("pong")
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConstructMessage(message_type string, message_data []string) (int, []byte) { //construct message from message type and message data
	message := message_type             //create message
	for _, data := range message_data { //add message data
		message += MSG_DELIMITER + data
	}
	return 1, []byte(message) //return message
}
