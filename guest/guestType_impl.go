package guest

import (
	"database/sql"
	"goapp/db"

	"github.com/google/uuid"
)

type guestRepository struct {
	db *db.UserDB
}

type GuestRepositoryRow struct {
	GuestID   string         `db:"guest_id"`
	GuestAge  sql.NullInt16  `db:"customer_age"`
	GuestName string         `db:"customer_name"`
	GuestSex  sql.NullString `db:"customer_sex"`
	GuestFrom sql.NullString `db:"customer_from"`
}

func NewGuestRepository(db *db.UserDB) GuestInterface {
	return &guestRepository{
		db: db,
	}
}
func (r *guestRepository) GetGuestAll(facilityID string) ([]*GuestOnly, error) {
	parsedFacilityID, err := uuid.Parse(facilityID)
	if err != nil {
		return nil, err
	}
	var guestList []*GuestRepositoryRow
	sql := `
		SELECT g.id AS customer_id,g.customer_name,g.customer_age,g.customer_sex,g.customer_from,middle.checkin FROM reservation 
		INNER JOIN reservation_guest AS middle ON reservation.id = middle.reservation_id
		INNER JOIN guest AS g ON middle.guest_id = g.id
		WHERE reservation.facility_id = ?
	`
	err = r.db.UserDB.Select(&guestList, sql, parsedFacilityID[:])
	if err != nil {
		return nil, err
	}
	var guests []*GuestOnly
	for _, g := range guestList {
		guests = append(guests, &GuestOnly{
			GuestID:   g.GuestID,
			GuestName: g.GuestName,
		})
		if g.GuestAge.Valid {
			age := int(g.GuestAge.Int16)
			guests[len(guests)-1].GuestAge = &age
		} else {
			guests[len(guests)-1].GuestAge = nil
		}
		if g.GuestSex.Valid {
			sex := g.GuestSex.String
			guests[len(guests)-1].GuestSex = &sex
		} else {
			guests[len(guests)-1].GuestSex = nil
		}
		if g.GuestFrom.Valid {
			from := g.GuestFrom.String
			guests[len(guests)-1].GuestFrom = &from
		} else {
			guests[len(guests)-1].GuestFrom = nil
		}
	}
	return guests, nil
}

func (r *guestRepository) GetGuestById(facilityID string, guestID string) (*GuestOnly, error) {
	parsedFacilityID, err := uuid.Parse(facilityID)
	if err != nil {
		return nil, err
	}
	parsedGuestID, err := uuid.Parse(guestID)
	if err != nil {
		return nil, err
	}
	var guest *GuestRepositoryRow
	sql := `
		SELECT g.id AS customer_id,g.customer_name,g.customer_age,g.customer_sex,g.customer_from,middle.checkin FROM reservation 
		INNER JOIN reservation_guest AS middle ON reservation.id = middle.reservation_id
		INNER JOIN guest AS g ON middle.guest_id = g.id
		WHERE reservation.facility_id = ? AND g.id = ?
	`
	err = r.db.UserDB.Get(&guest, sql, parsedFacilityID[:], parsedGuestID[:])
	if err != nil {
		return nil, err
	}
	var g *GuestOnly
	g = &GuestOnly{
		GuestID:   guest.GuestID,
		GuestName: guest.GuestName,
	}
	if guest.GuestAge.Valid {
		age := int(guest.GuestAge.Int16)
		g.GuestAge = &age
	} else {
		g.GuestAge = nil
	}
	if guest.GuestSex.Valid {
		sex := guest.GuestSex.String
		g.GuestSex = &sex
	} else {
		g.GuestSex = nil
	}
	if guest.GuestFrom.Valid {
		from := guest.GuestFrom.String
		g.GuestFrom = &from
	} else {
		g.GuestFrom = nil
	}
	return g, nil
}
