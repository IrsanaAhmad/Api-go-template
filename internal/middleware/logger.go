package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Logger mengembalikan middleware yang mencatat setiap HTTP request
// dengan informasi: method, path, status code, latency, dan client IP.
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Lanjutkan ke handler berikutnya
		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.OriginalURL()
		ip := c.IP()

		// Format log: [HTTP] 200 | 1.23ms | 127.0.0.1 | GET /api/v1/auth/login
		log.Printf("[HTTP] %d | %13v | %15s | %-7s %s\n",
			status, latency, ip, method, path,
		)

		return err
	}
}
