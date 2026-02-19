package main

import (
	"fitness-bot/internal/config"
	"fitness-bot/internal/database"
	adminhandler "fitness-bot/internal/handler/admin"
	"fitness-bot/internal/repository"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	// Database
	db, err := database.New(cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Migrations
	if err := database.RunMigrations(cfg.GetDatabaseURL()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepo(db.Pool)
	questRepo := repository.NewQuestionnaireRepo(db.Pool)
	moduleRepo := repository.NewModuleRepo(db.Pool)
	scoreRepo := repository.NewScoreRepo(db.Pool)

	// Admin router
	router := adminhandler.NewRouter(
		cfg.AdminAPIKey,
		moduleRepo,
		questRepo,
		userRepo,
		scoreRepo,
	)

	log.Printf("Admin API listening on :%s", cfg.AdminPort)
	if err := http.ListenAndServe(":"+cfg.AdminPort, router); err != nil {
		log.Fatalf("Admin server failed: %v", err)
	}
}
