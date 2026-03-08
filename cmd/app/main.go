package main

import (
	"log"
	"time"

	"github.com/IrsanaAhmad/go-starter-kit/internal/config"
	"github.com/IrsanaAhmad/go-starter-kit/internal/database"
	"github.com/IrsanaAhmad/go-starter-kit/internal/middleware"
	"github.com/IrsanaAhmad/go-starter-kit/internal/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Load konfigurasi dari .env / environment variables
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 2. Init DB LIMAU + pooling (modular, bisa diganti driver nanti)
	dbClient, err := database.InitDB()
	if err != nil {
		log.Fatal("DB Init failed:", err)
	}
	defer func() {
		if err := dbClient.Close(); err != nil {
			log.Println("DB close error:", err)
		}
	}()

	// 3. Setup Fiber app dengan nama app dari config
	app := fiber.New(fiber.Config{
		AppName: config.GetConfig().App.Name,
	})

	// 4. Global middleware
	app.Use(middleware.Logger())
	app.Use(middleware.RateLimiter(60, 1*time.Minute))

	router.Register(app, dbClient)

	// 5. Gunakan port dari .env
	port := config.GetConfig().App.Port
	log.Println("Server running on port", port)
	log.Fatal(app.Listen(":" + port))
}
