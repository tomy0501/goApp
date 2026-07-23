package room

import (
	"goapp/role"
)

type Room struct {
	RoomID   string `json:"RoomID" binding:"required"`
	RoomName string `json:"RoomName" binding:"required"`
	RoomType string `json:"RoomType"`
}

type GetRoomInterface interface {
	GetRoomAll(facilityID string) ([]*Room, error)
	GetRoomById(facilityID string, roomID string) (*Room, error)
}

type getRoomHandler struct {
	RoomHandler GetRoomInterface
	RoleHandler role.RoleInterface
}

func NewGetRoomHandler(roomHandler GetRoomInterface, roleHandler role.RoleInterface) *getRoomHandler {
	return &getRoomHandler{
		RoomHandler: roomHandler,
		RoleHandler: roleHandler,
	}
}
