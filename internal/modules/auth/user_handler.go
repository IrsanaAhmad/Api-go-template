package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/IrsanaAhmad/go-starter-kit/internal/config"
	"github.com/IrsanaAhmad/go-starter-kit/shared/response"

	"github.com/gofiber/fiber/v2"
)

const (
	CookieAccessToken  = "access_token"
	CookieRefreshToken = "refresh_token"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", nil)
	}

	if req.Username == "" || req.Email == "" || req.Password == "" || req.FullName == "" {
		return response.Error(c, http.StatusBadRequest, "semua field wajib diisi", nil)
	}

	if len(req.Password) < 8 {
		return response.Error(c, http.StatusBadRequest, "password minimal 8 karakter", nil)
	}

	ctx := context.Background()

	resp, err := h.service.Register(ctx, &req)
	if err != nil {
		if errors.Is(err, ErrConflict) {
			return response.Error(c, http.StatusConflict, "username atau email sudah terdaftar", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "gagal memproses registrasi", nil)
	}

	return response.Success(c, http.StatusCreated, "registrasi berhasil", resp)
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

	resp, tokens, err := h.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return response.Error(c, http.StatusUnauthorized, "username atau password salah", nil)
		}
		if errors.Is(err, ErrForbidden) {
			return response.Error(c, http.StatusForbidden, "akun tidak aktif", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "gagal memproses login", nil)
	}

	setTokenCookies(c, tokens)

	return response.Success(c, http.StatusOK, "login berhasil", resp)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Ambil refresh token dari cookie untuk revoke di DB
	refreshToken := c.Cookies(CookieRefreshToken)
	if refreshToken != "" {
		ctx := context.Background()
		_ = h.service.Logout(ctx, refreshToken)
	}

	clearTokenCookies(c)

	return response.Success(c, http.StatusOK, "logout berhasil", nil)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies(CookieRefreshToken)
	if refreshToken == "" {
		return response.Error(c, http.StatusUnauthorized, "refresh token not found", nil)
	}

	ctx := context.Background()

	resp, tokens, err := h.service.Refresh(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrUnauthorized) {
			clearTokenCookies(c)
			return response.Error(c, http.StatusUnauthorized, "refresh token tidak valid, silakan login ulang", nil)
		}
		if errors.Is(err, ErrForbidden) {
			clearTokenCookies(c)
			return response.Error(c, http.StatusForbidden, "akun tidak aktif", nil)
		}
		return response.Error(c, http.StatusInternalServerError, "gagal memproses refresh token", nil)
	}

	setTokenCookies(c, tokens)

	return response.Success(c, http.StatusOK, "token refreshed", resp)
}

// setTokenCookies menyimpan access & refresh token di HttpOnly cookies.
func setTokenCookies(c *fiber.Ctx, tokens *TokenPair) {
	cfg := config.GetConfig()
	isProduction := cfg.App.Env == "production"

	c.Cookie(&fiber.Cookie{
		Name:     CookieAccessToken,
		Value:    tokens.AccessToken,
		Path:     "/",
		MaxAge:   cfg.JWT.AccessTokenExp * 60, // menit → detik
		HTTPOnly: true,
		Secure:   isProduction,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     CookieRefreshToken,
		Value:    tokens.RefreshToken,
		Path:     "/api/v1/auth",                         // hanya dikirim ke auth endpoint
		MaxAge:   cfg.JWT.RefreshTokenExp * 24 * 60 * 60, // hari → detik
		HTTPOnly: true,
		Secure:   isProduction,
		SameSite: "Lax",
	})
}

// clearTokenCookies menghapus cookies token dengan set expired.
func clearTokenCookies(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     CookieAccessToken,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     CookieRefreshToken,
		Value:    "",
		Path:     "/api/v1/auth",
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})
}
