package role

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

func (h *roleHandler) GetUserRoles(userID string) ([]Roles, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		slog.Error("uuid is invalid, can't get user role")
		return nil, err
	}
	roleInRedis, err := h.redis.Load(userID)
	if err != nil {
		slog.Error("can't get Roll")
		return nil, err
	}
	if len(roleInRedis) == 0 {
		slog.Info("redis role is not found check databases")

		row, err := h.db.UserDB.Query("SELECT facility_id,role FROM facility_users WHERE user_id = ?", parsedUserID[:])
		if err != nil {
			slog.Error("failed in get facility", "error", err)
			return nil, err
		}
		err = row.Err()
		if err != nil {
			slog.Error("failed when fetch facility row", "error", err)
			return nil, err
		}
		defer row.Close()
		var roles []Roles
		for row.Next() {
			var role Roles
			err = row.Scan(&role.HotelID, &role.Role)
			if err != nil {
				slog.Error("can't get roll from mysql", "db error", err)
				return nil, err
			}
			roles = append(roles, role)
		}
		roleInRedis["userRoles"] = roles
		h.redis.Save(userID, roleInRedis)
		return roles, nil
	} else {
		if rolls, ok := roleInRedis["userRoles"]; ok {
			if paserdRoll, isRole := rolls.([]Roles); isRole {
				return paserdRoll, nil
			}
		}
		errors := errors.New("Role in redis is invalid please check")
		slog.Error(errors.Error())
		return nil, errors
	}
}
