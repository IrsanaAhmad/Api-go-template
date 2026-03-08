package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/IrsanaAhmad/go-starter-kit/internal/database"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	StoreRefreshToken(ctx context.Context, token *RefreshToken) error
	RevokeRefreshTokenByHash(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
	IsRefreshTokenValid(ctx context.Context, tokenHash string) (bool, error)
}

type SQLUserRepository struct {
	db database.DBClient
}

func NewSQLUserRepository(db database.DBClient) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, username, email, password_hash, full_name, role, is_active, created_at, updated_at
		FROM users WHERE username = $1`
	row := r.db.GetDB().QueryRowContext(ctx, query, username)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *SQLUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, username, email, password_hash, full_name, role, is_active, created_at, updated_at
		FROM users WHERE email = $1`
	row := r.db.GetDB().QueryRowContext(ctx, query, email)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *SQLUserRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	query := `INSERT INTO users (username, email, password_hash, full_name, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, email, full_name, role, is_active, created_at, updated_at`

	row := r.db.GetDB().QueryRowContext(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.FullName, user.Role,
	)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *SQLUserRepository) StoreRefreshToken(ctx context.Context, token *RefreshToken) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)`

	_, err := r.db.GetDB().ExecContext(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt)
	return err
}

func (r *SQLUserRepository) RevokeRefreshTokenByHash(ctx context.Context, tokenHash string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1 AND revoked_at IS NULL`

	result, err := r.db.GetDB().ExecContext(ctx, query, tokenHash)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("refresh token not found or already revoked")
	}

	return nil
}

func (r *SQLUserRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`

	_, err := r.db.GetDB().ExecContext(ctx, query, userID)
	return err
}

func (r *SQLUserRepository) FindByID(ctx context.Context, id string) (*User, error) {
	query := `SELECT id, username, email, password_hash, full_name, role, is_active, created_at, updated_at
		FROM users WHERE id = $1`
	row := r.db.GetDB().QueryRowContext(ctx, query, id)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *SQLUserRepository) IsRefreshTokenValid(ctx context.Context, tokenHash string) (bool, error) {
	query := `SELECT COUNT(1) FROM refresh_tokens
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()`

	var count int
	if err := r.db.GetDB().QueryRowContext(ctx, query, tokenHash).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}
