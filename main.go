package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Room struct {
	Name     string // The common name for a room. This is the name that is listed on the plate next to the room's door.
	Building string
	RoomId   string // The ID of the room as it is saved in the Campus Navigator site.
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

var db *sql.DB

func main() {
	// Load in the `.env` file
	_ = godotenv.Load(".env.local")

	// Open a connection to the database
	var err error
	db, err = sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatal("failed to open db connection", err)
	}

	r := gin.Default()
	r.GET("/info", getAllBuildings)
	r.GET("/freeRoom", findFreeRoom)
	r.GET("/updateRooms", updateRooms)
	r.GET("/rooms/:building", getAllRoomsForBuilding)
	r.GET("/rooms/:building/:room", getRoomInfo)
	r.Run(os.Getenv("URL")) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getAllBuildings(c *gin.Context) {
	buildings := []string{"APB"}
	c.JSON(http.StatusOK, buildings)
}

func updateRooms(c *gin.Context) {
	API_KEY := os.Getenv("API_KEY")
	authHeader := c.GetHeader("Authentication")

	if API_KEY != authHeader {
		log.Println("api keys dont match!")
		c.AbortWithStatus(401)
		return
	}

	rooms := fetchAllRooms()
	for _, room := range rooms {
		putRoomIntoDb(room)
	}

	c.Status(http.StatusNoContent)
}

func getRoomsForBuildings(buildings []string) (rs []Room) {
	for _, b := range buildings {
		dbRooms := queryRoomsBuilding(b, db)
		for _, dbRoom := range dbRooms {
			rs = append(rs, dbRoomIntoRoom(dbRoom))
		}
	}
	return
}

func findFreeRoom(c *gin.Context) {
	log.Println("testing")
	// check params and do error handling
	// TODO: check if building exists
	building := c.Query("building")
	periodParam := c.Query("period")
	period, err := strconv.Atoi(periodParam)
	if err != nil {
		// TODO: check if period is valid
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expected a number for the period"})
	}
	log.Println(building, period)

	rooms := getRoomsForBuildings([]string{building})

	freeRooms := FindFreeRooms(rooms, time.Now().Weekday(), period)

	if len(freeRooms) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, freeRooms)
}

func getRoomInfo(c *gin.Context) {
	building := strings.ToUpper(c.Param("building"))
	roomName := strings.ToUpper(c.Param("room"))
	dbRoom, err := queryOneRoom(building, roomName, db)

	if errors.Is(err, sql.ErrNoRows) {
		c.Status(http.StatusNotFound)
		return
	} else if err != nil {
		log.Fatal("(GetRoomInfo) db.Exec", err)
	}

	room := dbRoomIntoRoom(dbRoom)
	c.JSON(http.StatusOK, room)
}

func getAllRoomsForBuilding(c *gin.Context) {
	building := strings.ToUpper(c.Param("building"))
	if building == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var rooms []Room
	dbRooms := queryRoomsBuilding(building, db)
	for _, dbRoom := range dbRooms {
		rooms = append(rooms, dbRoomIntoRoom(dbRoom))
	}

	c.JSON(http.StatusOK, rooms)
}

type dbRoom struct {
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

func prepareRoomForDB(room Room) dbRoom {

	dbRoom := dbRoom{
		Name:     room.Name,
		Building: room.Building,
		RoomId:   room.RoomId,
		Ds1:      strings.Join(room.Ds1, "|"),
		Ds2:      strings.Join(room.Ds2, "|"),
		Ds3:      strings.Join(room.Ds3, "|"),
		Ds4:      strings.Join(room.Ds4, "|"),
		Ds5:      strings.Join(room.Ds5, "|"),
		Ds6:      strings.Join(room.Ds6, "|"),
		Ds7:      strings.Join(room.Ds7, "|"),
		Ds8:      strings.Join(room.Ds8, "|"),
		Ds9:      strings.Join(room.Ds9, "|"),
		Ds10:     strings.Join(room.Ds10, "|"),
	}

	return dbRoom
}

func dbRoomIntoRoom(dbRoom dbRoom) (room Room) {

	room = Room{
		Name:     dbRoom.Name,
		Building: dbRoom.Building,
		RoomId:   dbRoom.RoomId,
		Ds1:      strings.Split(dbRoom.Ds1, "|"),
		Ds2:      strings.Split(dbRoom.Ds2, "|"),
		Ds3:      strings.Split(dbRoom.Ds3, "|"),
		Ds4:      strings.Split(dbRoom.Ds4, "|"),
		Ds5:      strings.Split(dbRoom.Ds5, "|"),
		Ds6:      strings.Split(dbRoom.Ds6, "|"),
		Ds7:      strings.Split(dbRoom.Ds7, "|"),
		Ds8:      strings.Split(dbRoom.Ds8, "|"),
		Ds9:      strings.Split(dbRoom.Ds9, "|"),
		Ds10:     strings.Split(dbRoom.Ds10, "|"),
	}

	return room
}

func putRoomIntoDb(room Room) {
	dbRoom := prepareRoomForDB(room)
	var tempRoom Room

	// check row exists already first
	checkIfExistsQuery := `SELECT * FROM room_info_everything WHERE room_id = ?`
	var tmp interface{}
	err := db.QueryRow(checkIfExistsQuery, dbRoom.RoomId).Scan(&tmp, &tempRoom.Name, &tempRoom.Building, &tempRoom.RoomId,
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
		log.Printf("inserting entry for room_id: %v\n", dbRoom.RoomId)
		// doesn't exist yet, we have to create it
		insertQuery := `INSERT INTO room_info_everything (name, building, room_id, ds1, ds2, ds3, ds4, ds5, ds6, ds7, ds8, ds9,
                ds10) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		res, err := db.Exec(insertQuery, dbRoom.Name, dbRoom.Building, dbRoom.RoomId,
			dbRoom.Ds1, dbRoom.Ds2, dbRoom.Ds3, dbRoom.Ds4, dbRoom.Ds5, dbRoom.Ds6, dbRoom.Ds7, dbRoom.Ds8, dbRoom.Ds9, dbRoom.Ds10)
		if err != nil {
			log.Fatal("(CreateRoom) db.Exec", err)
		}
		_, err = res.LastInsertId()
		if err != nil {
			log.Fatal("(CreateRoom) res.LastInsertId", err)
		}
	} else if err != nil {
		log.Println(err)
		log.Println("epic fail occurred, sir")
	} else {
		log.Printf("updating table entry for room: %v\n", dbRoom.RoomId)
		// already exists, update record
		updateQuery := `UPDATE room_info_everything SET name = ?, building = ?, ds1 = ?, ds2 = ?, ds3 = ?, ds4 = ?, ds5 = ?, ds6 = ?, ds7 = ?,
                    ds8 = ?, ds9 = ?, ds10 = ? WHERE room_id = ?`

		res, err := db.Exec(updateQuery, dbRoom.Name, dbRoom.Building,
			dbRoom.Ds1, dbRoom.Ds2, dbRoom.Ds3, dbRoom.Ds4, dbRoom.Ds5, dbRoom.Ds6, dbRoom.Ds7, dbRoom.Ds8, dbRoom.Ds9, dbRoom.Ds10,
			dbRoom.RoomId)
		if err != nil {
			log.Fatal("(UpdateRoom) db.Exec", err)
		}
		_, err = res.LastInsertId()
		if err != nil {
			log.Fatal("(UpdateRoom) res.LastInsertId", err)
		}

	}

}
