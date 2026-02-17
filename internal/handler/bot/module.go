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

const moduleTTL = time.Hour

type ModuleHandler struct {
	convSvc  *service.ConversationService
	modSvc   *service.ModuleService
	subSvc   *service.SubscriptionService
	scoreSvc *service.ScoreService
}

func NewModuleHandler(convSvc *service.ConversationService, modSvc *service.ModuleService, subSvc *service.SubscriptionService, scoreSvc *service.ScoreService) *ModuleHandler {
	return &ModuleHandler{convSvc: convSvc, modSvc: modSvc, subSvc: subSvc, scoreSvc: scoreSvc}
}

func (h *ModuleHandler) HandleModuleList(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User) {
	modules, err := h.modSvc.ListModules(ctx)
	if err != nil || len(modules) == 0 {
		send(bot, chatID, "Модули временно недоступны.")
		return
	}

	text := "📚 Доступные модули:\n\n"
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, m := range modules {
		text += fmt.Sprintf("%s %s\n%s\n\n", m.Icon, m.Name, m.Description)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %s", m.Icon, m.Name),
				fmt.Sprintf("mod:%d", m.ID),
			),
		))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	sendWithKeyboard(bot, chatID, text, keyboard)
}

func (h *ModuleHandler) HandleCallback(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data string) {
	switch {
	case strings.HasPrefix(data, "mod:"):
		h.handleModuleSelect(ctx, bot, chatID, user, data)
	case strings.HasPrefix(data, "cat:"):
		h.handleCategorySelect(ctx, bot, chatID, user, data)
	case strings.HasPrefix(data, "les:"):
		h.handleLessonSelect(ctx, bot, chatID, user, data)
	case strings.HasPrefix(data, "les_complete:"):
		h.handleLessonComplete(ctx, bot, chatID, user, data)
	case strings.HasPrefix(data, "back_cats:"):
		idStr := strings.TrimPrefix(data, "back_cats:")
		moduleID, _ := strconv.Atoi(idStr)
		h.showCategories(ctx, bot, chatID, user, moduleID)
	}
}

func (h *ModuleHandler) handleModuleSelect(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data string) {
	idStr := strings.TrimPrefix(data, "mod:")
	moduleID, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}

	module, err := h.modSvc.GetModule(ctx, moduleID)
	if err != nil {
		send(bot, chatID, "Модуль не найден.")
		return
	}

	// Check subscription
	if module.RequiresSubscription {
		active, _ := h.subSvc.HasActiveSubscription(ctx, user.ID)
		if !active {
			send(bot, chatID, "⚠️ Для доступа к этому модулю нужна подписка.\n\n/subscribe — Оформить подписку")
			return
		}
	}

	h.showCategories(ctx, bot, chatID, user, moduleID)
}

