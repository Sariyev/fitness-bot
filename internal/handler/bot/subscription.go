package bot

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PaymentHandler struct {
	convSvc    *service.ConversationService
	paymentSvc *service.PaymentService
}

func NewPaymentHandler(convSvc *service.ConversationService, paymentSvc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{convSvc: convSvc, paymentSvc: paymentSvc}
}

// Step 1: /buy → Show product card
func (h *PaymentHandler) HandleBuy(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User) {
	if user.IsPaid {
		send(bot, chatID, "✅ У тебя уже есть полный доступ! Пользуйся на здоровье 💪\n\n/modules — Перейти к модулям")
		return
	}

	text := "🏋️ *Полный доступ к платформе*\n\n" +
		"Что входит:\n" +
		"🏥 ЛФК — упражнения при грыже, протрузиях, сколиозе и др.\n" +
		"💪 Тренировки — программы по группам мышц с видео\n" +
		"🥗 Питание — рецепты и планы питания\n\n" +
		"💰 *Стоимость: 5 000 ₸* (разовая оплата)"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Оплатить", "pay:start"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

// Step 2-4: Handle callback buttons
func (h *PaymentHandler) HandleCallback(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data string) {
	switch data {
	case "pay:start":
		h.showConfirmation(ctx, bot, chatID, user)
	case "pay:confirm":
		h.processPayment(ctx, bot, chatID, user)
	case "pay:cancel":
		h.convSvc.ClearState(ctx, user.TelegramID)
		send(bot, chatID, "Оплата отменена. Когда будешь готов — /buy")
	}
}

// Step 2: Confirmation screen
func (h *PaymentHandler) showConfirmation(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User) {
	if user.IsPaid {
		send(bot, chatID, "✅ У тебя уже есть полный доступ!")
		return
	}

	text := "📋 *Подтверждение оплаты*\n\n" +
		"Полный доступ к платформе\n" +
		"Сумма: *5 000 ₸*\n\n" +
		"Подтвердить оплату?"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить", "pay:confirm"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "pay:cancel"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

// Step 3-4: Process payment with loading animation
func (h *PaymentHandler) processPayment(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User) {
	if user.IsPaid {
		send(bot, chatID, "✅ У тебя уже есть полный доступ!")
		return
	}

	// Step 3: Show processing message
	processingMsg := tgbotapi.NewMessage(chatID, "⏳ Обработка платежа...")
	resp, err := bot.Send(processingMsg)
	if err != nil {
		send(bot, chatID, "Ошибка. Попробуй позже: /buy")
		return
	}

	// Simulate payment processing delay
	time.Sleep(2 * time.Second)

	// Process actual payment
	err = h.paymentSvc.ProcessPayment(ctx, user)
	if err != nil {
		log.Printf("Error processing payment for user %d: %v", user.ID, err)
		edit := tgbotapi.NewEditMessageText(chatID, resp.MessageID, "❌ Ошибка при оплате. Попробуй позже: /buy")
		bot.Send(edit)
		return
	}

	h.convSvc.ClearState(ctx, user.TelegramID)

	// Step 4: Edit processing message → success
	successText := "✅ Оплата прошла успешно!\n\n" +
		"Полный доступ ко всем модулям открыт. Приятных тренировок! 💪\n\n" +
		"/modules — Перейти к модулям"
	edit := tgbotapi.NewEditMessageText(chatID, resp.MessageID, successText)
	bot.Send(edit)
}

func (h *PaymentHandler) HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, state *models.ConversationState) {
	send(bot, msg.Chat.ID, "Используй кнопки для оплаты или /buy для начала.")
}
