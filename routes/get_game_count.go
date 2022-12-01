package routes

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"net/http"
	"scorepower4cours/api"
	badger2 "scorepower4cours/badger"
)

func GetGameCount(c *gin.Context, db *badger.DB) {
	count, err := countGameCount(db, c.Query("hostId"))
	if err != nil {
		c.String(http.StatusInternalServerError, "internal error: "+err.Error())
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("%d", count))
}

func countGameCount(db *badger.DB, hostId string) (uint, error) {
	var count uint
	if err := badger2.IterateRecords(db, func(record api.ScoreRecord, i int) {
		if record.HostId.String() == hostId {
			count++
		}
	}); err != nil {
		return 0, err
	}
	return count, nil
}
