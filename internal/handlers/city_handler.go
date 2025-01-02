package handlers

import (
	"net/http"
	"strconv"

	"server/internal/services"

	"github.com/gin-gonic/gin"
)

type CityHandler struct {
	service *services.CityService
}

// NewCityHandler creates a new CityHandler.
func NewCityHandler(service *services.CityService) *CityHandler {
	return &CityHandler{service: service}
}

// RegisterCityRoutes registers city-related routes.
func (h *CityHandler) RegisterCityRoutes(r *gin.RouterGroup) {
	r.GET("/cities/:id", h.GetCityByID)
}

// GetCityByID handles the GET request for fetching a city by ID.
func (h *CityHandler) GetCityByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	city, err := h.service.GetCityByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"city": city})
}
