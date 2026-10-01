package main

import (
	"context"
	"log"
	"os"

	"github.com/AdityaGudadhe/pair-coding/services/main/services/backend/internal"
	"github.com/AdityaGudadhe/pair-coding/services/main/services/backend/internal/handlers"
	"github.com/AdityaGudadhe/pair-coding/services/main/services/backend/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// minSecretLen is the shortest JWT signing key accepted: HS256 needs at
// least 32 bytes of key to be as strong as its hash.
const minSecretLen = 32

func main() {
	// A missing .env is fine; the variables may come from the environment.
	_ = godotenv.Load()

	secret := os.Getenv("SECRET_KEY")
	if len(secret) < minSecretLen {
		log.Fatalf("SECRET_KEY must be set and at least %d characters long", minSecretLen)
	}

	pool, err := internal.InitDatabase(context.Background())
	if err != nil {
		log.Fatalf("init database: %v", err)
	}
	defer pool.Close()

	router := gin.Default()
	routes.RegisterRoutes(router, handlers.NewAuthHandler(pool, []byte(secret)))
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
