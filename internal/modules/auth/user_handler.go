package auth

import (
	"context"
	"net/http"

	"github.com/IrsanaAhmad/go-starter-kit/shared/response"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", nil)
	}

	if req.Username == "" || req.Password == "" {
		return response.Error(c, http.StatusBadRequest, "username dan password wajib diisi", nil)
	}

	ctx := context.Background()

	resp, err := h.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		if err == fiber.ErrUnauthorized {
			return response.Error(c, http.StatusUnauthorized, "username atau password salah", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "gagal memproses login", nil)
	}

	return response.Success(c, http.StatusOK, "login berhasil", resp)
}
