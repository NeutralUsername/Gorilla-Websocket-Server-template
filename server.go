package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const MSG_DELIMITER = "<;>"                                                                //delimiter for messages
const VALID_NAME_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_" //valid characters for usernames and ids

type User struct {
	id             string //unchangable unique id for user 27 random characters a-z, A-Z, 0-9, _
	name           string //unique initially 15 random characters a-z, A-Z, 0-9, _
	password       string //hashed with SHA256
	email          string //optional
	secretQuestion string //optional
	secretAnswer   string //optional

	power      int       //power of user
	lastActive time.Time //timestamp when last active websocket was closed
	rating     float64   //intial value is 1200

	websockets map[*websocket.Conn]bool //map of websockets linked to user
}

var userids = make(map[string]*User)             //map-- userid : user
var usernames = make(map[string]*User)           //map-- username : user
var websockets = make(map[*websocket.Conn]*User) //map-- websocket : user

func main() {
	rand.Seed(time.Now().UnixNano())

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
		for {                                                                     //loop while connection is open
			_, message, err := conn.ReadMessage() //read message from client
			if err == nil {                       //if no error
				messageHandler(conn, message) //handle message
			} else { //if error
				log.Println(err) //log error
				break            //break loop and close connection
			}
		}
		defer func() { //deferred function to run when function ends
			logout(conn) //logout user and close connection
		}()
	}
}

func messageHandler(conn *websocket.Conn, message []byte) { //handle messages from client
	segments := strings.Split(string(message), MSG_DELIMITER) //split message into segments
	fmt.Println(segments)
	message_type := segments[0]        //first segment is the message type and the rest are the message data
	if message_type == "credentials" { //if message type is credentials
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

func login(conn *websocket.Conn, segments []string) {
	var user *User
	if len(segments) == 4 && validateCredentials(segments[1], segments[2], segments[3]) { //if message data count is valid
		user = userids[segments[1]] //set user pointer to user with userid
	} else {
		user = createNewUser() //else create new user
	}
	websockets[conn] = user      //add websocket to user's websockets
	user.websockets[conn] = true //add websocket to user's websockets
	conn.WriteMessage(1, constructMessage("user_data", []string{user.id, user.name, user.password, user.email, user.secretQuestion, user.secretAnswer, strconv.Itoa(user.power), strconv.FormatFloat(user.rating, 'f', 6, 64)}))
}
func logout(conn *websocket.Conn) {
	user := websockets[conn]
	delete(websockets, conn)       //delete websocket from websockets map
	delete(user.websockets, conn)  //delete websocket from user's websockets map
	if len(user.websockets) == 0 { //if user has no more websockets
		user.lastActive = time.Now() //update last active time
	}
	conn.Close() //close connection
}

func createNewUser() *User {
	username := genRandString(15) //generate random username
	for {
		if _, ok := usernames[username]; !ok { //if username is not taken
			break //break loop and continue
		} else {
			username = genRandString(15) //generate random username
		}
	}
	password := SHA256(genRandString(15)) //hash password with SHA256
	userid := genRandString(27)           //generate random userid
	for {
		if _, ok := userids[userid]; !ok { //if userid is not taken
			break //break loop and continue
		} else {
			userid = genRandString(27) //generate random userid
		}
	}
	user := User{ //create user with random credentials and default values
		id:             userid,
		name:           username,
		password:       password,
		email:          "",
		secretQuestion: "",
		secretAnswer:   "",

		power:      0,
		lastActive: time.Now(),
		rating:     1200,

		websockets: make(map[*websocket.Conn]bool),
	}
	userids[user.id] = &user
	usernames[user.name] = &user
	return &user
}
func validateCredentials(userid string, username string, password string) bool {
	if user, ok := userids[userid]; ok { //if userid belongs to a user
		if user.name == username && user.password == password { //if username and password match
			return true
		}
	}
	return false
}
func genRandString(length int) string {
	str := ""
	for i := 0; i < length; i++ {
		str += string(VALID_NAME_CHARS[rand.Intn(len(VALID_NAME_CHARS))])
	}
	return str
}
func SHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
