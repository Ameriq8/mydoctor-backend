package utils

import "github.com/gin-gonic/gin"

func StandardErrorResponse(message string, errorDetails []string) gin.H {
	return gin.H{
		"message":      message,
		"errorDetails": errorDetails,
	}
}
