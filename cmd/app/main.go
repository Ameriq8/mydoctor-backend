package main

import (
	"log"
	"net/http"
	"server/config"
	"server/internal/repositories"
	"server/pkg/logger"
	"server/pkg/middlewares"
	pg "server/pkg/utils"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	gin.SetMode(gin.DebugMode) // Set Gin to release mode
	r := gin.Default()
	r.SetTrustedProxies(nil) // Disables proxy trusting

	logger.InitLogger("development")
	defer logger.Sync()

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

	// Add the monitoring middleware
	r.Use(middlewares.MonitoringMiddleware())

	// Expose Prometheus metrics at "/metrics"
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Register routes
	r.GET("/ping", func(c *gin.Context) {
		// Create a FacilityRepository instance using the existing db connection
		fRepo := repositories.NewFacilityRepository(db)

		// Fetch a facility from the repository
		f, err := fRepo.Find(130)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Facility not found"})
			return
		}

		c.JSON(200, gin.H{
			"message":  "pong",
			"facility": f,
		})
	})

	logger.Info("Application started", zap.String("env", "development"))

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
