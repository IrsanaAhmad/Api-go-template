package router

import (
	"github.com/IrsanaAhmad/go-starter-kit/internal/database"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, db database.DBClient) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	RegisterAuthRoutes(v1, db)
}
