package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var (
	appCfg Config
	mu     sync.RWMutex
)

func LoadConfig() error {
	mu.Lock()
	defer mu.Unlock()

	// Load .env file jika ada (opsional, di production biasanya env sudah di-inject oleh orchestrator)
	if err := godotenv.Load(); err != nil {
		log.Println("[CONFIG] .env file not found, reading from system environment variables")
	}

	if err := env.Parse(&appCfg); err != nil {
		return err
	}

	log.Println("[CONFIG] Loaded. App:", appCfg.App.Name, "| ENV:", appCfg.App.Env)
	return nil
}

func GetConfig() Config {
	mu.RLock()
	defer mu.RUnlock()
	return appCfg
}
