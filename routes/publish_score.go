package routes

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"net/http"
	"scorepower4cours/api"
	"time"
)

// PublishScore publishes a score to the database and returns the player points
func PublishScore(c *gin.Context, db *badger.DB) {
	var scoreRecord api.ScoreRecord
	if err := c.ShouldBindJSON(&scoreRecord); err != nil {
		c.String(http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}

	if err := scoreRecord.Validate(); err != nil {
		c.String(http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}

	scoreRecord.Date = time.Now()

	marshalledBody, err := json.Marshal(scoreRecord)
	if err != nil {
		c.String(http.StatusInternalServerError, "internal error: "+err.Error())
		return
	}

	if err = db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(fmt.Sprintf("%d", scoreRecord.Date.UnixNano())), marshalledBody)
	}); err != nil {
		c.String(http.StatusInternalServerError, "internal error: "+err.Error())
		return
	}

	points, err := getHostPoints(db, scoreRecord.HostId.String())
	if err != nil {
		c.String(http.StatusInternalServerError, "internal error: "+err.Error())
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("%d", points))
}
