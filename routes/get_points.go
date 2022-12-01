package routes

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"net/http"
	"scorepower4cours/api"
	badger2 "scorepower4cours/badger"
	"scorepower4cours/utils"
)

func GetPoints(c *gin.Context, db *badger.DB) {
	hostId := c.Query("hostId")
	if hostId == "" {
		c.String(http.StatusBadRequest, "hostId is required")
		return
	}

	points, err := getHostPoints(db, hostId)
	if err != nil {
		c.String(http.StatusInternalServerError, "internal error: "+err.Error())
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("%d", points))
	return
}

// iterateRecords iterates over all records in the database and calls the callback for each record
func getHostPoints(db *badger.DB, hostId string) (uint, error) {
	var points uint
	var gameCount int
	if err := badger2.IterateRecords(db, func(record api.ScoreRecord, i int) {
		fmt.Println(record.HostId.String(), hostId, record.HostId.String() == hostId)
		if record.HostId.String() == hostId && record.IsHostWinner() {
			points += utils.IfThenElse[uint](gameCount%3 == 0, 2, 1)
			gameCount++
		}
	}); err != nil {
		return 0, err
	}
	return points, nil
}
