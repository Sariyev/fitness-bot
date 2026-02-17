package bot

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ProfileHandler struct {
	userSvc *service.UserService
	subSvc  *service.SubscriptionService
}

func NewProfileHandler(userSvc *service.UserService, subSvc *service.SubscriptionService) *ProfileHandler {
	return &ProfileHandler{userSvc: userSvc, subSvc: subSvc}
}

func (h *ProfileHandler) HandleProfile(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User) {
	profile, err := h.userSvc.GetProfile(ctx, user.ID)
	if err != nil {
		send(bot, chatID, "Профиль не найден. Пройдите регистрацию: /start")
		return
	}

	subStatus := "Нет активной подписки"
	sub, _ := h.subSvc.GetActiveSubscription(ctx, user.ID)
	if sub != nil {
		subStatus = fmt.Sprintf("Активна до %s", sub.ExpiresAt.Format("02.01.2006"))
	}

	text := fmt.Sprintf(
		"👤 *Ваш профиль*\n\n"+
			"Имя: %s %s\n"+
			"Возраст: %d\n"+
			"Рост: %d см\n"+
			"Вес: %.1f кг\n"+
			"Пол: %s\n"+
			"Уровень: %s\n"+
			"Цель: %s\n\n"+
			"📋 Подписка: %s",
		user.FirstName, user.LastName,
		profile.Age, profile.HeightCm, profile.WeightKg,
		genderLabel(profile.Gender),
		fitnessLabel(profile.FitnessLevel),
		goalLabel(profile.Goal),
		subStatus,
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func (h *ProfileHandler) HandleHelp(bot *tgbotapi.BotAPI, chatID int64) {
	text := "📖 *Доступные команды:*\n\n" +
		"/start — Начало / главное меню\n" +
		"/modules — Модули обучения\n" +
		"/app — Открыть приложение\n" +
		"/subscribe — Оформить подписку\n" +
		"/profile — Мой профиль\n" +
		"/help — Эта справка"

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
