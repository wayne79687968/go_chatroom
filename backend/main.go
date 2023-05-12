package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

var ctx = context.Background()

func generateUsername(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(9999) + 1
	username := fmt.Sprintf("Anonymous#%d", num)

	c.JSON(200, gin.H{"username": username})
}


func main() {
	r := gin.Default()

	hub := newHub()
	go hub.run()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // Allow all origins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	r.Use(cors.New(corsConfig))

	r.GET("/api/getUsername", func(c *gin.Context) {
		generateUsername(c)
	})

  	r.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	})

	r.Run()
}