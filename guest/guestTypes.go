package guest

import (
	"goapp/room"
)

type Guest struct {
	GuestID   string    `json:"GuestID" binding:"required"`
	GuestName string    `json:"GuestName" binding:"required"`
	GuestAge  *int      `json:"GuestAge" binding:"gte=0,lte=130"`
	GuestSex  *string   `json:"GuestSex" binding:"oneof=male female other unknown"`
	GuestFrom *string   `json:"GuestFrom"`
	Room      room.Room `json:"Room" binding:"required"`
	Checkin   *string   `json:"Checkin" binding:"datetime=2006-01-02"`
	Checkout  *string   `json:"Checkout" binding:"datetime=2006-01-02"`
}
