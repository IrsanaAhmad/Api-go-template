package auth

import (
	"github.com/IrsanaAhmad/go-starter-kit/internal/database"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, db database.DBClient) {
	authRepo := NewSQLUserRepository(db)
	authService := NewAuthService(authRepo)
	authHandler := NewAuthHandler(authService)

	authGroup := router.Group("/auth")
	authGroup.Post("/login", authHandler.Login)
}
