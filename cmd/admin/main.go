package main

import (
	"context"
	"errors"
	"fitness-bot/internal/config"
	"fitness-bot/internal/database"
	adminhandler "fitness-bot/internal/handler/admin"
	"fitness-bot/internal/repository"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	if err := cfg.Validate(config.RoleAdmin); err != nil {
		log.Fatalf("config: %v", err)
	}

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

	srv := &http.Server{
		Addr:              ":" + cfg.AdminPort,
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Admin API listening on :%s", cfg.AdminPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Admin server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down admin...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Admin shutdown error: %v", err)
	}
}
