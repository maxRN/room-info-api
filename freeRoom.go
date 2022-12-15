package main

import (
	"log"
	"strconv"
	"time"
)

type roomTime struct {
	date time.Time
	room Room
	ds   []int
}

// freeRoomBuilding returns a list of rooms that are free right now
// and for how long they will be free.
func freeRoomBuilding(building string) []roomTime {
	return []roomTime{}

}

func getScheduleForDs(room Room, ds int) []string {
	switch ds {
	case 1:
		return room.Ds1
	case 2:
		return room.Ds2
	case 3:
		return room.Ds3
	case 4:
		return room.Ds4
	case 5:
		return room.Ds5
	case 6:
		return room.Ds6
	case 7:
		return room.Ds7
	case 8:
		return room.Ds8
	case 9:
		return room.Ds9
	case 10:
		return room.Ds10
	}

	panic("can't happen")
}

func FindFreeRooms(rooms []Room, d time.Weekday, ds int) []Room {
	log.Println(d)
	if d == 0 || d == 6 {
		return rooms
	}

	freeRooms := []Room{}
	for _, room := range rooms {
		currSched := getScheduleForDs(room, ds)
		log.Println(currSched)
		lec := currSched[d-1]
		log.Println(lec)
		if lec == "" {
			freeRooms = append(freeRooms, room)
		}
	}

	return freeRooms
}

// This is also wasted I think....
func isDuringDs(ds int, date time.Time) bool {
	n := time.Now()
	firstDs := time.Date(n.Year(), n.Month(), n.Day(), 7, 30, 0, 0, n.Location())
	timeSinceFirstDs := 110 * (ds - 1)
	dur, err := time.ParseDuration(strconv.Itoa(timeSinceFirstDs) + "m")
	if err != nil {
		log.Printf("big oof")
	}
	firstDs = firstDs.Add(dur)

	dsEnd := time.Date(n.Year(), n.Month(), n.Day(), 9, 0, 0, 0, n.Location())
	dsEnd = dsEnd.Add(dur)

	if date.After(firstDs) && date.Before(dsEnd) {
		return true
	} else {
		return false
	}
}

func getCurrentDs(date time.Time) int {
	currentDS := 1
	// find current DS
	if isDuringDs(1, date) {
		currentDS = 1
	} else if isDuringDs(2, date) {
		currentDS = 2
	} else if isDuringDs(3, date) {
		currentDS = 3
	} else if isDuringDs(4, date) {
		currentDS = 4
	} else if isDuringDs(5, date) {
		currentDS = 5
	} else if isDuringDs(6, date) {
		currentDS = 6
	} else if isDuringDs(7, date) {
		currentDS = 7
	} else if isDuringDs(8, date) {
		currentDS = 8
	} else if isDuringDs(9, date) {
		currentDS = 9
	} else if isDuringDs(10, date) {
		currentDS = 10
	} else {
		currentDS = 0
	}
	return currentDS
}
