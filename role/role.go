package role

import (
	"goapp/db"
	"goapp/myRedis"
	"slices"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
type RoleInterface interface {
	GetUserRoles(userID string) ([]Roles, error)
	GetUserRoleByFacility(userID string, facilityID string) (*Roles, error)
}

type roleHandler struct {
	db    db.UserDB
	redis myRedis.RedisPool
}

func NewRoleHandler(db db.UserDB, redis myRedis.RedisPool) RoleInterface {
	return &roleHandler{
		db:    db,
		redis: redis,
	}
}

type Roles struct {
	HotelID string
	Role    HotelRoles
}

func GetRoles(r []Roles, hotelID string) *Roles {
	sliceRole := slices.DeleteFunc(r, func(r Roles) bool {
		return r.HotelID != hotelID
	})
	if len(sliceRole) == 0 {
		return nil
	}
	return &sliceRole[0]
}

type HotelRoles string

const (
	StatusOwner        HotelRoles = "owner"
	StatusManager      HotelRoles = "manager"
	StatusFront        HotelRoles = "front"
	StatusHousekeeping HotelRoles = "housekeeping"
	StatusStaff        HotelRoles = "staff"
)
