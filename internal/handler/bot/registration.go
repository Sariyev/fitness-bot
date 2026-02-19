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

	text := fmt.Sprintf("Привет, %s! 👋\n\n"+
		"Я — Андрей, твой персональный тренер. "+
		"Помогу тебе привести тело в форму и улучшить здоровье.\n\n"+
		"У меня есть для тебя:\n"+
		"🏥 *ЛФК* — упражнения при проблемах со спиной и суставами\n"+
		"💪 *Тренировки* — программы по группам мышц\n"+
		"🥗 *Питание* — рецепты и планы питания\n\n"+
		"Давай для начала заполним твой профиль.\n"+
		"Сколько тебе лет?",
		user.FirstName,
	)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
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
		// Text messages ignored in goal state — use inline buttons
		send(bot, msg.Chat.ID, "Выбери цели кнопками выше и нажми «✅ Готово»")
	}
}

func (h *RegistrationHandler) handleAge(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, data models.RegistrationData) {
	age, err := strconv.Atoi(msg.Text)
	if err != nil || age < 10 || age > 120 {
		send(bot, msg.Chat.ID, "Введи корректный возраст (от 10 до 120):")
		return
	}
	data.Age = age
	h.convSvc.SetState(ctx, user.TelegramID, models.StateRegHeight, data, registrationTTL)
	send(bot, msg.Chat.ID, "Отлично! Какой у тебя рост в сантиметрах?")
}

func (h *RegistrationHandler) handleHeight(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, data models.RegistrationData) {
	height, err := strconv.Atoi(msg.Text)
	if err != nil || height < 100 || height > 250 {
		send(bot, msg.Chat.ID, "Введи корректный рост (от 100 до 250 см):")
		return
	}
	data.HeightCm = height
	h.convSvc.SetState(ctx, user.TelegramID, models.StateRegWeight, data, registrationTTL)
	send(bot, msg.Chat.ID, "Принял! А вес в килограммах?")
}

