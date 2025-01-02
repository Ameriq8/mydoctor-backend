package handlers

import (
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

// Services groups all the service instances.
type Services struct {
	CityService     *services.CityService
	FacilityService *services.FacilityService
	// Add other services here as needed
}

// RegisterHandlers initializes and registers all application handlers with the provided Gin router.
func RegisterHandlers(router *gin.Engine, services *Services) {
	api := router.Group("/api")

	// Initialize handlers
	cityHandler := NewCityHandler(services.CityService)
	facilityHandler := NewFacilityHandler(services.FacilityService)

	// Register routes
	cityHandler.RegisterCityRoutes(api)
	facilityHandler.RegisterFacilityRoutes(api)
}
