package main

import (
	"database/sql"
	"errors"
	"log"
)

func queryRoomsBuilding(building string, db *sql.DB) []dbRoom {
	query := `SELECT * FROM room_info_everything WHERE building = ?`

	rows, err := db.Query(query, building)
	if err != nil {
		log.Fatalf("error during query: %v\n", err)
	}

	rooms := []dbRoom{}
	for rows.Next() {
		var room dbRoom
		var tmp interface{}
		err = rows.Scan(&tmp, &room.Name, &room.Building, &room.RoomId,
			&room.Ds1,
			&room.Ds2,
			&room.Ds3,
			&room.Ds4,
			&room.Ds5,
			&room.Ds6,
			&room.Ds7,
			&room.Ds8,
			&room.Ds9,
			&room.Ds10)
		if err != nil {
			log.Fatal("(GetAllRoomsForBuilding) res.Scan", err)
		}
		rooms = append(rooms, room)
	}

	return rooms
}

func queryOneRoom(building, roomName string, db *sql.DB) (dbRoom, error) {
	query := `SELECT * FROM room_info_everything WHERE name = ? AND building = ?`
	var room dbRoom

	// tmp *interface, because we don't need the ID from the database row
	var tmp interface{}
	err := db.QueryRow(query, roomName, building).Scan(&tmp, &room.Name, &room.Building, &room.RoomId,
		&room.Ds1,
		&room.Ds2,
		&room.Ds3,
		&room.Ds4,
		&room.Ds5,
		&room.Ds6,
		&room.Ds7,
		&room.Ds8,
		&room.Ds9,
		&room.Ds10)

	if errors.Is(err, sql.ErrNoRows) {
		return dbRoom{}, sql.ErrNoRows
	} else if err != nil {
		log.Fatal("(GetRoomInfo) db.Exec", err)
	}

	return room, nil
}

// FetchAllRooms fetches the HTML code and parses it into the Room datastructure for all rooms.
func fetchAllRooms() (rooms []Room) {

	rawRooms := GetAllRooms()
	for _, rawR := range rawRooms {
		rooms = append(rooms, parseIntoRoom(rawR))
	}

	return rooms
}

// RoomInfo holds the lecture plan information for a room.
// Each Ds is a slice of length 5 containing the lectures or exercices for each weekday.
// An empty strings means that the room is free to use in the time slot.
type RoomInfo struct {
	Ds1  []string
	Ds2  []string
	Ds3  []string
	Ds4  []string
	Ds5  []string
	Ds6  []string
	Ds7  []string
	Ds8  []string
	Ds9  []string
	Ds10 []string
}
