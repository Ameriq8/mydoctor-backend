package main

import (
	"log"
	"net/http"
	"server/config"
	"server/internal/handlers"
	"server/internal/repositories"
	"server/internal/services"
	"server/pkg/logger"
	"server/pkg/middlewares"
	pg "server/pkg/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	gin.SetMode(gin.ReleaseMode) // Set Gin to release mode for production
	r := gin.Default()

	// Error handling for SetTrustedProxies method
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Error setting trusted proxies: %v", err)
	}

	logger.InitLogger("development")
	defer logger.Sync()

	// Load server port from config
	port := config.LoadConfig().ServerPort
	if port == "" {
		port = "8080"
	}

	// Setup CORS config & middleware
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true, // Allow specific origins instead of all
		// AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		// ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add XSS Middleware for sanitizing input data
	// r.Use(middlewares.XSSMiddleware())

	// Add the monitoring middleware for Prometheus metrics
	r.Use(middlewares.MonitoringMiddleware())
	r.Use(middlewares.Logging())

	// Add secure headers middleware for enhanced security
	// secureConfig := secure.New(secure.Config{
	// 	SSLRedirect:           true,                                                        // Force HTTPS
	// 	FrameDeny:             true,                                                        // Deny iframe embedding (Clickjacking)
	// 	ContentTypeNosniff:    true,                                                        // Prevent MIME type sniffing
	// 	BrowserXssFilter:      true,                                                        // Enable browser's XSS filter
	// 	ContentSecurityPolicy: "default-src 'self'; script-src 'self'; object-src 'none';", // Apply Content Security Policy (CSP)
	// })
	// r.Use(secureConfig)

	// Database connection
	db, err := pg.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Expose Prometheus metrics at "/metrics"
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Register routes
	r.GET("/ping", func(c *gin.Context) {
		// Create a FacilityRepository instance using the existing db connection
		fRepo := repositories.NewFacilityRepository(db)

		// Fetch a facility from the repository
		f, err := fRepo.Find(1)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Facility not found"})
			return
		}

		c.JSON(200, gin.H{
			"message":  "pong",
			"facility": f,
		})
	})

	// Initialize repositories
	cityRepo := repositories.NewCitiesRepository(db)
	facilityRepo := repositories.NewFacilityRepository(db)
	authRepo := repositories.NewAuthRepository(db)

	// Initialize services
	serviceGroup := &handlers.Services{
		CityService:     services.NewCityService(cityRepo),
		FacilityService: services.NewFacilityService(facilityRepo),
		AuthService:     services.NewAuthService(authRepo),
	}

	// Register handlers
	handlers.RegisterHandlers(r, serviceGroup)

	logger.Info("Application started", zap.String("env", "development"))

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
