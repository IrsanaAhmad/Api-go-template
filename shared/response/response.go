package response

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	resp := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	return c.Status(statusCode).JSON(resp)
}

func Error(c *fiber.Ctx, statusCode int, message string, errors interface{}) error {
	resp := APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}
	return c.Status(statusCode).JSON(resp)
}
