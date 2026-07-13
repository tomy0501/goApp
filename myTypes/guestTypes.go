package myTypes

import "database/sql"

type Guest struct {
	GuestID   string         `json:"guestID" binding:"required"`
	GuestName string         `json:"guestName" binding:"required"`
	GuestAge  sql.NullInt16  `json:"guestAge" binding:"gte=0,lte=130"`
	GuestSex  sql.NullString `json:"guestSex" binding:"oneof=male female other unknown"`
	GuestFrom sql.NullString `json:"guestFrom"`
	Room      Room           `json:"room" binding:"required"`
	Checkin   sql.NullString `json:"checkin" binding:"datetime=2006-01-02"`
	Checkout  sql.NullString `json:"checkout" binding:"datetime=2006-01-02"`
}
