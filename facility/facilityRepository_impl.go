package facility

import (
	"goapp/db"
	"log/slog"

	"github.com/google/uuid"
)

type facilityRepository struct {
	db db.UserDB
}

func NewFacilityRepository(db db.UserDB) FacilityInterface {
	return &facilityRepository{
		db: db,
	}
}
func (d *facilityRepository) GetFacilityInfo(userID string) (*[]FacilityJSON, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		slog.Error("uuid can't parse", "error", err)
		return nil, err
	}
	row, err := d.db.UserDB.Query(`SELECT 
	facilities.id,facilities.name,facilities.address,facilities.phone 
	FROM facilities INNER JOIN facility_users 
	ON facilities.id = facility_users.facility_id 
	WHERE facility_users.user_id = ?;`, parsedUserID[:])
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
	var facilities []FacilityJSON
	for row.Next() {
		var facility FacilityJSON
		var rawFacilityID []byte
		if err := row.Scan(&rawFacilityID, &facility.FacilityName, &facility.Address, &facility.Phone); err != nil {
			slog.Error("can't read facility row", "error", err)
			return nil, err
		}
		parsedFacilityID, err := uuid.FromBytes(rawFacilityID)
		if err != nil {
			slog.Error("can't read faclity ID", "error", err)
			return nil, err
		}
		facility.FacilityID = parsedFacilityID.String()
		facilities = append(facilities, facility)
	}
	return &facilities, nil
}
