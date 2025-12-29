package router

import (
	"github.com/IrsanaAhmad/go-starter-kit/internal/database"
	"github.com/IrsanaAhmad/go-starter-kit/internal/modules/users"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, db database.DBClient) {
	authRepo := users.NewSQLUserRepository(db)
	authService := users.NewAuthService(authRepo)
	authHandler := users.NewAuthHandler(authService)

	api := app.Group("/api")
	v1 := api.Group("/v1")
	authGroup := v1.Group("/auth")

	authGroup.Post("/login", authHandler.Login)
}
