package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/IrsanaAhmad/go-starter-kit/internal/database"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
}

type SQLUserRepository struct {
	db database.DBClient
}

func NewSQLUserRepository(db database.DBClient) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	query := "SELECT id, username, password_hash, full_name FROM users WHERE username = @p1"
	row := r.db.GetDB().QueryRowContext(ctx, query, username)

	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}
