package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IrsanaAhmad/go-starter-kit/internal/config"
	_ "github.com/microsoft/go-mssqldb"
)

type Client struct {
	db *sql.DB
}

// InitDB koneksi ke database LIMAU untuk login
func InitDB() (*Client, error) {
	cfg := config.GetConfig()

	targetDB := cfg.Database.DBName.LIMAU // ← LIMAU dipakai untuk login

	dsn := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=%s;trustServerCertificate=%t",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Port,
		targetDB,
		cfg.Database.Encrypt,
		cfg.Database.TrustServerCertificate,
	)

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}

	// Pooling sesuai config.json
	db.SetMaxIdleConns(cfg.Database.MaxIdleConn)
	db.SetMaxOpenConns(cfg.Database.MaxOpenConn)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.MaxLifetime) * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, errors.New("database unreachable")
	}

	log.Println("Connected to database:", targetDB)
	return &Client{db: db}, nil
}

func (c *Client) GetDB() *sql.DB {
	return c.db
}

func (c *Client) Close() error {
	if c.db == nil {
		return errors.New("db nil")
	}
	return c.db.Close()
}

func (c *Client) HealthCheck() error {
	if c.db == nil {
		return errors.New("db nil")
	}
	return c.db.Ping()
}
