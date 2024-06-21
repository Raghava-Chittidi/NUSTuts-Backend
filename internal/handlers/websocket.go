package handlers

import (
	"NUSTuts-Backend/internal/api"
	"NUSTuts-Backend/internal/util"
	"NUSTuts-Backend/internal/websockets"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

// Configuration to upgrade to websocket protocol
var Upgrader = websocket.Upgrader {
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Host check
		origin := r.Header.Get("origin")
		return origin == "http://localhost:5173"
	},
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	// Create and add room into main hub if room does not exist yet
	_, ok := websockets.MainHub.Rooms[tutorialId]
	if !ok {
		websockets.MainHub.Rooms[tutorialId] = &websockets.Room{
			TutorialID: tutorialId, 
			Users: make(map[int]*websockets.User),
		}
	}

	util.WriteJSON(w, api.Response{Message: "Created!"}, http.StatusOK)
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	tutorialId, err := strconv.Atoi(chi.URLParam(r, "tutorialId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	name := r.URL.Query().Get("name")
	userType := r.URL.Query().Get("userType")

	// Upgrade to websocket protocol
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		util.ErrorJSON(w, err)
		return
	}

	defer conn.Close()

	newUser := &websockets.User{
		Socket: conn,
		Receive: make(chan *websockets.Message, 10),
		ID: userId,
		Name: name,
		UserType: userType,
		RoomID: tutorialId,
	}

	// Register user who wants to join a room
	websockets.MainHub.Register <- newUser

	// Write and read message
	go newUser.Write()
	newUser.Read()
}