package websockets

type Hub struct {
	// RoomID keys mapped to their Room references
	Rooms map[int]*Room

	// Register channel that recieves users who want to join a room
	Register chan *User

	// Unregister channel that recieves users who want to leave a room
	Unregister chan *User

	// Broadcast channel that recieves messages that should be broadcasted to all the other users in that room
	Broadcast chan *Message
}

var MainHub *Hub

func InitialiseHub() {
	MainHub = &Hub{
		Rooms: make(map[int]*Room),
		Register: make(chan *User),
		Unregister: make(chan *User),
		Broadcast: make(chan *Message, 5),
	}
}

func RunHub() {
	for {
		select {
			case message := <- MainHub.Broadcast:
				// Broadcast the messages to all users in the room if the room exists
				_, ok := MainHub.Rooms[message.RoomID]
				if ok {
					room := MainHub.Rooms[message.RoomID]
					for _, user := range room.Users {
						user.Receive <- message
					}
				}
			case user := <- MainHub.Register:
				// Check whether room exists
				_, ok := MainHub.Rooms[user.RoomID]
				if ok {
					room := MainHub.Rooms[user.RoomID]
					// Add user if he is not inside the room
					_, ok := room.Users[user.ID]
					if !ok {
						room.Users[user.ID] = user
					}
				}
			case user := <- MainHub.Unregister:
				_, ok := MainHub.Rooms[user.RoomID]
				if ok {
					room := MainHub.Rooms[user.RoomID]
					// Remove user if he is inside the room
					_, ok := room.Users[user.ID]
					if ok {
						delete(room.Users, user.ID)
						// Close the channel where this user receives messages
						close(user.Receive)
					}
				}
		}
	}
}