package auth

import "time"

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	FullName     string
	Role         string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time // nil jika belum di-revoke
	CreatedAt time.Time
}
