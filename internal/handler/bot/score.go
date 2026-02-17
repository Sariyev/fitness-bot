package bot

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const scoreTTL = 10 * time.Minute

type ScoreHandler struct {
	convSvc  *service.ConversationService
	scoreSvc *service.ScoreService
}

func NewScoreHandler(convSvc *service.ConversationService, scoreSvc *service.ScoreService) *ScoreHandler {
	return &ScoreHandler{convSvc: convSvc, scoreSvc: scoreSvc}
}

func (h *ScoreHandler) RequestScore(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, scoreType, refType string, refID int) {
	// Check if already scored
	scored, _ := h.scoreSvc.HasScored(ctx, user.ID, refType, refID)
	if scored {
		return
	}

	flowData := models.ScoreFlowData{
		ScoreType:     scoreType,
		ReferenceType: refType,
		ReferenceID:   refID,
	}
	h.convSvc.SetState(ctx, user.TelegramID, models.StateScoreAwaitRating, flowData, scoreTTL)

	var buttons []tgbotapi.InlineKeyboardButton
	for i := 1; i <= 10; i++ {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			strconv.Itoa(i),
			fmt.Sprintf("score:%d", i),
		))
	}
	row1 := buttons[:5]
	row2 := buttons[5:]
	keyboard := tgbotapi.NewInlineKeyboardMarkup(row1, row2)
	sendWithKeyboard(bot, chatID, "⭐ Оцените от 1 до 10:", keyboard)
}

func (h *ScoreHandler) HandleCallback(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data string) {
	if data == "score_skip" {
		h.saveScore(ctx, bot, chatID, user, "")
		return
	}

	if !strings.HasPrefix(data, "score:") {
		return
	}

	valStr := strings.TrimPrefix(data, "score:")
	val, err := strconv.Atoi(valStr)
	if err != nil || val < 1 || val > 10 {
		return
	}

	state, err := h.convSvc.GetState(ctx, user.TelegramID)
	if err != nil || state == nil {
		return
	}

	var flowData models.ScoreFlowData
	h.convSvc.GetData(state, &flowData)
	flowData.Score = val

	h.convSvc.SetState(ctx, user.TelegramID, models.StateScoreAwaitComment, flowData, scoreTTL)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пропустить", "score_skip"),
		),
	)
	sendWithKeyboard(bot, chatID, "💬 Хотите оставить комментарий? Напишите или нажмите «Пропустить»:", keyboard)
}

func (h *ScoreHandler) HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, state *models.ConversationState) {
	if state.State == models.StateScoreAwaitComment {
		h.saveScore(ctx, bot, msg.Chat.ID, user, msg.Text)
		return
	}

	if state.State == models.StateScoreAwaitRating {
		val, err := strconv.Atoi(msg.Text)
		if err != nil || val < 1 || val > 10 {
			send(bot, msg.Chat.ID, "Пожалуйста, используйте кнопки или введите число от 1 до 10:")
			return
		}

		var flowData models.ScoreFlowData
		h.convSvc.GetData(state, &flowData)
		flowData.Score = val
		h.convSvc.SetState(ctx, user.TelegramID, models.StateScoreAwaitComment, flowData, scoreTTL)

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Пропустить", "score_skip"),
			),
		)
		sendWithKeyboard(bot, msg.Chat.ID, "💬 Хотите оставить комментарий?", keyboard)
	}
}

func (h *ScoreHandler) saveScore(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, comment string) {
	state, err := h.convSvc.GetState(ctx, user.TelegramID)
	if err != nil || state == nil {
		return
	}

	var flowData models.ScoreFlowData
	h.convSvc.GetData(state, &flowData)

	if flowData.Score == 0 {
		flowData.Score = 5 // default
	}

	score := &models.UserScore{
		UserID:        user.ID,
		ScoreType:     flowData.ScoreType,
		ReferenceType: flowData.ReferenceType,
		ReferenceID:   flowData.ReferenceID,
		Score:         flowData.Score,
		Comment:       comment,
	}
	h.scoreSvc.SaveScore(ctx, score)
	h.convSvc.ClearState(ctx, user.TelegramID)

	send(bot, chatID, "Спасибо за оценку! 🙏\n\n/modules — Вернуться к модулям")
}
