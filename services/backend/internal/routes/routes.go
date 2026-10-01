package routes

import (
	"github.com/AdityaGudadhe/pair-coding/services/main/services/backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/signup", handlers.SignupHandler)
	router.GET("/login", handlers.LoginHandler)
}
