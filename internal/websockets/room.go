package websockets

// Tutorial chat room
type Room struct {
	TutorialID int
	// UserID keys mapped to their User references
	Users map[int]*User
}