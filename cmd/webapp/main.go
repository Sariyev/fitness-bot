package main

import (
	"fitness-bot/internal/config"
	"fitness-bot/internal/database"
	webapphandler "fitness-bot/internal/handler/webapp"
	"fitness-bot/internal/payment"
	"fitness-bot/internal/repository"
	"fitness-bot/internal/service"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
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
	moduleRepo := repository.NewModuleRepo(db.Pool)
	paymentRepo := repository.NewPaymentRepo(db.Pool)
	programRepo := repository.NewProgramRepo(db.Pool)
	workoutRepo := repository.NewWorkoutRepo(db.Pool)
	exerciseRepo := repository.NewExerciseRepo(db.Pool)
	completionRepo := repository.NewCompletionRepo(db.Pool)
	rehabRepo := repository.NewRehabRepo(db.Pool)
	nutritionRepo := repository.NewNutritionRepo(db.Pool)
	foodLogRepo := repository.NewFoodLogRepo(db.Pool)
	progressRepo := repository.NewProgressRepo(db.Pool)
	achievementRepo := repository.NewAchievementRepo(db.Pool)

	// Services
	userSvc := service.NewUserService(userRepo)
	moduleSvc := service.NewModuleService(moduleRepo)
	paymentProvider := payment.NewDummyProvider()
	paymentSvc := service.NewPaymentService(paymentRepo, userRepo, paymentProvider, 5000)
	workoutSvc := service.NewWorkoutService(programRepo, workoutRepo, exerciseRepo, completionRepo)
	rehabSvc := service.NewRehabService(rehabRepo)
	nutritionSvc := service.NewNutritionService(nutritionRepo, foodLogRepo)
	progressSvc := service.NewProgressService(progressRepo, completionRepo, achievementRepo)
	dashboardSvc := service.NewDashboardService(userSvc, workoutSvc, rehabSvc, nutritionSvc)
	recommendSvc := service.NewRecommendationService(programRepo, rehabRepo, nutritionRepo)

	// WebApp router
	router := webapphandler.NewRouter(
		cfg.TelegramToken,
		userSvc,
		moduleSvc,
		paymentSvc,
		workoutSvc,
		rehabSvc,
		nutritionSvc,
		progressSvc,
		dashboardSvc,
		recommendSvc,
		"./static",
	)

	log.Printf("WebApp server listening on :%s", cfg.WebAppPort)
	if err := http.ListenAndServe(":"+cfg.WebAppPort, router); err != nil {
		log.Fatalf("WebApp server failed: %v", err)
	}
}
