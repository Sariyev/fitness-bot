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

const questionnaireTTL = 30 * time.Minute

type QuestionnaireHandler struct {
	convSvc  *service.ConversationService
	questSvc *service.QuestionnaireService
}

func NewQuestionnaireHandler(convSvc *service.ConversationService, questSvc *service.QuestionnaireService) *QuestionnaireHandler {
	return &QuestionnaireHandler{convSvc: convSvc, questSvc: questSvc}
}

func (h *QuestionnaireHandler) StartQuestionnaire(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, quiz *models.Questionnaire) {
	questions, err := h.questSvc.GetQuestions(ctx, quiz.ID)
	if err != nil || len(questions) == 0 {
		log.Printf("Error loading questions: %v", err)
		send(bot, chatID, "Тест временно недоступен. Используйте /modules.")
		return
	}

	sub, err := h.questSvc.StartQuestionnaire(ctx, user.ID, quiz.ID)
	if err != nil {
		log.Printf("Error starting questionnaire: %v", err)
		return
	}

	flowData := models.QuizFlowData{
		SubmissionID:       sub.ID,
		QuestionnaireID:    quiz.ID,
		CurrentQuestionIdx: 0,
		TotalQuestions:     len(questions),
	}

	h.convSvc.SetState(ctx, user.TelegramID, models.StateQuizInProgress, flowData, questionnaireTTL)
	send(bot, chatID, fmt.Sprintf("📋 %s\n\n%s\n\nВопрос 1 из %d:", quiz.Title, quiz.Description, len(questions)))
	h.sendQuestion(bot, chatID, questions[0])
}

func (h *QuestionnaireHandler) HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, user *models.User, state *models.ConversationState) {
	var flowData models.QuizFlowData
	h.convSvc.GetData(state, &flowData)

	questions, err := h.questSvc.GetQuestions(ctx, flowData.QuestionnaireID)
	if err != nil || flowData.CurrentQuestionIdx >= len(questions) {
		send(bot, msg.Chat.ID, "Произошла ошибка. Попробуйте /start")
		h.convSvc.ClearState(ctx, user.TelegramID)
		return
	}

	currentQ := questions[flowData.CurrentQuestionIdx]

	// For text type questions, save the text answer
	if currentQ.QuestionType == "text" {
		answer := &models.QuestionnaireAnswer{
			SubmissionID: flowData.SubmissionID,
			QuestionID:   currentQ.ID,
			AnswerText:   &msg.Text,
		}
		h.questSvc.SaveAnswer(ctx, answer)
		h.advanceQuestion(ctx, bot, msg.Chat.ID, user, flowData, questions)
		return
	}

	// For scale type, try to parse number
	if currentQ.QuestionType == "scale" {
		val, err := strconv.Atoi(msg.Text)
		if err != nil || val < 1 || val > 10 {
			send(bot, msg.Chat.ID, "Пожалуйста, введите число от 1 до 10:")
			return
		}
		valStr := msg.Text
		answer := &models.QuestionnaireAnswer{
			SubmissionID: flowData.SubmissionID,
			QuestionID:   currentQ.ID,
			AnswerValue:  &valStr,
		}
		h.questSvc.SaveAnswer(ctx, answer)
		h.advanceQuestion(ctx, bot, msg.Chat.ID, user, flowData, questions)
		return
	}

	send(bot, msg.Chat.ID, "Пожалуйста, используйте кнопки для ответа.")
}

func (h *QuestionnaireHandler) HandleChoiceCallback(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, questionID int, value string) {
	state, err := h.convSvc.GetState(ctx, user.TelegramID)
	if err != nil || state == nil || state.State != models.StateQuizInProgress {
		return
	}

	var flowData models.QuizFlowData
	h.convSvc.GetData(state, &flowData)

	answer := &models.QuestionnaireAnswer{
		SubmissionID: flowData.SubmissionID,
		QuestionID:   questionID,
		AnswerValue:  &value,
	}
	h.questSvc.SaveAnswer(ctx, answer)

	questions, err := h.questSvc.GetQuestions(ctx, flowData.QuestionnaireID)
	if err != nil {
		return
	}

	h.advanceQuestion(ctx, bot, chatID, user, flowData, questions)
}

func (h *QuestionnaireHandler) advanceQuestion(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, flowData models.QuizFlowData, questions []models.Question) {
	flowData.CurrentQuestionIdx++

	if flowData.CurrentQuestionIdx >= flowData.TotalQuestions {
		h.completeQuestionnaire(ctx, bot, chatID, user, flowData)
		return
	}

	h.convSvc.SetState(ctx, user.TelegramID, models.StateQuizInProgress, flowData, questionnaireTTL)
	nextQ := questions[flowData.CurrentQuestionIdx]
	send(bot, chatID, fmt.Sprintf("Вопрос %d из %d:", flowData.CurrentQuestionIdx+1, flowData.TotalQuestions))
	h.sendQuestion(bot, chatID, nextQ)
}

func (h *QuestionnaireHandler) completeQuestionnaire(ctx context.Context, bot *tgbotapi.BotAPI, chatID int64, user *models.User, flowData models.QuizFlowData) {
	h.questSvc.CompleteSubmission(ctx, flowData.SubmissionID)
	h.convSvc.ClearState(ctx, user.TelegramID)

	text := "✅ Тест завершён!\n\n" +
		"На основе твоих ответов я подготовлю персональные рекомендации.\n" +
		"Для полного доступа к программам оплати доступ.\n\n" +
		"/buy — Оплатить доступ\n" +
		"/modules — Перейти к модулям"
	send(bot, chatID, text)
}

func (h *QuestionnaireHandler) sendQuestion(bot *tgbotapi.BotAPI, chatID int64, q models.Question) {
	switch q.QuestionType {
	case "single_choice":
		var rows [][]tgbotapi.InlineKeyboardButton
		for _, opt := range q.Options {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(opt.Text, fmt.Sprintf("quiz_ans:%d:%s", q.ID, opt.Value)),
			))
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
		sendWithKeyboard(bot, chatID, q.Text, keyboard)

	case "multiple_choice":
		var rows [][]tgbotapi.InlineKeyboardButton
		for _, opt := range q.Options {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(opt.Text, fmt.Sprintf("quiz_ans:%d:%s", q.ID, opt.Value)),
			))
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
		sendWithKeyboard(bot, chatID, q.Text+"\n\n(Выберите один вариант)", keyboard)

	case "scale":
		var buttons []tgbotapi.InlineKeyboardButton
		for i := 1; i <= 10; i++ {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
				strconv.Itoa(i),
				fmt.Sprintf("quiz_ans:%d:%d", q.ID, i),
			))
		}
		row1 := buttons[:5]
		row2 := buttons[5:]
		keyboard := tgbotapi.NewInlineKeyboardMarkup(row1, row2)
		sendWithKeyboard(bot, chatID, q.Text, keyboard)

	case "text":
		send(bot, chatID, q.Text+"\n\n(Введите ваш ответ текстом)")
	}
}
