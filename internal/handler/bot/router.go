package bot

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	userSvc    *service.UserService
	convSvc    *service.ConversationService
	questSvc   *service.QuestionnaireService
	subSvc     *service.SubscriptionService
	moduleSvc  *service.ModuleService
	scoreSvc   *service.ScoreService
	webAppURL  string

	regHandler     *RegistrationHandler
	quizHandler    *QuestionnaireHandler
	subHandler     *SubscriptionHandler
	modHandler     *ModuleHandler
	scoreHandler   *ScoreHandler
	profileHandler *ProfileHandler
}

func NewRouter(
	userSvc *service.UserService,
	convSvc *service.ConversationService,
	questSvc *service.QuestionnaireService,
	subSvc *service.SubscriptionService,
	moduleSvc *service.ModuleService,
	scoreSvc *service.ScoreService,
	webAppURL string,
) *Router {
	r := &Router{
		userSvc:   userSvc,
		convSvc:   convSvc,
		questSvc:  questSvc,
		subSvc:    subSvc,
		moduleSvc: moduleSvc,
		scoreSvc:  scoreSvc,
		webAppURL: webAppURL,
	}

	r.regHandler = NewRegistrationHandler(userSvc, convSvc, questSvc)
	r.quizHandler = NewQuestionnaireHandler(convSvc, questSvc)
	r.subHandler = NewSubscriptionHandler(convSvc, subSvc)
	r.modHandler = NewModuleHandler(convSvc, moduleSvc, subSvc, scoreSvc)
	r.scoreHandler = NewScoreHandler(convSvc, scoreSvc)
	r.profileHandler = NewProfileHandler(userSvc, subSvc)

	return r
}

func (r *Router) HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	ctx := context.Background()

	if update.CallbackQuery != nil {
		r.handleCallback(ctx, bot, update.CallbackQuery)
		return
	}

	if update.Message == nil {
		return
	}

	msg := update.Message
	user, err := r.userSvc.GetOrCreateUser(ctx, msg.From.ID, msg.From.UserName, msg.From.FirstName, msg.From.LastName)
	if err != nil {
		log.Printf("Error getting/creating user: %v", err)
		return
	}

	if msg.IsCommand() {
		r.handleCommand(ctx, bot, msg, user)
		return
	}

	r.handleMessage(ctx, bot, msg, user)
}

func (r *Router) handleCommand(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User) {
	switch strings.ToLower(msg.Command()) {
	case "start":
		r.handleStart(ctx, bot, msg, user)
	case "help":
		r.profileHandler.HandleHelp(bot, msg.Chat.ID)
	case "profile":
		if !r.requireRegistration(bot, msg.Chat.ID, user) {
			return
		}
		r.profileHandler.HandleProfile(ctx, bot, msg.Chat.ID, user)
	case "modules":
		if !r.requireRegistration(bot, msg.Chat.ID, user) {
			return
		}
		r.modHandler.HandleModuleList(ctx, bot, msg.Chat.ID, user)
	case "subscribe":
		if !r.requireRegistration(bot, msg.Chat.ID, user) {
			return
		}
		r.subHandler.HandleShowPlans(ctx, bot, msg.Chat.ID, user)
	case "app":
		if !r.requireRegistration(bot, msg.Chat.ID, user) {
			return
		}
		r.handleWebApp(bot, msg.Chat.ID)
	default:
		send(bot, msg.Chat.ID, "Неизвестная команда. Используйте /help для списка команд.")
	}
}

func (r *Router) handleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User) {
	state, err := r.convSvc.GetState(ctx, user.TelegramID)
	if err != nil {
		log.Printf("Error getting state: %v", err)
		return
	}

	if state == nil {
		send(bot, msg.Chat.ID, "Используйте /help для списка доступных команд.")
		return
	}

	switch {
	case service.IsFlowActive(state.State, "reg:"):
		r.regHandler.HandleMessage(ctx, bot, msg, user, state)
	case service.IsFlowActive(state.State, "quiz:"):
		r.quizHandler.HandleMessage(ctx, bot, msg, user, state)
	case service.IsFlowActive(state.State, "sub:"):
		r.subHandler.HandleMessage(ctx, bot, msg, user, state)
	case service.IsFlowActive(state.State, "mod:"):
		r.modHandler.HandleMessage(ctx, bot, msg, user, state)
	case service.IsFlowActive(state.State, "score:"):
		r.scoreHandler.HandleMessage(ctx, bot, msg, user, state)
	default:
		send(bot, msg.Chat.ID, "Используйте /help для списка доступных команд.")
	}
}

