package bot

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
	"fmt"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const registrationTTL = 24 * time.Hour

type RegistrationHandler struct {
	userSvc  *service.UserService
	convSvc  *service.ConversationService
	questSvc *service.QuestionnaireService
}

func NewRegistrationHandler(userSvc *service.UserService, convSvc *service.ConversationService, questSvc *service.QuestionnaireService) *RegistrationHandler {
	return &RegistrationHandler{userSvc: userSvc, convSvc: convSvc, questSvc: questSvc}
}

func (h *RegistrationHandler) StartRegistration(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User) {
	h.convSvc.SetState(ctx, user.TelegramID, models.StateRegAge, models.RegistrationData{}, registrationTTL)
	send(bot, chatID, "Давайте создадим ваш профиль!\n\nСколько вам лет?")
}

func (h *RegistrationHandler) HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, state *models.ConversationState) {
	var data models.RegistrationData
	h.convSvc.GetData(state, &data)

	switch state.State {
	case models.StateRegAge:
		h.handleAge(ctx, bot, msg, user, data)
	case models.StateRegHeight:
		h.handleHeight(ctx, bot, msg, user, data)
	case models.StateRegWeight:
		h.handleWeight(ctx, bot, msg, user, data)
	case models.StateRegGender:
		h.handleGender(ctx, bot, msg, user, data)
	case models.StateRegFitnessLevel:
		h.handleFitnessLevel(ctx, bot, msg, user, data)
	case models.StateRegGoal:
		h.handleGoal(ctx, bot, msg, user, data)
	}
}

func (h *RegistrationHandler) handleAge(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, data models.RegistrationData) {
	age, err := strconv.Atoi(msg.Text)
	if err != nil || age < 10 || age > 120 {
		send(bot, msg.Chat.ID, "Пожалуйста, введите корректный возраст (от 10 до 120):")
		return
	}
	data.Age = age
	h.convSvc.SetState(ctx, user.TelegramID, models.StateRegHeight, data, registrationTTL)
	send(bot, msg.Chat.ID, "Укажите ваш рост в сантиметрах:")
}

func (h *RegistrationHandler) handleHeight(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, data models.RegistrationData) {
	height, err := strconv.Atoi(msg.Text)
	if err != nil || height < 100 || height > 250 {
		send(bot, msg.Chat.ID, "Пожалуйста, введите корректный рост (от 100 до 250 см):")
		return
	}
	data.HeightCm = height
	h.convSvc.SetState(ctx, user.TelegramID, models.StateRegWeight, data, registrationTTL)
	send(bot, msg.Chat.ID, "Укажите ваш вес в килограммах:")
}

func (h *RegistrationHandler) handleWeight(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, data models.RegistrationData) {
	weight, err := strconv.ParseFloat(msg.Text, 64)
	if err != nil || weight < 30 || weight > 300 {
		send(bot, msg.Chat.ID, "Пожалуйста, введите корректный вес (от 30 до 300 кг):")
		return
	}
	data.WeightKg = weight
	h.convSvc.SetState(ctx, user.TelegramID, models.StateRegGender, data, registrationTTL)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Мужской", "gender:male"),
			tgbotapi.NewInlineKeyboardButtonData("Женский", "gender:female"),
		),
	)
	sendWithKeyboard(bot, msg.Chat.ID, "Укажите ваш пол:", keyboard)
}

func (h *RegistrationHandler) handleGender(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, data models.RegistrationData) {
	text := msg.Text
	switch text {
	case "Мужской":
		data.Gender = "male"
	case "Женский":
		data.Gender = "female"
	default:
		send(bot, msg.Chat.ID, "Пожалуйста, выберите пол:")
		return
	}
	h.advanceToFitness(ctx, bot, msg.Chat.ID, user, data)
}

func (h *RegistrationHandler) HandleGenderCallback(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, value string) {
	state, err := h.convSvc.GetState(ctx, user.TelegramID)
	if err != nil || state == nil || state.State != models.StateRegGender {
		return
	}
	var data models.RegistrationData
	h.convSvc.GetData(state, &data)
	data.Gender = value
	h.advanceToFitness(ctx, bot, chatID, user, data)
}

func (h *RegistrationHandler) advanceToFitness(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data models.RegistrationData) {
	h.convSvc.SetState(ctx, user.TelegramID, models.StateRegFitnessLevel, data, registrationTTL)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Новичок", "fitness:beginner"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Средний", "fitness:intermediate"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Продвинутый", "fitness:advanced"),
		),
	)
	sendWithKeyboard(bot, chatID, "Ваш уровень физической подготовки:", keyboard)
}

