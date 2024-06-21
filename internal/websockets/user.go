package websockets

import (
	"github.com/gorilla/websocket"
)

// Messages sent in the tutorial chat room
type Message struct {
	Sender string `json:"sender"`
	SenderID int `json:"senderId"`
	Content string `json:"content"`
	RoomID int `json:"roomId"`
	UserType string `json:"userType"`
}

// User that joins the tutorial chat room
type User struct {
	ID int
	Name string
	RoomID int
	UserType string

	// Websocket Connection with the user
	Socket *websocket.Conn
	
	// Receive channel where the user receive messages from other users 
	Receive chan *Message
}

func (u *User) Read() {
	defer func() {
		// Remove user from whichever room he is in and close websocket connection to the user
		MainHub.Unregister <- u
		u.Socket.Close()
	}()

	// Keep reading messages and broadcast them to all the users in the room
	for {
		_, message, err := u.Socket.ReadMessage()
		if err != nil {
			return
		}

		msg := &Message{
			Sender: u.Name,
			SenderID: u.ID,
			Content: string(message),
			RoomID: u.RoomID,
			UserType: u.UserType,
		}

		MainHub.Broadcast <- msg
	}
}

func (u *User) Write() {
	defer u.Socket.Close()
	for message := range u.Receive {
		err := u.Socket.WriteJSON(message)
		if err != nil {
			return
		}
	}
}
