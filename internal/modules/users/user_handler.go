package users

import (
	"context"
	"net/http"

	"github.com/IrsanaAhmad/go-starter-kit/internal/auth"
	"github.com/IrsanaAhmad/go-starter-kit/shared/response"
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Login(ctx context.Context, username string, password string) (*LoginResponse, error)
}

type authService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Login(ctx context.Context, username string, password string) (*LoginResponse, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fiber.ErrUnauthorized
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fiber.ErrUnauthorized
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	resp := &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserDTO{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
		},
	}

	return resp, nil
}

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
