package middlewares

import (
	"net/http"
	"server/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateRequest is a middleware that validates request payloads using Go's validator package.
func ValidateRequest(c *gin.Context, structData interface{}) {
	if err := validate.Struct(structData); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, err.Error()) // Reassign the result to errorMessages
		}

		// Use the errorMessages in the response
		c.JSON(http.StatusBadRequest, utils.StandardErrorResponse("Invalid request data", errorMessages))
		c.Abort()
		return
	}
	c.Next()
}
