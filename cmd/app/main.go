package main

import (
	"log"
	"server/config"
	"server/internal/repositories"
	pg "server/pkg/utils"
	"time"

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

	facilityRepo := repositories.NewFacilityRepository(db)

	// Start timer
	start := time.Now()

	// Fetch all facilities
	facilities, err := facilityRepo.FindMany(nil) // Pass nil to fetch all facilities
	if err != nil {
		log.Printf("Error fetching facilities: %v", err)
	} else {
		for _, facility := range facilities {
			log.Printf("Facility Details:")
			log.Printf("  ID: %d", facility.ID)
			log.Printf("  Name: %s", facility.Name)
			log.Printf("  Type: %s", facility.Type)
			log.Printf("  Location: %s", facility.Location)
			log.Printf("  Phone: %s", facility.Phone)
			log.Printf("  Email: %s", facility.Email)

			if facility.Coordinates != nil {
				log.Printf("  Coordinates: %s", *facility.Coordinates)
			} else {
				log.Printf("  Coordinates: NULL")
			}

			if facility.Description != nil {
				log.Printf("  Description: %s", *facility.Description)
			} else {
				log.Printf("  Description: NULL")
			}

			log.Println("----") // Separator between facilities for better readability
		}
	}

	// End timer
	duration := time.Since(start)
	log.Printf("Fetching and printing facilities took %s", duration)

	// Start the server
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	} else {
		log.Printf("Server started on port %s", port)
	}
}
