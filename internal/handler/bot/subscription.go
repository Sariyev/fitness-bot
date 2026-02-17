package bot

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SubscriptionHandler struct {
	convSvc *service.ConversationService
	subSvc  *service.SubscriptionService
}

func NewSubscriptionHandler(convSvc *service.ConversationService, subSvc *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{convSvc: convSvc, subSvc: subSvc}
}

func (h *SubscriptionHandler) HandleShowPlans(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User) {
	// Check if already subscribed
	active, _ := h.subSvc.HasActiveSubscription(ctx, user.ID)
	if active {
		sub, _ := h.subSvc.GetActiveSubscription(ctx, user.ID)
		if sub != nil {
			send(bot, chatID, fmt.Sprintf(
				"✅ У вас активная подписка!\n\nДействует до: %s\n\n/modules — Перейти к модулям",
				sub.ExpiresAt.Format("02.01.2006"),
			))
			return
		}
	}

	plans, err := h.subSvc.ListPlans(ctx)
	if err != nil || len(plans) == 0 {
		send(bot, chatID, "Планы подписки временно недоступны.")
		return
	}

	text := "💳 Выберите план подписки:\n\n"
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, p := range plans {
		text += fmt.Sprintf("• %s — %d ₸ (%d дней)\n  %s\n\n", p.Name, p.PriceKZT, p.DurationDays, p.Description)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s — %d ₸", p.Name, p.PriceKZT),
				fmt.Sprintf("plan:%d", p.ID),
			),
		))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	sendWithKeyboard(bot, chatID, text, keyboard)
}

func (h *SubscriptionHandler) HandleCallback(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data string) {
	if !strings.HasPrefix(data, "plan:") {
		return
	}

	planIDStr := strings.TrimPrefix(data, "plan:")
	planID, err := strconv.Atoi(planIDStr)
	if err != nil {
		return
	}

	sub, err := h.subSvc.Subscribe(ctx, user.ID, planID)
	if err != nil {
		log.Printf("Error subscribing user %d: %v", user.ID, err)
		send(bot, chatID, "Произошла ошибка при оформлении подписки. Попробуйте позже.")
		return
	}

	h.convSvc.ClearState(ctx, user.TelegramID)

	send(bot, chatID, fmt.Sprintf(
		"✅ Подписка оформлена!\n\n"+
			"Действует до: %s\n\n"+
			"/modules — Перейти к модулям",
		sub.ExpiresAt.Format("02.01.2006"),
	))
}

func (h *SubscriptionHandler) HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, state *models.ConversationState) {
	_ = time.Now()
	send(bot, msg.Chat.ID, "Пожалуйста, используйте кнопки для выбора плана подписки.")
}
