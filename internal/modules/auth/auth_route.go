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
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/logout", authHandler.Logout)
	authGroup.Post("/refresh", authHandler.Refresh)
}