func (h *ModuleHandler) showCategories(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, moduleID int) {
	module, err := h.modSvc.GetModule(ctx, moduleID)
	if err != nil {
		return
	}

	categories, err := h.modSvc.ListCategories(ctx, moduleID)
	if err != nil || len(categories) == 0 {
		send(bot, chatID, "В этом модуле пока нет категорий.")
		return
	}

	text := fmt.Sprintf("%s %s\n\nВыберите категорию:\n", module.Icon, module.Name)
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, c := range categories {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s %s", c.Icon, c.Name),
				fmt.Sprintf("cat:%d", c.ID),
			),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("« Назад к модулям", "back_modules"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	sendWithKeyboard(bot, chatID, text, keyboard)
}

func (h *ModuleHandler) handleCategorySelect(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data string) {
	idStr := strings.TrimPrefix(data, "cat:")
	catID, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}

	category, err := h.modSvc.GetCategory(ctx, catID)
	if err != nil {
		send(bot, chatID, "Категория не найдена.")
		return
	}

	h.modSvc.SelectCategory(ctx, user.ID, catID)

	lessons, err := h.modSvc.ListLessons(ctx, catID)
	if err != nil || len(lessons) == 0 {
		send(bot, chatID, "В этой категории пока нет занятий.")
		return
	}

	completed, total, _ := h.modSvc.GetCategoryProgress(ctx, user.ID, catID)

	text := fmt.Sprintf("%s %s\n\nПрогресс: %d/%d занятий\n\nВыберите занятие:\n",
		category.Icon, category.Name, completed, total)

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, l := range lessons {
		label := l.Title
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, fmt.Sprintf("les:%d", l.ID)),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("« Назад", fmt.Sprintf("back_cats:%d", category.ModuleID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	sendWithKeyboard(bot, chatID, text, keyboard)
}

func (h *ModuleHandler) handleLessonSelect(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data string) {
	idStr := strings.TrimPrefix(data, "les:")
	lessonID, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}

	lesson, err := h.modSvc.GetLesson(ctx, lessonID)
	if err != nil {
		send(bot, chatID, "Занятие не найдено.")
		return
	}

	h.modSvc.StartLesson(ctx, user.ID, lessonID)

	contents, err := h.modSvc.GetLessonContent(ctx, lessonID)
	if err != nil {
		log.Printf("Error loading lesson content: %v", err)
	}

	// Send lesson header
	text := fmt.Sprintf("📖 %s\n\n%s", lesson.Title, lesson.Description)
	send(bot, chatID, text)

	// Send all content pieces
	for _, c := range contents {
		h.sendContent(ctx, bot, chatID, c)
	}

	// Completion button
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Занятие завершено", fmt.Sprintf("les_complete:%d", lessonID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Назад", fmt.Sprintf("cat:%d", lesson.CategoryID)),
		),
	)
	sendWithKeyboard(bot, chatID, "Нажмите когда завершите занятие:", keyboard)

	flowData := models.ModuleBrowseData{
		LessonID:   lessonID,
		CategoryID: lesson.CategoryID,
	}
	h.convSvc.SetState(ctx, user.TelegramID, models.StateModViewLesson, flowData, moduleTTL)
}

func (h *ModuleHandler) handleLessonComplete(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data string) {
	idStr := strings.TrimPrefix(data, "les_complete:")
	lessonID, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}

	h.modSvc.CompleteLesson(ctx, user.ID, lessonID)
	h.convSvc.ClearState(ctx, user.TelegramID)

	send(bot, chatID, "🎉 Отличная работа! Занятие завершено!")

	// Trigger score collection
	scoreHandler := NewScoreHandler(h.convSvc, h.scoreSvc)
	scoreHandler.RequestScore(ctx, bot, chatID, user, "lesson_complete", "lesson", lessonID)
}

func (h *ModuleHandler) sendContent(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, c models.LessonContent) {
	switch c.ContentType {
	case "text":
		if c.Body != "" {
			msg := tgbotapi.NewMessage(chatID, c.Body)
			msg.ParseMode = "Markdown"
			bot.Send(msg)
		}

	case "video":
		if c.TelegramFileID != "" {
			video := tgbotapi.NewVideo(chatID, tgbotapi.FileID(c.TelegramFileID))
			if c.Title != "" {
				video.Caption = c.Title
			}
			bot.Send(video)
			return
		}
		if c.VideoURL != "" {
			text := c.VideoURL
			if c.Title != "" {
				text = fmt.Sprintf("🎬 *%s*\n\n%s", c.Title, c.VideoURL)
			}
			msg := tgbotapi.NewMessage(chatID, text)
			msg.ParseMode = "Markdown"
			bot.Send(msg)
			return
		}
		if c.FileURL != "" {
			video := tgbotapi.NewVideo(chatID, tgbotapi.FileURL(c.FileURL))
			if c.Title != "" {
				video.Caption = c.Title
			}
			resp, err := bot.Send(video)
			if err == nil && resp.Video != nil {
				h.modSvc.UpdateTelegramFileID(ctx, c.ID, resp.Video.FileID)
			}
		}

	case "image":
		if c.TelegramFileID != "" {
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(c.TelegramFileID))
			if c.Title != "" {
				photo.Caption = c.Title
			}
			bot.Send(photo)
		} else if c.FileURL != "" {
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(c.FileURL))
			if c.Title != "" {
				photo.Caption = c.Title
			}
			resp, err := bot.Send(photo)
			if err == nil && resp.Photo != nil && len(resp.Photo) > 0 {
				h.modSvc.UpdateTelegramFileID(ctx, c.ID, resp.Photo[len(resp.Photo)-1].FileID)
			}
		}
	}
}

func (h *ModuleHandler) HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, state *models.ConversationState) {
	send(bot, msg.Chat.ID, "Пожалуйста, используйте кнопки для навигации.")
}
