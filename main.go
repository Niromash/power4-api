package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"time"
)

type ScoreRecord struct {
	PlayerRed    int       `json:"playerRed" validate:"required"`
	PlayerYellow int       `json:"playerYellow" validate:"required"`
	BoardSizeX   int       `json:"boardSizeX" validate:"required"`
	BoardSizeY   int       `json:"boardSizeY" validate:"required"`
	UserId       uuid.UUID `json:"userId" validate:"required"`
}

func (s ScoreRecord) Validate() error {
	return validator.New().Struct(s)
}

func main() {
	db, err := badger.Open(badger.DefaultOptions("./db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.POST("/score", func(c *gin.Context) {
		var scoreRecord ScoreRecord
		if err = c.ShouldBindJSON(&scoreRecord); err != nil {
			c.String(http.StatusBadRequest, "invalid body: "+err.Error())
			return
		}

		if err = scoreRecord.Validate(); err != nil {
			c.String(http.StatusBadRequest, "invalid body: "+err.Error())
			return
		}

		marshalledBody, err := json.Marshal(scoreRecord)
		if err != nil {
			c.String(http.StatusInternalServerError, "internal error: "+err.Error())
			return
		}

		if err = db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(fmt.Sprintf("%d", time.Now().UnixNano())), marshalledBody)
		}); err != nil {
			c.String(http.StatusInternalServerError, "internal error: "+err.Error())
			return
		}
	})

	iterateRecords := func(callback func(record ScoreRecord)) error {
		return db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				err := item.Value(func(val []byte) error {
					var record ScoreRecord
					if err = json.Unmarshal(val, &record); err != nil {
						return err
					}
					callback(record)
					return nil
				})
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	r.GET("/score", func(c *gin.Context) {
		scores := []ScoreRecord{}
		if err = iterateRecords(func(record ScoreRecord) {
			userIdQuery := c.Query("userId")
			if len(userIdQuery) > 0 && record.UserId.String() != userIdQuery {
				return
			}
			scores = append(scores, record)
		}); err != nil {
			c.String(http.StatusInternalServerError, "internal error: "+err.Error())
			return
		}
		c.JSON(http.StatusOK, scores)
	})

	r.GET("/generateUserId", func(c *gin.Context) {
		c.String(http.StatusCreated, uuid.New().String())
	})

	r.GET("/getGameCount", func(c *gin.Context) {
		var count uint
		if err = iterateRecords(func(record ScoreRecord) {
			if record.UserId.String() == c.Query("userId") {
				count++
			}
		}); err != nil {
			c.String(http.StatusInternalServerError, "internal error: "+err.Error())
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("%d", count))
	})

	r.DELETE("/clear", func(c *gin.Context) {
		apiKey := c.GetHeader("apiKey")
		if apiKey != os.Getenv("API_KEY") {
			c.String(http.StatusUnauthorized, "unauthorized")
			return
		}
		if err = db.DropAll(); err != nil {
			c.String(http.StatusInternalServerError, "internal error: "+err.Error())
			return
		}
		c.Status(http.StatusNoContent)
	})

	r.Run()
}
