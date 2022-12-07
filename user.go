package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	id          string   //unchangable unique id for user 27 random characters a-z, A-Z, 0-9, _
	websockets  sync.Map //websockets user is connected with
	socketCount int      //number of websockets user is connected with
}

func login(conn *websocket.Conn, segments []string) {
	var user *User
	if len(segments) == 4 && validateUserCredentials(segments[1], segments[2], segments[3]) { //if message data count is valid and credentials are valid
		if userRef, ok := activeUsersSync.Load(segments[1]); ok { //if user is already logged in
			user = userRef.(*User)
		} else {
			user = &User{ //create new user
				id:         segments[1],
				websockets: sync.Map{},
			}
			activeUsersSync.Store(segments[1], user) //add user to activeUsers map
		}
	} else {
		user = createNewUser()               //create new user
		activeUsersSync.Store(user.id, user) //add user to activeUsers map
	}
	activeSocketsSync.Store(conn, user) //add websocket to user's websockets
	user.websockets.Store(conn, true)   //add user to websocket's users
	user.socketCount++                  //increment socket count
	userData := selectUser(user.id, false)
	jso, _ := json.Marshal(userData)
	conn.WriteMessage(1, []byte("user_data"+MSG_DELIMITER+string(jso))) //send user data to client
}

func createNewUser() *User {
	username := genRandString(20) //generate random username
	for {
		if !usernameInUse(username) { //if username is not taken
			break //break loop
		}
		username = genRandString(20) //generate new username
	}
	password := SHA256(genRandString(20)) //hash password with SHA256
	userid := genRandString(27)           //generate random userid
	for {
		if !useridInUse(userid) { //if userid is not taken
			break //break loop
		}
		userid = genRandString(27) //generate new userid
	}
	insertUser(userid, username, password, 0, time.Now(), time.Now()) //insert user into database
	user := User{                                                     //create user with random credentials and default values
		id:         userid,
		websockets: sync.Map{},
	}
	return &user
}

func logout(conn *websocket.Conn) {
	userRef, ok := activeSocketsSync.Load(conn)
	if ok {
		user := userRef.(*User)
		activeSocketsSync.Delete(conn) //delete websocket from activeSockets map
		user.websockets.Delete(conn)   //delete websocket from user's websockets
		user.socketCount--             //decrement socket count
		if user.socketCount == 0 {     //if user has no more websockets
			activeUsersSync.Delete(user.id)                                                                                                    //delete user from activeUsers map
			changeProperty(user, "lastActiveAt", time.Now().Format(time.RFC3339), selectUser(user.id, true).LastActiveAt.Format(time.RFC3339)) //update lastActiveAt
		}
	}
	conn.Close() //close connection
}

func changeProperty(user *User, property string, newValue string, oldValue string) bool {
	_, err := db.Exec("UPDATE users SET "+property+" = ? WHERE id = ?", newValue, user.id) //update property in database
	if err == nil {                                                                        //if property was updated successfully
		_, err := db.Exec("INSERT INTO userChanges (userId, columnN, oldValue, newValue, createdAt) VALUES (?, ?, ?, ?, ?)", user.id, property, oldValue, newValue, time.Now().Format(time.RFC3339)) //insert change into database
		return err == nil
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
