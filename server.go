package main

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var activeUsersSync = sync.Map{}
var activeSocketsSync = sync.Map{}

const VALID_NAME_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_" //valid characters for usernames and ids
const MSG_DELIMITER = "<;>"

func main() {
	rand.Seed(time.Now().UnixNano())
	connectToDB()
	dropTables()
	initDB()
	http.HandleFunc("/", servePublic)      //serve files from public folder on root
	http.HandleFunc("/ws", serveWebsocket) //serve websocket connections

	err := http.ListenAndServe(":8080", nil) //start server on port 8080
	if err != nil {
		log.Fatal(err)
	}

	log.Println("http server started on :8080")
}
func servePublic(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./public")).ServeHTTP(w, r) //serve files from public folder
}
func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // upgrade http connection to websocket
	if err == nil {                          //if upgrade was successful
		conn.WriteMessage(1, constructMessage("request_credentials", []string{})) //request credentials from client. the application is based on conn being linked to a user
		for {                                                                     //loop while connection is open.
			_, message, err := conn.ReadMessage() //read message from client
			if err == nil {                       //if no error
				messageHandler(conn, message) //handle message
			} else { //if error
				break //break loop and close connection
			}
		}
		defer func() { //deferred function to run when function ends
			logout(conn) //logout user and close connection
		}()
	}
}

func messageHandler(conn *websocket.Conn, message []byte) { //handle messages from client
	segments := strings.Split(string(message), MSG_DELIMITER) //split message into segments
	message_type := segments[0]                               //first segment is the message type and the rest are the message data
	if message_type == "credentials" {                        //if message type is credentials
		login(conn, segments) //login user
	}
	if message_type == "placeholder" {
		println("testing placeholder implementation")
	}
}

func constructMessage(message_type string, message_data []string) []byte { //construct message from message type and message data
	message := message_type + MSG_DELIMITER //create message
	for _, data := range message_data {     //add message data
		message += data + MSG_DELIMITER
	}
	return []byte(message) //return message
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
