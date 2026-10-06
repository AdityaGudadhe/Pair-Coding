package routes

import (
	"github.com/AdityaGudadhe/pair-coding/services/main/services/backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, auth *handlers.AuthHandler) {
	router.POST("/signup", auth.SignupHandler)
	router.POST("/login", auth.LoginHandler)
}
