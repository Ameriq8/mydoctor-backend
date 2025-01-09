package middlewares

import (
	"net/http"
	"server/internal/services"
	"server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks if the user has a valid authorization token
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from headers or query params
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Validate the token (you can implement more specific logic here)
		_, err := authService.ValidateAuthToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Proceed to the next handler if valid
		c.Next()
	}
}

// RoleBasedAccessControl checks the user's role based on a required role
func RoleBasedAccessControl(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Assume `role` is stored in JWT claims or session
		userRole := c.GetHeader("X-User-Role")

		if userRole != requiredRole {
			// Pass an empty slice for errorDetails
			c.JSON(http.StatusForbidden, utils.StandardErrorResponse("Forbidden: Insufficient permissions", []string{}))
			c.Abort()
			return
		}

		c.Next()
	}
}
