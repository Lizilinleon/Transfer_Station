package main

import (
	"log"

	"ai-chat-platform/backend/config"
	"ai-chat-platform/backend/database"
	"ai-chat-platform/backend/router"
)

// main wires together configuration, database initialization, routes, and startup.
func main() {
	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	server := router.New(cfg, db)

	log.Printf("backend listening on :%s", cfg.Port)
	if err := server.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
