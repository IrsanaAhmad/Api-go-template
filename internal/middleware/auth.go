package middleware

import (
	"net/http"
	"strings"

	"github.com/IrsanaAhmad/go-starter-kit/internal/auth"
	"github.com/IrsanaAhmad/go-starter-kit/shared/response"
	"github.com/gofiber/fiber/v2"
)

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, http.StatusUnauthorized, "missing authorization header", nil)
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return response.Error(c, http.StatusUnauthorized, "invalid authorization header format", nil)
		}

		tokenStr := parts[1]

		claims, err := auth.ParseToken(tokenStr)
		if err != nil {
			return response.Error(c, http.StatusUnauthorized, "invalid or expired token", nil)
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}
