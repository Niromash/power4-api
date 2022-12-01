package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func GenerateUserId(c *gin.Context) {
	c.String(http.StatusCreated, uuid.New().String())
}
