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

// RequestScore starts the review flow for specific content.
func (h *ScoreHandler) RequestScore(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, scoreType, refType string, refID int) {
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
	h.sendRatingKeyboard(bot, chatID, "⭐ Оцените от 1 до 5:")
}

// StartBotReview starts the review flow for the bot itself (/review command).
func (h *ScoreHandler) StartBotReview(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User) {
	flowData := models.ScoreFlowData{
		ScoreType:     "rating",
		ReferenceType: "bot",
		ReferenceID:   0,
	}
	h.convSvc.SetState(ctx, user.TelegramID, models.StateScoreAwaitRating, flowData, scoreTTL)
	h.sendRatingKeyboard(bot, chatID, "Мы ценим ваше мнение! Оцените бота от 1 до 5:")
}

func (h *ScoreHandler) sendRatingKeyboard(bot *tgbotapi.BotAPI, chatID int64, text string) {
	stars := []string{"1", "2", "3", "4", "5"}
	var buttons []tgbotapi.InlineKeyboardButton
	for _, s := range stars {
		n, _ := strconv.Atoi(s)
		label := strings.Repeat("⭐", n)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(label, "score:"+s))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
	sendWithKeyboard(bot, chatID, text, keyboard)
}

func (h *ScoreHandler) HandleCallback(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data string) {
	state, err := h.convSvc.GetState(ctx, user.TelegramID)
	if err != nil || state == nil {
		return
	}

	var flowData models.ScoreFlowData
	h.convSvc.GetData(state, &flowData)

	switch {
	// Rating selection
	case strings.HasPrefix(data, "score:") && state.State == models.StateScoreAwaitRating:
		valStr := strings.TrimPrefix(data, "score:")
		val, err := strconv.Atoi(valStr)
		if err != nil || val < 1 || val > 5 {
			return
		}
		flowData.Score = val
		h.convSvc.SetState(ctx, user.TelegramID, models.StateScoreAwaitTags, flowData, scoreTTL)
		h.sendTagKeyboard(bot, chatID, flowData.ReferenceType, flowData.Tags)

	// Tag toggle
	case strings.HasPrefix(data, "tag:") && state.State == models.StateScoreAwaitTags:
		tag := strings.TrimPrefix(data, "tag:")
		flowData.Tags = toggleTag(flowData.Tags, tag)
		h.convSvc.SetState(ctx, user.TelegramID, models.StateScoreAwaitTags, flowData, scoreTTL)
		h.sendTagKeyboard(bot, chatID, flowData.ReferenceType, flowData.Tags)

	// Done selecting tags
	case data == "tags_done" && state.State == models.StateScoreAwaitTags:
		h.convSvc.SetState(ctx, user.TelegramID, models.StateScoreAwaitComment, flowData, scoreTTL)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Пропустить", "score_skip"),
			),
		)
		sendWithKeyboard(bot, chatID, "💬 Хотите оставить комментарий? Напишите или нажмите «Пропустить»:", keyboard)

	// Skip comment
	case data == "score_skip" && state.State == models.StateScoreAwaitComment:
		h.saveScore(ctx, bot, chatID, user, "")
	}
}

func (h *ScoreHandler) HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, state *models.ConversationState) {
	switch state.State {
	case models.StateScoreAwaitRating:
		val, err := strconv.Atoi(msg.Text)
		if err != nil || val < 1 || val > 5 {
			send(bot, msg.Chat.ID, "Пожалуйста, используйте кнопки или введите число от 1 до 5:")
			return
		}
		var flowData models.ScoreFlowData
		h.convSvc.GetData(state, &flowData)
		flowData.Score = val
		h.convSvc.SetState(ctx, user.TelegramID, models.StateScoreAwaitTags, flowData, scoreTTL)
		h.sendTagKeyboard(bot, msg.Chat.ID, flowData.ReferenceType, flowData.Tags)

	case models.StateScoreAwaitComment:
		h.saveScore(ctx, bot, msg.Chat.ID, user, msg.Text)
	}
}

func (h *ScoreHandler) sendTagKeyboard(bot *tgbotapi.BotAPI, chatID int64, refType string, selected []string) {
	tags, ok := models.ReviewTags[refType]
	if !ok {
		tags = models.ReviewTags["bot"]
	}

	selectedSet := make(map[string]bool)
	for _, t := range selected {
		selectedSet[t] = true
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	for i, tag := range tags {
		label := tag
		if selectedSet[tag] {
			label = "✅ " + tag
		}
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(label, "tag:"+tag))
		if (i+1)%2 == 0 || i == len(tags)-1 {
			rows = append(rows, row)
			row = nil
		}
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Готово ➡️", "tags_done"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	text := "Выберите подходящие теги (можно несколько):"
	if len(selected) > 0 {
		text = fmt.Sprintf("Выбрано: %d. Выберите ещё или нажмите «Готово»:", len(selected))
	}
	sendWithKeyboard(bot, chatID, text, keyboard)
}

func (h *ScoreHandler) saveScore(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, comment string) {
	state, err := h.convSvc.GetState(ctx, user.TelegramID)
	if err != nil || state == nil {
		return
	}

	var flowData models.ScoreFlowData
	h.convSvc.GetData(state, &flowData)

	if flowData.Score == 0 {
		flowData.Score = 3
	}

	tags := flowData.Tags
	if tags == nil {
		tags = []string{}
	}

	score := &models.UserScore{
		UserID:        user.ID,
		ScoreType:     flowData.ScoreType,
		ReferenceType: flowData.ReferenceType,
		ReferenceID:   flowData.ReferenceID,
		Score:         flowData.Score,
		Comment:       comment,
		Tags:          tags,
	}
	h.scoreSvc.SaveScore(ctx, score)
	h.convSvc.ClearState(ctx, user.TelegramID)

	stars := strings.Repeat("⭐", flowData.Score)
	thank := fmt.Sprintf("Спасибо за отзыв! 🙏\n\nВаша оценка: %s", stars)
	if len(tags) > 0 {
		thank += "\nТеги: " + strings.Join(tags, ", ")
	}
	if comment != "" {
		thank += "\nКомментарий: " + comment
	}
	send(bot, chatID, thank)
}

func toggleTag(tags []string, tag string) []string {
	for i, t := range tags {
		if t == tag {
			return append(tags[:i], tags[i+1:]...)
		}
	}
	return append(tags, tag)
}
