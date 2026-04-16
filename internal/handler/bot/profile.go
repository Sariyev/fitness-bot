package bot

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ProfileHandler struct {
	userSvc *service.UserService
}

func NewProfileHandler(userSvc *service.UserService) *ProfileHandler {
	return &ProfileHandler{userSvc: userSvc}
}

func (h *ProfileHandler) HandleProfile(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User) {
	profile, err := h.userSvc.GetProfile(ctx, user.ID)
	if err != nil {
		send(bot, chatID, "Профиль не найден. Пройди регистрацию: /start")
		return
	}

	paymentStatus := "Не оплачено"
	if user.IsPaid {
		paymentStatus = "✅ Оплачено"
	}

	// Parse comma-separated goals
	goals := strings.Split(profile.Goal, ",")
	goalsText := goalsLabel(goals)

	text := fmt.Sprintf(
		"👤 *Твой профиль*\n\n"+
			"Имя: %s %s\n"+
			"Возраст: %d\n"+
			"Рост: %d см\n"+
			"Вес: %.1f кг\n"+
			"Пол: %s\n"+
			"Уровень: %s\n"+
			"Цели: %s\n\n"+
			"💳 Доступ: %s",
		user.FirstName, user.LastName,
		profile.Age, profile.HeightCm, profile.WeightKg,
		genderLabel(profile.Gender),
		fitnessLabel(profile.FitnessLevel),
		goalsText,
		paymentStatus,
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func (h *ProfileHandler) HandleHelp(bot *tgbotapi.BotAPI, chatID int64) {
	text := "📖 *Доступные команды:*\n\n" +
		"/start — Начало / главное меню\n" +
		"/modules — Модули и тренировки\n" +
		"/app — Открыть приложение\n" +
		"/buy — Оплатить доступ\n" +
		"/profile — Мой профиль\n" +
		"/review — Оставить отзыв о боте\n" +
		"/help — Эта справка"

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
