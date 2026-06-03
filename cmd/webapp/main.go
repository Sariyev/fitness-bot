package main

import (
	"context"
	"errors"
	"fitness-bot/internal/config"
	"fitness-bot/internal/database"
	webapphandler "fitness-bot/internal/handler/webapp"
	"fitness-bot/internal/payment"
	"fitness-bot/internal/repository"
	"fitness-bot/internal/service"
	"fitness-bot/internal/storage"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	if err := cfg.Validate(config.RoleWebApp); err != nil {
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
	moduleRepo := repository.NewModuleRepo(db.Pool)
	paymentRepo := repository.NewPaymentRepo(db.Pool)
	scoreRepo := repository.NewScoreRepo(db.Pool)
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

	// Payment provider: Robokassa if configured, otherwise Dummy
	var paymentProvider payment.Provider
	var verifier payment.CallbackVerifier
	if cfg.RobokassaMerchantLogin != "" {
		roboCfg := payment.RobokassaConfig{
			MerchantLogin: cfg.RobokassaMerchantLogin,
			Password1:     cfg.RobokassaPassword1,
			Password2:     cfg.RobokassaPassword2,
			IsTest:        cfg.RobokassaIsTest,
		}
		rp := payment.NewRobokassaProvider(roboCfg, paymentRepo)
		paymentProvider = rp
		verifier = rp
		log.Println("Payment provider: Robokassa (test:", cfg.RobokassaIsTest, ")")
	} else {
		paymentProvider = payment.NewDummyProvider()
		log.Println("Payment provider: Dummy (instant)")
	}

	// Per-category access (three-bucket content gating: free / trial / paid).
	// Constructed before paymentSvc so the payment service can read prices and
	// grant access on confirm.
	pricingRepo := repository.NewPricingRepo(db.Pool)
	accessRepo := repository.NewAccessRepo(db.Pool)
	accessSvc := service.NewAccessService(pricingRepo, accessRepo)

	scoreSvc := service.NewScoreService(scoreRepo)
	paymentSvc := service.NewPaymentService(paymentRepo, userRepo, accessSvc, paymentProvider, 5000)
	workoutSvc := service.NewWorkoutService(workoutRepo, exerciseRepo, completionRepo)
	rehabSvc := service.NewRehabService(rehabRepo)
	nutritionSvc := service.NewNutritionService(nutritionRepo, foodLogRepo)
	progressSvc := service.NewProgressService(progressRepo, completionRepo, achievementRepo)
	dashboardSvc := service.NewDashboardService(userSvc, workoutSvc, rehabSvc, nutritionSvc)
	recommendSvc := service.NewRecommendationService(workoutRepo, rehabRepo, nutritionRepo)

	// Media (R2) — optional. Empty AccessKeyID disables; routes won't register.
	var mediaSvc *service.MediaService
	if cfg.R2AccessKeyID != "" {
		r2Provider, err := storage.NewR2Provider(storage.R2Config{
			AccountID:       cfg.R2AccountID,
			AccessKeyID:     cfg.R2AccessKeyID,
			SecretAccessKey: cfg.R2SecretAccessKey,
			BucketPrivate:   cfg.R2BucketPrivate,
			BucketPublic:    cfg.R2BucketPublic,
			PublicURL:       cfg.R2PublicURL,
		})
		if err != nil {
			log.Fatalf("Failed to init R2 provider: %v", err)
		}
		mediaRepo := repository.NewMediaRepo(db.Pool)
		mediaSvc = service.NewMediaService(mediaRepo, r2Provider, cfg.MediaQuotaBytes)
		log.Printf("Media: R2 enabled (private=%s, public=%s, quota=%d bytes)",
			cfg.R2BucketPrivate, cfg.R2BucketPublic, cfg.MediaQuotaBytes)
	} else {
		log.Println("Media: R2 not configured — media endpoints disabled")
	}

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
		scoreSvc,
		mediaSvc,
		accessSvc,
		"./static",
		verifier,
		cfg.WebAppURL,
	)

	srv := &http.Server{
		Addr:              ":" + cfg.WebAppPort,
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
		log.Printf("WebApp server listening on :%s", cfg.WebAppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("WebApp server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down webapp...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("WebApp shutdown error: %v", err)
	}
}
