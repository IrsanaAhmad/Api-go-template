package auth

import (
	"time"

	"github.com/IrsanaAhmad/go-starter-kit/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID string, username string) (string, error) {
	cfg := config.GetConfig()
	secret := []byte(cfg.JWT.SecretKey)
	expiresAt := time.Now().Add(time.Duration(cfg.JWT.AccessTokenExp) * time.Minute)

	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GenerateRefreshToken(userID string) (string, error) {
	cfg := config.GetConfig()
	secret := []byte(cfg.JWT.SecretKey)
	expiresAt := time.Now().Add(time.Duration(cfg.JWT.RefreshTokenExp) * 24 * time.Hour)

	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenStr string) (*CustomClaims, error) {
	cfg := config.GetConfig()
	secret := []byte(cfg.JWT.SecretKey)

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
