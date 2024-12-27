package main

import (
	"log"
	"server/config"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Load server port from config
	port := config.LoadConfig().ServerPort
	if port == "" {
		port = "8080"
	}

	// Register routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Start the server
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	} else {
		log.Printf("Server started on port %s", port)
	}
}
