package main

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"scorepower4cours/routes"
)

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

	r.POST("/score", func(c *gin.Context) { routes.PublishScore(c, db) })
	r.GET("/score", func(c *gin.Context) { routes.GetScores(c, db) })
	r.GET("/generateUserId", routes.GenerateUserId)
	r.DELETE("/clear", func(c *gin.Context) { routes.Clear(c, db) })

	log.Fatalln(r.Run())
}
