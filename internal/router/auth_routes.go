package router

import (
	"github.com/IrsanaAhmad/go-starter-kit/internal/database"
	"github.com/IrsanaAhmad/go-starter-kit/internal/modules/auth"
	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(v1 fiber.Router, db database.DBClient) {
	auth.RegisterRoutes(v1, db)
}
