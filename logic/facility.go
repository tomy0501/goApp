package logic

import (
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
)

type FacilityInterface interface {
	GetFacilityInfo(userID string) (*[]FacilityJSON, error)
}

type FacilityJSON struct {
	FacilityName string `json:"name"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
}

type Facility struct {
	UserDB *sql.DB
}

func (f *Facility) GetFacilityInfo(userID string) (*[]FacilityJSON, error) {
	uuid := uuid.Parse(userID)
	row, err := f.UserDB.Query(`SELECT 
	facilities.name,facilities.address,facilities.phone 
	FROM facilities INNER JOIN facility_users 
	ON facilities.id = facility_users.facility_id 
	WHERE facility_users.user_id = ?;`, uuid[:])
	if err != nil {
		slog.Error("failed in get facility", "error", err)
		return nil, err
	}
	var facilities []FacilityJSON
	for row.Next() {
		var facility FacilityJSON
		if err := row.Scan(&facility.FacilityName, &facility.Address, &facility.Phone); err != nil {
			slog.Error("can't read facility row", "error", err)
			return nil, err
		}
		facilities = append(facilities, facility)
	}
	return &facilities, nil
}