func (h *RegistrationHandler) handleWeight(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, data models.RegistrationData) {
	weight, err := strconv.ParseFloat(msg.Text, 64)
	if err != nil || weight < 30 || weight > 300 {
		send(bot, msg.Chat.ID, "Введи корректный вес (от 30 до 300 кг):")
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
	sendWithKeyboard(bot, msg.Chat.ID, "Укажи свой пол:", keyboard)
}

func (h *RegistrationHandler) handleGender(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, data models.RegistrationData) {
	switch msg.Text {
	case "Мужской":
		data.Gender = "male"
	case "Женский":
		data.Gender = "female"
	default:
		send(bot, msg.Chat.ID, "Выбери пол кнопкой выше:")
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
	sendWithKeyboard(bot, chatID, "Какой у тебя уровень подготовки?", keyboard)
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
		send(bot, msg.Chat.ID, "Выбери уровень кнопкой выше:")
		return
	}
	h.advanceToGoal(ctx, bot, msg.Chat.ID, user, data)
}

// --- Goal multi-select ---

var goalOptions = []struct {
	Key   string
	Label string
}{
	{"weight_loss", "Похудеть"},
	{"muscle_gain", "Набрать массу"},
	{"strength", "Больше силы"},
	{"endurance", "Выносливость"},
	{"maintenance", "Поддержание формы"},
	{"hernia", "Грыжа"},
	{"protrusion", "Протрузии"},
	{"scoliosis", "Сколиоз"},
	{"kyphosis", "Кифоз"},
	{"lordosis", "Лордоз"},
}

func (h *RegistrationHandler) advanceToGoal(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data models.RegistrationData) {
	data.Goals = nil
	h.convSvc.SetState(ctx, user.TelegramID, models.StateRegGoal, data, registrationTTL)

	keyboard := buildGoalKeyboard(nil)
	msg := tgbotapi.NewMessage(chatID, "Выбери свои цели (можно несколько):\n\n"+
		"💪 — Фитнес-цели\n🏥 — Проблемы со здоровьем (ЛФК)")
	msg.ReplyMarkup = keyboard
	resp, _ := bot.Send(msg)

	// Save the message ID so we can edit the keyboard later
	data.GoalMsgID = resp.MessageID
	h.convSvc.SetState(ctx, user.TelegramID, models.StateRegGoal, data, registrationTTL)
}

func buildGoalKeyboard(selected []string) tgbotapi.InlineKeyboardMarkup {
	selectedSet := make(map[string]bool)
	for _, s := range selected {
		selectedSet[s] = true
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, opt := range goalOptions {
		label := opt.Label
		if selectedSet[opt.Key] {
			label = "✓ " + label
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, "goal:"+opt.Key),
		))
	}
	// Done button
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅ Готово", "goal:done"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (h *RegistrationHandler) HandleGoalCallback(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, messageID int, user *models.User, value string) {
	state, err := h.convSvc.GetState(ctx, user.TelegramID)
	if err != nil || state == nil || state.State != models.StateRegGoal {
		return
	}
	var data models.RegistrationData
	h.convSvc.GetData(state, &data)

	if value == "done" {
		if len(data.Goals) == 0 {
			cb := tgbotapi.NewCallback("", "Выбери хотя бы одну цель!")
			bot.Send(cb)
			return
		}
		h.completeRegistration(ctx, bot, chatID, user, data)
		return
	}

	// Toggle goal
	found := false
	var newGoals []string
	for _, g := range data.Goals {
		if g == value {
			found = true
			continue
		}
		newGoals = append(newGoals, g)
	}
	if !found {
		newGoals = append(newGoals, value)
	}
	data.Goals = newGoals

	h.convSvc.SetState(ctx, user.TelegramID, models.StateRegGoal, data, registrationTTL)

	// Edit the keyboard in-place
	msgID := data.GoalMsgID
	if msgID == 0 {
		msgID = messageID
	}
	keyboard := buildGoalKeyboard(data.Goals)
	edit := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, keyboard)
	bot.Send(edit)
}

func (h *RegistrationHandler) completeRegistration(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, data models.RegistrationData) {
	if err := h.userSvc.CreateProfile(ctx, user.ID, data); err != nil {
		log.Printf("Error creating profile: %v", err)
		send(bot, chatID, "Ошибка при сохранении профиля. Попробуй снова: /start")
		return
	}

	if err := h.userSvc.MarkRegistered(ctx, user); err != nil {
		log.Printf("Error marking user registered: %v", err)
	}

	h.convSvc.ClearState(ctx, user.TelegramID)

	goalsText := goalsLabel(data.Goals)

	send(bot, chatID, fmt.Sprintf(
		"Отлично, профиль готов! 🎉\n\n"+
			"Возраст: %d\n"+
			"Рост: %d см\n"+
			"Вес: %.1f кг\n"+
			"Пол: %s\n"+
			"Уровень: %s\n"+
			"Цели: %s\n\n"+
			"Сейчас пройдём короткий тест здоровья, чтобы я мог подобрать тебе программу...",
		data.Age, data.HeightCm, data.WeightKg,
		genderLabel(data.Gender), fitnessLabel(data.FitnessLevel), goalsText,
	))

	// Auto-trigger health questionnaire
	quiz, err := h.questSvc.GetBySlug(ctx, "health_test")
	if err != nil {
		log.Printf("Error loading health test: %v", err)
		send(bot, chatID, "Используй /modules для доступа к модулям.")
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
	for _, opt := range goalOptions {
		if opt.Key == g {
			return opt.Label
		}
	}
	return g
}

func goalsLabel(goals []string) string {
	if len(goals) == 0 {
		return "—"
	}
	var labels []string
	for _, g := range goals {
		labels = append(labels, goalLabel(g))
	}
	return strings.Join(labels, ", ")
}
