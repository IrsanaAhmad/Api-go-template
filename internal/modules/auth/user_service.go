package auth

import (
	"context"

	"github.com/IrsanaAhmad/go-starter-kit/internal/auth"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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
