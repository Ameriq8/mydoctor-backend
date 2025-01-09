package validators

import (
	"server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// FacilityRequest represents the structure of the request body for creating or updating a facility.
type FacilityRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Type        string `json:"type" validate:"required"`
	CityID      int    `json:"city_id" validate:"required"`
}

// ValidateFacilityRequest validates the request body for creating or updating a facility.
func ValidateFacilityRequest(c *gin.Context, facilityRequest FacilityRequest) {
	// Use the ValidateRequest middleware to validate the request
	utils.ValidateRequest(c, facilityRequest)
}
