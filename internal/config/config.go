package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

var (
	AppConfig Config
	mu        sync.RWMutex
)

func LoadConfig() error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.ReadFile("internal/config/config.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, &AppConfig); err != nil {
		return err
	}

	log.Println("[CONFIG] Loaded. App:", AppConfig.App.Name, "| ENV:", AppConfig.App.Env)
	return nil
}

func GetConfig() Config {
	mu.RLock()
	defer mu.RUnlock()
	return AppConfig
}

func WatchConfigChanges() {
	var lastMod time.Time

	for {
		info, err := os.Stat("internal/config/config.json")
		if err != nil {
			log.Println("[CONFIG] Watch error:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if info.ModTime().After(lastMod) {
			log.Println("[CONFIG] Change detected → Reloading...")
			if err := LoadConfig(); err != nil {
				log.Println("[CONFIG] Reload failed:", err)
			} else {
				lastMod = info.ModTime()
				log.Println("[CONFIG] Reload success!")
			}
		}

		time.Sleep(2 * time.Second)
	}
}
