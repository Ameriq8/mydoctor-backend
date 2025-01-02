package middlewares

import (
	"net/http"
	"server/pkg/utils"

	"github.com/gin-gonic/gin"
)

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
