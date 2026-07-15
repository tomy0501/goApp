package guest

import (
	"database/sql"
	"goapp/room"
)

type Guest struct {
	GuestID   string         `json:"GuestID" binding:"required"`
	GuestName string         `json:"GuestName" binding:"required"`
	GuestAge  sql.NullInt16  `json:"GuestAge" binding:"gte=0,lte=130"`
	GuestSex  sql.NullString `json:"GuestSex" binding:"oneof=male female other unknown"`
	GuestFrom sql.NullString `json:"GuestFrom"`
	Room      room.Room      `json:"Room" binding:"required"`
	Checkin   sql.NullString `json:"Checkin" binding:"datetime=2006-01-02"`
	Checkout  sql.NullString `json:"Checkout" binding:"datetime=2006-01-02"`
}
