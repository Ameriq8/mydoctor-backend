package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	r.GET("/users", GetUsers)
}

func GetUsers(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"users": []string{"User1", "User2"},
	})
}
