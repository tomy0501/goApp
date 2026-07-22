package hotelindex

import (
	"goapp/db"

	"github.com/google/uuid"
)

type hotelIndexRepository struct {
	db db.UserDB
}

func NewHotelIndexRepository(db db.UserDB) HotelIndexInterface {
	return &hotelIndexRepository{
		db: db,
	}
}

type CombineHotelIndexRow struct {
	FacilityRole string `db:"role"`
	FacilityName string `db:"name"`
}

func (h *hotelIndexRepository) GetHotelIndexInfo(facilityID string, userID string) (*HotelIndex, error) {
	parsedFacilityID, err := uuid.Parse(facilityID)
	if err != nil {
		return nil, err
	}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	var combineRow CombineHotelIndexRow

	err = h.db.UserDB.Select(&combineRow.FacilityRole, combineRow.FacilityName, `SELECT facility_users.role,facilities.name FROM facility_users
INNER JOIN facilities on facilities.id = facility_users.facility_id WHERE facility_users.user_id = ? AND facility_users.facilityID = ?`,
		parsedUserID[:], parsedFacilityID)
	if err != nil {
		return nil, err
	}
	hotelIndex := HotelIndex{
		HotelName: combineRow.FacilityName,
		Authority: combineRow.FacilityRole,
	}
	return &hotelIndex, nil
}
