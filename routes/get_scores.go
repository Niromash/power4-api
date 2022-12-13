package routes

import (
	badger2 "github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"net/http"
	"scorepower4cours/api"
	"scorepower4cours/badger"
)

func GetScores(c *gin.Context, db *badger2.DB) {
	scores := []api.GameRecord{}
	if err := badger.IterateRecords(db, func(record api.GameRecord, i int) {
		userIdQuery := c.Query("hostId")
		gameIdQuery := c.Query("gameId")
		if len(gameIdQuery) > 0 && record.GameId.String() != gameIdQuery {
			return
		}
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
