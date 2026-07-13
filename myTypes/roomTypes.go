package myTypes

type Room struct {
	RoomID   string `json:"roomID" binding:"required"`
	RoomName string `json:"roomName" binding:"required"`
	RoomType string `json:"roomType"`
}
