package middlewares

import (
	"server/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logging() gin.HandlerFunc {
	logger.InitLogger("development")
	defer logger.Sync()

	return func(c *gin.Context) {
		logger.Info("%s %s", zap.Any(c.Request.Method, c.Request.URL.Path))
		c.Next()
	}
}
