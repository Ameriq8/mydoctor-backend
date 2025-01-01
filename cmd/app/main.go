package main

import (
	"log"
	"server/config"
	pg "server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode) // Set Gin to release mode
	r := gin.Default()
	r.SetTrustedProxies(nil) // Disables proxy trusting

	// Load server port from config
	port := config.LoadConfig().ServerPort
	if port == "" {
		port = "8080"
	}

	db, err := pg.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