func (h *RegistrationHandler) HandleFitnessCallback(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, value string) {
	state, err := h.convSvc.GetState(ctx, user.TelegramID)
	if err != nil || state == nil || state.State != models.StateRegFitnessLevel {
		return
	}
	var data models.RegistrationData
	h.convSvc.GetData(state, &data)
	data.FitnessLevel = value
	h.advanceToGoal(ctx, bot, chatID, user, data)
}

func (h *RegistrationHandler) handleFitnessLevel(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, data models.RegistrationData) {
	switch msg.Text {
	case "Новичок":
		data.FitnessLevel = "beginner"
	case "Средний":
		data.FitnessLevel = "intermediate"
	case "Продвинутый":
		data.FitnessLevel = "advanced"
	default:
		send(bot, msg.Chat.ID, "Пожалуйста, выберите уровень подготовки:")
		return
	}
	h.advanceToGoal(ctx, bot, msg.Chat.ID, user, data)
}

func (h *RegistrationHandler) advanceToGoal(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data models.RegistrationData) {
	h.convSvc.SetState(ctx, user.TelegramID, models.StateRegGoal, data, registrationTTL)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Похудеть", "goal:weight_loss"),
			tgbotapi.NewInlineKeyboardButtonData("Набрать массу", "goal:muscle_gain"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Больше силы", "goal:strength"),
			tgbotapi.NewInlineKeyboardButtonData("Выносливость", "goal:endurance"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Поддержание формы", "goal:maintenance"),
		),
	)
	sendWithKeyboard(bot, chatID, "Какая ваша основная цель?", keyboard)
}

func (h *RegistrationHandler) HandleGoalCallback(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, value string) {
	state, err := h.convSvc.GetState(ctx, user.TelegramID)
	if err != nil || state == nil || state.State != models.StateRegGoal {
		return
	}
	var data models.RegistrationData
	h.convSvc.GetData(state, &data)
	data.Goal = value
	h.completeRegistration(ctx, bot, chatID, user, data)
}

func (h *RegistrationHandler) handleGoal(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, data models.RegistrationData) {
	goalMap := map[string]string{
		"Похудеть":            "weight_loss",
		"Набрать массу":      "muscle_gain",
		"Больше силы":        "strength",
		"Выносливость":       "endurance",
		"Поддержание формы":  "maintenance",
	}
	goal, ok := goalMap[msg.Text]
	if !ok {
		send(bot, msg.Chat.ID, "Пожалуйста, выберите цель:")
		return
	}
	data.Goal = goal
	h.completeRegistration(ctx, bot, msg.Chat.ID, user, data)
}

func (h *RegistrationHandler) completeRegistration(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data models.RegistrationData) {
	if err := h.userSvc.CreateProfile(ctx, user.ID, data); err != nil {
		log.Printf("Error creating profile: %v", err)
		send(bot, chatID, "Произошла ошибка при сохранении профиля. Попробуйте снова: /start")
		return
	}

	if err := h.userSvc.MarkRegistered(ctx, user); err != nil {
		log.Printf("Error marking user registered: %v", err)
	}

	h.convSvc.ClearState(ctx, user.TelegramID)

	send(bot, chatID, fmt.Sprintf(
		"Профиль создан!\n\n"+
			"Возраст: %d\n"+
			"Рост: %d см\n"+
			"Вес: %.1f кг\n"+
			"Пол: %s\n"+
			"Уровень: %s\n"+
			"Цель: %s\n\n"+
			"Теперь пройдите тест здоровья...",
		data.Age, data.HeightCm, data.WeightKg,
		genderLabel(data.Gender), fitnessLabel(data.FitnessLevel), goalLabel(data.Goal),
	))

	// Auto-trigger health questionnaire
	quiz, err := h.questSvc.GetBySlug(ctx, "health_test")
	if err != nil {
		log.Printf("Error loading health test: %v", err)
		send(bot, chatID, "Используйте /modules для доступа к модулям.")
		return
	}
	quizHandler := NewQuestionnaireHandler(h.convSvc, h.questSvc)
	quizHandler.StartQuestionnaire(ctx, bot, chatID, user, quiz)
}

func genderLabel(g string) string {
	if g == "male" {
		return "Мужской"
	}
	return "Женский"
}

func fitnessLabel(f string) string {
	labels := map[string]string{
		"beginner":     "Новичок",
		"intermediate": "Средний",
		"advanced":     "Продвинутый",
	}
	if l, ok := labels[f]; ok {
		return l
	}
	return f
}

func goalLabel(g string) string {
	labels := map[string]string{
		"weight_loss":  "Похудеть",
		"muscle_gain":  "Набрать массу",
		"strength":     "Больше силы",
		"endurance":    "Выносливость",
		"maintenance":  "Поддержание формы",
	}
	if l, ok := labels[g]; ok {
		return l
	}
	return g
}
