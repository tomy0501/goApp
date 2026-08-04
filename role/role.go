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
	db    *db.UserDB
	redis myRedis.RedisPool
}

func NewRoleHandler(db *db.UserDB, redis myRedis.RedisPool) RoleInterface {
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

func (r HotelRoles) CanAccessReadRole() bool {
	switch r {
	case StatusOwner, StatusManager, StatusFront, StatusStaff:
		return true
	case StatusHousekeeping:
		return false
	default:
		return false
	}
}

func (r HotelRoles) CanAccessWriteRole() bool {
	switch r {
	case StatusOwner, StatusManager, StatusStaff:
		return true
	case StatusFront, StatusHousekeeping:
		return false
	default:
		return false
	}
}

func (r HotelRoles) CanAccessConfidential() bool {
	switch r {
	case StatusOwner, StatusManager:
		return true
	case StatusFront, StatusStaff, StatusHousekeeping:
		return false
	default:
		return false
	}
}

func (r HotelRoles) CanAccessClean() bool {
	switch r {
	case StatusOwner, StatusHousekeeping:
		return true
	case StatusFront, StatusManager, StatusStaff:
		return false
	default:
		return false
	}
}
