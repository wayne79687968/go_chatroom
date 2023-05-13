package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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

	dbUser := "user"
	dbPassword := "password"
	dbName := "chatroom"
	dbHost := "localhost:3306"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
		serveWs(hub, c.Writer, c.Request, db)
	})

	r.Run()
}
