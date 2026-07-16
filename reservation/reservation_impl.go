package reservation

import (
	"database/sql"
	"goapp/db"
	"goapp/guest"
	"goapp/room"
	"log/slog"

	"github.com/google/uuid"
)

type reservationRepository struct {
	db db.UserDB
}

func NewReservationRepository(db db.UserDB) ReservationInterface {
	return &reservationRepository{
		db: db,
	}
}

type CombineReservationRow struct {
	ReservationID string `db:"reservation_id"`
	RepresentID   string `db:"represent_id"`
	GuestID       string `db:"guest_id"`
	GuestAge      sql.NullInt16
	GuestName     string         `db:"customer_name"`
	GuestSex      sql.NullString `db:"customer_sex"`
	GuestFrom     sql.NullString `db:"customer_from"`
	RoomID        string         `db:"room_id"`
	RooomName     string         `db:"room_name"`
	RoomType      string         `db:"room_type_name"`
	CheckIn       sql.NullString `db:"checkin"`
	CheckOut      sql.NullString `db:"checkout"`
}

func (mysql *reservationRepository) GetAllReservation(facilityID string) ([]*Reservation, error) {
	var rawReservationData []CombineReservationRow
	parsedFacilityID, err := uuid.Parse(facilityID)
	if err != nil {
		slog.Error("can't parse facilityID to byte")
		return nil, err
	}
	err = mysql.db.UserDB.Select(&rawReservationData,
		`
		SELECT reservation.id AS reservation_id,reservation.represent_guest_id AS represent_id,
		guest.id AS guest_id,
		guest.customer_name,guest.customer_sex,guest.customer_from,
		room.id AS room_id,room.room_name,reservation_guest.checkin,reservation_guest.checkout,
		room_type.room_type_name
		FROM reservation INNER JOIN reservation_guest ON reservation.id = reservation_guest.reservation_id
		INNER JOIN guest ON reservation_guest.guest_id = guest.id
		INNER JOIN room ON reservation_guest.room_id = room.id
		INNER JOIN room_type ON room.id = room_type.room_id
		WHERE reservation.facility_id = ?
		ORDER BY reservation_id
		`, parsedFacilityID[:])
	if err != nil {
		slog.Error("can't get reservation", "error", err)
		return nil, err
	}
	resMap := make(map[string]*Reservation)
	var reservations []*Reservation

	for _, row := range rawReservationData {
		parsedReservationID, err := uuid.FromBytes([]byte(row.ReservationID))
		if err != nil {
			slog.Error("invalid reservationID type please check mysql")
			return nil, err
		}
		parsedGuestID, err := uuid.FromBytes([]byte(row.GuestID))
		if err != nil {
			slog.Error("invalid guestID type please check mysql")
			return nil, err
		}
		res, resExist := resMap[parsedFacilityID.String()]
		if !resExist {

			res = &Reservation{
				ReservationID:  parsedReservationID.String(),
				RepresentGuest: nil,
				Guests:         []guest.Guest{},
			}
			resMap[parsedFacilityID.String()] = res
			reservations = append(reservations, res)
		}
		room := room.Room{
			RoomID:   row.RoomID,
			RoomName: row.RooomName,
			RoomType: row.RoomType,
		}
		if row.RepresentID == row.GuestID {
			res.RepresentGuest = &guest.Guest{
				GuestID:   parsedGuestID.String(),
				GuestName: row.GuestName,
				GuestSex:  row.GuestSex,
				GuestAge:  row.GuestAge,
				GuestFrom: row.GuestFrom,
				Room:      room,
			}
		}
		res.Guests = append(res.Guests, guest.Guest{
			GuestID:   parsedGuestID.String(),
			GuestName: row.GuestName,
			GuestSex:  row.GuestSex,
			GuestAge:  row.GuestAge,
			GuestFrom: row.GuestFrom,
			Room:      room,
		})
	}
	return reservations, nil
}
