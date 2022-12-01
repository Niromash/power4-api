package routes

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func Clear(c *gin.Context, db *badger.DB) {
	apiKey := c.GetHeader("apiKey")
	if apiKey != os.Getenv("API_KEY") {
		c.String(http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := db.DropAll(); err != nil {
		c.String(http.StatusInternalServerError, "internal error: "+err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
