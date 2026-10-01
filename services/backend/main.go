package main

import (
	"github.com/AdityaGudadhe/pair-coding/services/main/services/backend/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	routes.RegisterRoutes(router)
	router.Run(":3000")
}