func (r *Router) handleCallback(ctx context.Context, bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery) {
	user, err := r.userSvc.GetOrCreateUser(ctx, cb.From.ID, cb.From.UserName, cb.From.FirstName, cb.From.LastName)
	if err != nil {
		log.Printf("Error getting user in callback: %v", err)
		return
	}

	callback := tgbotapi.NewCallback(cb.ID, "")
	bot.Send(callback)

	data := cb.Data
	chatID := cb.Message.Chat.ID

	switch {
	// Registration callbacks
	case strings.HasPrefix(data, "gender:"):
		value := strings.TrimPrefix(data, "gender:")
		r.regHandler.HandleGenderCallback(ctx, bot, chatID, user, value)
	case strings.HasPrefix(data, "fitness:"):
		value := strings.TrimPrefix(data, "fitness:")
		r.regHandler.HandleFitnessCallback(ctx, bot, chatID, user, value)
	case strings.HasPrefix(data, "goal:"):
		value := strings.TrimPrefix(data, "goal:")
		r.regHandler.HandleGoalCallback(ctx, bot, chatID, user, value)

	// Questionnaire callbacks
	case strings.HasPrefix(data, "quiz_ans:"):
		parts := strings.SplitN(strings.TrimPrefix(data, "quiz_ans:"), ":", 2)
		if len(parts) == 2 {
			qID, _ := strconv.Atoi(parts[0])
			r.quizHandler.HandleChoiceCallback(ctx, bot, chatID, user, qID, parts[1])
		}

	// Module callbacks
	case strings.HasPrefix(data, "mod:"):
		r.modHandler.HandleCallback(ctx, bot, chatID, user, data)
	case strings.HasPrefix(data, "cat:"):
		r.modHandler.HandleCallback(ctx, bot, chatID, user, data)
	case strings.HasPrefix(data, "les:"):
		r.modHandler.HandleCallback(ctx, bot, chatID, user, data)
	case strings.HasPrefix(data, "les_complete:"):
		r.modHandler.HandleCallback(ctx, bot, chatID, user, data)
	case data == "back_modules":
		r.modHandler.HandleModuleList(ctx, bot, chatID, user)
	case strings.HasPrefix(data, "back_cats:"):
		r.modHandler.HandleCallback(ctx, bot, chatID, user, data)

	// Subscription callbacks
	case strings.HasPrefix(data, "plan:"):
		r.subHandler.HandleCallback(ctx, bot, chatID, user, data)

	// Score callbacks
	case strings.HasPrefix(data, "score:"):
		r.scoreHandler.HandleCallback(ctx, bot, chatID, user, data)
	case data == "score_skip":
		r.scoreHandler.HandleCallback(ctx, bot, chatID, user, data)
	}
}

func (r *Router) handleStart(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User) {
	if !user.IsRegistered {
		r.regHandler.StartRegistration(ctx, bot, msg.Chat.ID, user)
		return
	}

	// Check if user has completed health test
	quiz, err := r.questSvc.GetBySlug(ctx, "health_test")
	if err == nil && quiz != nil {
		completed, _ := r.questSvc.HasCompleted(ctx, user.ID, quiz.ID)
		if !completed {
			r.quizHandler.StartQuestionnaire(ctx, bot, msg.Chat.ID, user, quiz)
			return
		}
	}

	text := "Добро пожаловать назад, " + user.FirstName + "!\n\n" +
		"Выберите раздел:\n" +
		"/modules — Модули обучения\n" +
		"/app — Открыть приложение\n" +
		"/subscribe — Подписка\n" +
		"/profile — Мой профиль\n" +
		"/help — Помощь"
	send(bot, msg.Chat.ID, text)
}

func (r *Router) requireRegistration(bot *tgbotapi.BotAPI, chatID int64, user *models.User) bool {
	if user.IsRegistered {
		return true
	}
	send(bot, chatID, "Для начала пройдите регистрацию: /start")
	return false
}

// telegram-bot-api v5.5.1 doesn't have native WebApp support, so we build the keyboard manually.
type webAppInfo struct {
	URL string `json:"url"`
}

type webAppKeyboardButton struct {
	Text   string      `json:"text"`
	WebApp *webAppInfo `json:"web_app,omitempty"`
}

type webAppReplyKeyboard struct {
	Keyboard        [][]webAppKeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool                     `json:"resize_keyboard"`
	OneTimeKeyboard bool                     `json:"one_time_keyboard"`
}

func (r *Router) handleWebApp(bot *tgbotapi.BotAPI, chatID int64) {
	if r.webAppURL == "" {
		send(bot, chatID, "Мини-приложение временно недоступно.")
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Нажмите кнопку ниже, чтобы открыть приложение:")
	msg.ReplyMarkup = webAppReplyKeyboard{
		Keyboard: [][]webAppKeyboardButton{
			{
				{
					Text:   "📱 Открыть приложение",
					WebApp: &webAppInfo{URL: r.webAppURL},
				},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	bot.Send(msg)
}

func send(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}

func sendWithKeyboard(bot *tgbotapi.BotAPI, chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}
