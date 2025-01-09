package handlers

import (
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

// Services groups all the service instances.
type Services struct {
	CityService     *services.CityService
	FacilityService *services.FacilityService
	AuthService     *services.AuthService // Add AuthService
	// Add other services here as needed
}

// RegisterHandlers initializes and registers all application handlers with the provided Gin router.
func RegisterHandlers(router *gin.Engine, services *Services) {
	api := router.Group("/api")

	// Initialize handlers
	cityHandler := NewCityHandler(services.CityService)
	facilityHandler := NewFacilityHandler(services.FacilityService)
	authHandler := NewAuthHandler(services.AuthService) // Initialize the AuthHandler

	// Register routes
	cityHandler.RegisterCityRoutes(api)
	facilityHandler.RegisterFacilityRoutes(api)
	authHandler.RegisterAuthRoutes(api)
}
