package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindFreeRoom(c *gin.Context, db *sql.DB) {
	building := c.Query("building")
	ds := c.Query("ds")

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("you requested free rooms in building(s) %s during lecture(s) %s", building, ds)})
}
func getCurrentDs() {

}
