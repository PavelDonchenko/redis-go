package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Car struct {
	Model string `json:"model"`
	Year  uint32 `json:"year"`
	Color string `json:"color"`
	Email string `json:"email"`
}

func main() {
	var ctx = context.Background()
	server := gin.Default()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	fmt.Println("Redis client connected successfully...")

	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		// Sleep for 3 seconds and wait for Redis to initialize
		time.Sleep(3 * time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(pong)

	router := server.Group("/api")
	router.POST("/car", func(c *gin.Context) {
		var car Car

		err := c.ShouldBindJSON(&car)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		}

		payload, err := json.Marshal(car)
		if err != nil {
			fmt.Println(err)
		}

		err = redisClient.Publish(ctx, "send-car-data", payload).Err()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error sending in channel": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"response": car})
	})

	log.Fatal(server.Run(":8002"))
}
