package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
)

var db *sql.DB

type Room struct {
	Id       int64
	Name     string
	Building string
	RoomId   string
	Ds1      string
	Ds2      string
	Ds3      string
	Ds4      string
	Ds5      string
	Ds6      string
	Ds7      string
	Ds8      string
	Ds9      string
	Ds10     string
}

func main() {
	// Load in the `.env` file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env", err)
	}

	// Open a connection to the database
	db, err = sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatal("failed to open db connection", err)
	}

	r := gin.Default()
	r.GET("/info", GetAllBuildings)
	r.GET("/updateRooms", FetchNewRoomInfoAndPersistToDB)
	r.GET("/rooms/:building", GetAllRoomsForBuilding)
	r.GET("/rooms/:building/:room", GetRoomInfo)
	r.Run(os.Getenv("URL")) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func GetAllRoomsForBuilding(c *gin.Context) {
	building := c.Param("building")
	building = strings.ToUpper(building)
	query := `SELECT * FROM room_info_everything WHERE building = ?`

	rows, err := db.Query(query, building)
	if err != nil {
		log.Fatalf("error during query: %v\n", err)
	}

	rooms := []Room{}
	for rows.Next() {
		var room Room
		err = rows.Scan(&room.Id, &room.Name, &room.Building, &room.RoomId,
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
			log.Fatal("(GetProducts) res.Scan", err)
		}
		rooms = append(rooms, room)
	}

	c.JSON(http.StatusOK, rooms)

}

