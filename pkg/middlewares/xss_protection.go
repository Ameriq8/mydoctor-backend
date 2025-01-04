package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware() gin.HandlerFunc {
	policy := bluemonday.UGCPolicy() // Define a policy for safe HTML

	return func(c *gin.Context) {
		// Loop through all POST/PUT/QUERY parameters and sanitize them
		for key, values := range c.Request.Form {
			for i, value := range values {
				c.Request.Form[key][i] = policy.Sanitize(value) // Sanitize the input value
			}
		}

		c.Next()
	}
}
