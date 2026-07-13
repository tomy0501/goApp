package logic

import (
	"database/sql"
	"goapp/myTypes"
)

type AuthHandlers struct {
	AuthStorage ApplyLogic
}

func NewUserDB(db *sql.DB) UserDB {
	return UserDB{
		userDB: db,
	}
}

type UserDB struct {
	userDB *sql.DB
}
type ApplyLogic interface {
	myTypes.FacilityInterface
	myTypes.ReservationInterface
}
