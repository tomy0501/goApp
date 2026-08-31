package room

import (
	"goapp/db"

	"github.com/google/uuid"
)

type getRoomStruct struct {
	db *db.UserDB
}

func NewGetRoomStruct(db *db.UserDB) GetRoomInterface {
	return &getRoomStruct{
		db: db,
	}
}

type roomDatabaseValue struct {
	RoomID     string `db:"id"`
	RoomName   string `db:"room_name"`
	RoomType   string `db:"room_type_name"`
	RoomTypeID string `db:"room_type_id"`
}

func (roomStruct *getRoomStruct) GetRoomAll(facilityID string) ([]*Room, error) {
	parsedFalicityID, err := uuid.Parse(facilityID)
	if err != nil {
		return nil, err
	}
	sql := `
	SELECT room.id,room.room_name,room_type.room_type_name,room_type.id as room_type_id FROM room
	INNER JOIN room_type ON room.id = room_type.room_id
	WHERE room_type.facility_id = ?
	`
	var roomDatabases []roomDatabaseValue

	err = roomStruct.db.UserDB.Select(&roomDatabases, sql, parsedFalicityID[:])
	if err != nil {
		return nil, err
	}
	var Rooms []*Room
	for _, row := range roomDatabases {
		parsedRoomID, err := uuid.FromBytes([]byte(row.RoomID))
		if err != nil {
			return nil, err
		}
		parsedRoomTypeID, err := uuid.FromBytes([]byte(row.RoomTypeID))
		if err != nil {
			return nil, err
		}
		room := &Room{
			RoomID:     parsedRoomID.String(),
			RoomName:   row.RoomName,
			RoomType:   row.RoomType,
			RoomTypeID: parsedRoomTypeID.String(),
		}
		Rooms = append(Rooms, room)
	}
	return Rooms, nil
}

func (roomStruct getRoomStruct) GetRoomById(facilityID string, roomID string) (*Room, error) {
	parsedFacilityID, err := uuid.Parse(facilityID)
	if err != nil {
		return nil, err
	}
	parsedRoomID, err := uuid.Parse(roomID)
	if err != nil {
		return nil, err
	}
	sql := `
	SELECT room.id,room.room_name,room_type.room_type_name,room_type.id as room_type_id FROM room
	INNER JOIN room_type ON room.id = room_type.room_id
	WHERE room_type.facility_id = ? AND room.id = ?
	`
	var RoomVal *roomDatabaseValue
	err = roomStruct.db.UserDB.Get(&RoomVal, sql, parsedFacilityID[:], parsedRoomID)
	if err != nil {
		return nil, err
	}
	parsedRoomID2, err := uuid.FromBytes([]byte(RoomVal.RoomID))
	if err != nil {
		return nil, err
	}
	parsedRoomTypeID, err := uuid.FromBytes([]byte(RoomVal.RoomTypeID))
	if err != nil {
		return nil, err
	}
	RoomAnser := &Room{
		RoomID:     parsedRoomID2.String(),
		RoomName:   RoomVal.RoomName,
		RoomType:   RoomVal.RoomType,
		RoomTypeID: parsedRoomTypeID.String(),
	}

	return RoomAnser, nil
}
