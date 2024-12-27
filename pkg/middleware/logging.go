package middleware

import (
    "github.com/gin-gonic/gin"
    "log"
)

func Logging() gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Printf("%s %s", c.Request.Method, c.Request.URL.Path)
        c.Next()
    }
}
