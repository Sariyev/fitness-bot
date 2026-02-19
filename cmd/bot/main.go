package main

import (
	"context"
	"fitness-bot/internal/config"
	"fitness-bot/internal/database"
	bothandler "fitness-bot/internal/handler/bot"
	"fitness-bot/internal/payment"
	"fitness-bot/internal/repository"
	"fitness-bot/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	convRepo := repository.NewConversationRepo(db.Pool)
	questRepo := repository.NewQuestionnaireRepo(db.Pool)
	paymentRepo := repository.NewPaymentRepo(db.Pool)
	moduleRepo := repository.NewModuleRepo(db.Pool)
	scoreRepo := repository.NewScoreRepo(db.Pool)

	// Payment provider
	paymentProvider := payment.NewDummyProvider()

	// Services
	userSvc := service.NewUserService(userRepo)
	convSvc := service.NewConversationService(convRepo)
	questSvc := service.NewQuestionnaireService(questRepo)
	paymentSvc := service.NewPaymentService(paymentRepo, userRepo, paymentProvider, 5000)
	moduleSvc := service.NewModuleService(moduleRepo)
	scoreSvc := service.NewScoreService(scoreRepo)

	// Bot
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	bot.Debug = cfg.Debug

	log.Printf("Bot started: @%s", bot.Self.UserName)

	// Router
	router := bothandler.NewRouter(userSvc, convSvc, questSvc, paymentSvc, moduleSvc, scoreSvc, cfg.WebAppURL)

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Shutting down...")
		cancel()
	}()

	// Start polling
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			bot.StopReceivingUpdates()
			log.Println("Bot stopped")
			return
		case update := <-updates:
			go router.HandleUpdate(bot, update)
		}
	}
}
