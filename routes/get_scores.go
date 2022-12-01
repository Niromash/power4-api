package routes

import (
	badger2 "github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"net/http"
	"scorepower4cours/api"
	"scorepower4cours/badger"
)

func GetScores(c *gin.Context, db *badger2.DB) {
	scores := []api.ScoreRecord{}
	if err := badger.IterateRecords(db, func(record api.ScoreRecord, i int) {
		userIdQuery := c.Query("hostId")
		if len(userIdQuery) > 0 && record.HostId.String() != userIdQuery {
			return
		}
		scores = append(scores, record)
	}); err != nil {
		c.String(http.StatusInternalServerError, "internal error: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, scores)
}
