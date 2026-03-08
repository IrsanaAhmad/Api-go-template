package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IrsanaAhmad/go-starter-kit/internal/config"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
)

type Client struct {
	db     *sql.DB
	driver string
}

// InitDB membuat koneksi database secara dinamis.
// Prioritas: DATABASE_URL > DB_* fields individual.
func InitDB() (*Client, error) {
	cfg := config.GetConfig()

	driver, dsn, err := resolveDSN(cfg.DatabaseURL, cfg.Database)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Connection pooling
	db.SetMaxIdleConns(cfg.Database.MaxIdleConn)
	db.SetMaxOpenConns(cfg.Database.MaxOpenConn)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.MaxLifetime) * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	log.Printf("[DB] Connected via %s driver\n", driver)
	return &Client{db: db, driver: driver}, nil
}

// resolveDSN menentukan driver dan DSN berdasarkan konfigurasi.
func resolveDSN(databaseURL string, dbCfg config.DatabaseConfig) (driver string, dsn string, err error) {
	// Prioritas 1: DATABASE_URL langsung
	if databaseURL != "" {
		driver = detectDriverFromURL(databaseURL)
		return driver, databaseURL, nil
	}

	// Prioritas 2: Bangun DSN dari field individual
	switch strings.ToLower(dbCfg.Connection) {
	case "postgres", "postgresql":
		driver = "postgres"
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
			dbCfg.Host, dbCfg.Port, dbCfg.Username, dbCfg.Password, dbCfg.Name,
		)
	case "sqlserver", "mssql":
		driver = "sqlserver"
		dsn = fmt.Sprintf(
			"server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=%s;trustServerCertificate=%t",
			dbCfg.Host, dbCfg.Username, dbCfg.Password, dbCfg.Port, dbCfg.Name,
			dbCfg.Encrypt, dbCfg.TrustServerCertificate,
		)
	default:
		return "", "", fmt.Errorf("unsupported DB_CONNECTION: %s", dbCfg.Connection)
	}

	return driver, dsn, nil
}

// detectDriverFromURL menebak driver dari prefix URL.
func detectDriverFromURL(url string) string {
	switch {
	case strings.HasPrefix(url, "postgres://"), strings.HasPrefix(url, "postgresql://"):
		return "postgres"
	case strings.HasPrefix(url, "sqlserver://"):
		return "sqlserver"
	default:
		return "postgres" // default fallback
	}
}

func (c *Client) GetDB() *sql.DB {
	return c.db
}

func (c *Client) GetDriver() string {
	return c.driver
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