func GetRoomInfo(c *gin.Context) {
	var room Room
	building := c.Param("building")
	building = strings.ToUpper(building)
	roomName := c.Param("room")
	roomName = strings.ToUpper(roomName)

	query := `SELECT * FROM room_info_everything WHERE name = ? AND building = ?`

	err := db.QueryRow(query, roomName, building).Scan(&room.Id, &room.Name, &room.Building, &room.RoomId,
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
		c.Status(http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		log.Fatal("(UpdateProduct) db.Exec", err)
	}

	// split the lectures into array
	var roomInfoResponse RoomInfoResponse
	roomInfoResponse.Ds1 = strings.Split(room.Ds1, "|")
	roomInfoResponse.Ds2 = strings.Split(room.Ds2, "|")
	roomInfoResponse.Ds3 = strings.Split(room.Ds3, "|")
	roomInfoResponse.Ds4 = strings.Split(room.Ds4, "|")
	roomInfoResponse.Ds5 = strings.Split(room.Ds5, "|")
	roomInfoResponse.Ds6 = strings.Split(room.Ds6, "|")
	roomInfoResponse.Ds7 = strings.Split(room.Ds7, "|")
	roomInfoResponse.Ds8 = strings.Split(room.Ds8, "|")
	roomInfoResponse.Ds9 = strings.Split(room.Ds9, "|")
	roomInfoResponse.Ds10 = strings.Split(room.Ds10, "|")

	c.JSON(http.StatusOK, roomInfoResponse)

}

type RoomInfoResponse struct {
	Id       int64
	Name     string
	Building string
	RoomId   string
	Ds1      []string
	Ds2      []string
	Ds3      []string
	Ds4      []string
	Ds5      []string
	Ds6      []string
	Ds7      []string
	Ds8      []string
	Ds9      []string
	Ds10     []string
}

func putRoomIntoDb(room Room, c *gin.Context) {
	var tempRoom Room

	log.Printf("checking room: %v\n", room)
	// check row exists already first
	checkIfExistsQuery := `SELECT * FROM room_info_everything WHERE room_id = ?`
	err := db.QueryRow(checkIfExistsQuery, room.RoomId).Scan(&tempRoom.Id, &tempRoom.Name, &tempRoom.Building, &tempRoom.RoomId,
		&tempRoom.Ds1,
		&tempRoom.Ds2,
		&tempRoom.Ds3,
		&tempRoom.Ds4,
		&tempRoom.Ds5,
		&tempRoom.Ds6,
		&tempRoom.Ds7,
		&tempRoom.Ds8,
		&tempRoom.Ds9,
		&tempRoom.Ds10)

	if errors.Is(sql.ErrNoRows, err) {
		log.Printf("inserting entry for room_id: %v\n", room.RoomId)
		// doesn't exist yet, we have to create it
		insertQuery := `INSERT INTO room_info_everything (name, building, room_id, ds1, ds2, ds3, ds4, ds5, ds6, ds7, ds8, ds9, ds10) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		res, err := db.Exec(insertQuery, room.Name, room.Building, room.RoomId,
			room.Ds1, room.Ds2, room.Ds3, room.Ds4, room.Ds5, room.Ds6, room.Ds7, room.Ds8, room.Ds9, room.Ds10)
		if err != nil {
			log.Fatal("(CreateRoom) db.Exec", err)
		}
		room.Id, err = res.LastInsertId()
		if err != nil {
			log.Fatal("(CreateRoom) res.LastInsertId", err)
		}
	} else if err != nil {
		log.Println(err)
		log.Println("epic fail occurred, sir")
	} else {
		log.Printf("updating table entry for room: %v\n", room.RoomId)
		// already exists, update record
		updateQuery := `UPDATE room_info_everything SET name = ?, building = ?, ds1 = ?, ds2 = ?, ds3 = ?, ds4 = ?, ds5 = ?, ds6 = ?, ds7 = ?, ds8 = ?, ds9 = ?, ds10 = ? WHERE room_id = ?`

		res, err := db.Exec(updateQuery, room.Name, room.Building,
			room.Ds1, room.Ds2, room.Ds3, room.Ds4, room.Ds5, room.Ds6, room.Ds7, room.Ds8, room.Ds9, room.Ds10,
			room.RoomId)
		if err != nil {
			log.Fatal("(UpdateRoom) db.Exec", err)
		}
		room.Id, err = res.LastInsertId()
		if err != nil {
			log.Fatal("(UpdateRoom) res.LastInsertId", err)
		}

	}

}

func FetchNewRoomInfoAndPersistToDB(c *gin.Context) {
	newRoom := Room{}

	for _, room := range GetRoomsWithWebPage() {
		rawData := parse(room.WebPage)
		info := transformIntoStruct(rawData)
		room := Room{
			Name:     room.Name,
			Building: room.Building,
			RoomId:   room.Id,
			Ds1:      strings.Join(info.Ds1, "|"),
			Ds2:      strings.Join(info.Ds2, "|"),
			Ds3:      strings.Join(info.Ds3, "|"),
			Ds4:      strings.Join(info.Ds4, "|"),
			Ds5:      strings.Join(info.Ds5, "|"),
			Ds6:      strings.Join(info.Ds6, "|"),
			Ds7:      strings.Join(info.Ds7, "|"),
			Ds8:      strings.Join(info.Ds8, "|"),
			Ds9:      strings.Join(info.Ds9, "|"),
			Ds10:     strings.Join(info.Ds10, "|"),
		}
		newRoom = room

		putRoomIntoDb(room, c)

	}

	c.JSON(http.StatusOK, newRoom)
}

func GetAllBuildings(c *gin.Context) {
	query := "SELECT * FROM room_info_everything"
	res, err := db.Query(query)
	defer res.Close()
	if err != nil {
		log.Fatal("(GetProducts) db.Query", err)
	}

	rooms := []Room{}
	for res.Next() {
		var room Room
		err := res.Scan(&room.Id, &room.Name, &room.Building, &room.RoomId,
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
			log.Fatal("(GetProducts) res.Scan", err)
		}
		rooms = append(rooms, room)
	}

	var distinctBuildings []string
	for _, room := range rooms {
		if !slices.Contains(distinctBuildings, room.Building) {
			distinctBuildings = append(distinctBuildings, room.Building)
		}
	}

	c.JSON(http.StatusOK, distinctBuildings)
}

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
