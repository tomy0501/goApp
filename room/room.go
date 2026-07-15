package room

type Room struct {
	RoomID   string `json:"RoomID" binding:"required"`
	RoomName string `json:"RoomName" binding:"required"`
	RoomType string `json:"RoomType"`
}
