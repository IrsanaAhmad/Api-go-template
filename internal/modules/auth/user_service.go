package auth

import (
	"context"
	"errors"
	"time"

	"github.com/IrsanaAhmad/go-starter-kit/internal/auth"
	"github.com/IrsanaAhmad/go-starter-kit/internal/config"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("username or email already exists")
	ErrInvalidToken = errors.New("invalid or expired refresh token")
)

// TokenPair menyimpan access dan refresh token (digunakan internal, tidak di-expose ke JSON)
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, username string, password string) (*LoginResponse, *TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (*RefreshResponse, *TokenPair, error)
}

type authService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Cek apakah username sudah ada
	existing, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrConflict
	}

	// Cek apakah email sudah ada
	existingEmail, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, ErrConflict
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Role:         "user",
	}

	created, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		User: UserDTO{
			ID:       created.ID,
			Username: created.Username,
			Email:    created.Email,
			FullName: created.FullName,
			Role:     created.Role,
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, username string, password string) (*LoginResponse, *TokenPair, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, ErrUnauthorized
	}

	if !user.IsActive {
		return nil, nil, ErrForbidden
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, ErrUnauthorized
	}

	tokens, err := s.generateAndStoreTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	resp := &LoginResponse{
		User: UserDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}

	return resp, tokens, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := auth.HashToken(refreshToken)
	return s.repo.RevokeRefreshTokenByHash(ctx, tokenHash)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*RefreshResponse, *TokenPair, error) {
	// 1. Validasi JWT refresh token
	claims, err := auth.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, nil, ErrInvalidToken
	}

	userID := claims.Subject

	// 2. Cek apakah token ada di DB dan belum di-revoke
	tokenHash := auth.HashToken(refreshToken)
	valid, err := s.repo.IsRefreshTokenValid(ctx, tokenHash)
	if err != nil || !valid {
		return nil, nil, ErrInvalidToken
	}

	// 3. Revoke token lama
	_ = s.repo.RevokeRefreshTokenByHash(ctx, tokenHash)

	// 4. Ambil data user terbaru
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, nil, ErrUnauthorized
	}

	if !user.IsActive {
		return nil, nil, ErrForbidden
	}

	// 5. Generate token baru (rotation)
	tokens, err := s.generateAndStoreTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	resp := &RefreshResponse{
		User: UserDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}

	return resp, tokens, nil
}

// generateAndStoreTokens membuat access + refresh token pair baru dan menyimpan hash refresh token di DB.
func (s *authService) generateAndStoreTokens(ctx context.Context, user *User) (*TokenPair, error) {
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	cfg := config.GetConfig()
	tokenHash := auth.HashToken(refreshToken)
	expiresAt := time.Now().Add(time.Duration(cfg.JWT.RefreshTokenExp) * 24 * time.Hour)

	rt := &RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	if err := s.repo.StoreRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
