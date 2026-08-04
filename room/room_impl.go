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
	roomID   string `db:"id"`
	roomName string `db:"room_name"`
	roomType string `db:"room_type_name"`
}

func (roomStruct *getRoomStruct) GetRoomAll(facilityID string) ([]*Room, error) {
	parsedFalicityID, err := uuid.Parse(facilityID)
	if err != nil {
		return nil, err
	}
	sql := `
	SELECT room.id,room.room_name,room_type.room_type_name FROM room
	INNER JOIN room_type ON room.id = room_type.room_id
	WHERE room.facility_id = ?
	`
	var roomDatabases []roomDatabaseValue

	err = roomStruct.db.UserDB.Select(&roomDatabases, sql, parsedFalicityID[:])
	if err != nil {
		return nil, err
	}
	var Rooms []*Room
	for _, row := range roomDatabases {
		parsedRoomID, err := uuid.FromBytes([]byte(row.roomID))
		if err != nil {
			return nil, err
		}
		room := &Room{
			RoomID:   parsedRoomID.String(),
			RoomName: row.roomName,
			RoomType: row.roomType,
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
	SELECT room.id,room.room_name,room_type.room_type_name FROM room
	INNER JOIN room_type ON room.id = room_type.room_id
	WHERE room.facility_id = ? AND room.id = ?
	`
	var RoomVal *roomDatabaseValue
	err = roomStruct.db.UserDB.Get(&RoomVal, sql, parsedFacilityID[:], parsedRoomID)
	if err != nil {
		return nil, err
	}
	parsedRoomID2, err := uuid.FromBytes([]byte(RoomVal.roomID))
	if err != nil {
		return nil, err
	}
	RoomAnser := &Room{
		RoomID:   parsedRoomID2.String(),
		RoomName: RoomVal.roomName,
		RoomType: RoomVal.roomType,
	}

	return RoomAnser, nil
}
