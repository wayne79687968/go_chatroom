package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gin-contrib/cors"
)

var ctx = context.Background()

func generateUsername(c *gin.Context, rdb *redis.Client) {
	val, err := rdb.Get(ctx, "username").Result()
	if err == redis.Nil {
		fmt.Println("No username, generating new one")
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(9999) + 1
		val = fmt.Sprintf("Anonymous#%d", num)

		err := rdb.Set(ctx, "username", val, 0).Err()
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"username": val,
	})
}


func main() {
	r := gin.Default()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	hub := newHub()
	go hub.run()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // Allow all origins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	r.Use(cors.New(corsConfig))

	r.GET("/api/getUsername", func(c *gin.Context) {
		generateUsername(c, rdb)
	})

  r.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	})

	r.Run()
}