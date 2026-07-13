package role

import (
	"goapp/db"
	"goapp/myRedis"
)

type RoleInterface interface {
	GetUserRoles(userID string) ([]Roles, error)
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
	Role    string
}
