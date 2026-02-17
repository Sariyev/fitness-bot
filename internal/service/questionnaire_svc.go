package service

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
)

type QuestionnaireService struct {
	repo repository.QuestionnaireRepository
}

func NewQuestionnaireService(repo repository.QuestionnaireRepository) *QuestionnaireService {
	return &QuestionnaireService{repo: repo}
}

func (s *QuestionnaireService) GetBySlug(ctx context.Context, slug string) (*models.Questionnaire, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *QuestionnaireService) GetQuestions(ctx context.Context, questionnaireID int) ([]models.Question, error) {
	return s.repo.GetQuestionsByQuestionnaireID(ctx, questionnaireID)
}

func (s *QuestionnaireService) GetQuestionByID(ctx context.Context, id int) (*models.Question, error) {
	return s.repo.GetQuestionByID(ctx, id)
}

func (s *QuestionnaireService) StartQuestionnaire(ctx context.Context, userID int64, questionnaireID int) (*models.QuestionnaireSubmission, error) {
	sub := &models.QuestionnaireSubmission{
		UserID:          userID,
		QuestionnaireID: questionnaireID,
	}
	if err := s.repo.CreateSubmission(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *QuestionnaireService) SaveAnswer(ctx context.Context, answer *models.QuestionnaireAnswer) error {
	return s.repo.SaveAnswer(ctx, answer)
}

func (s *QuestionnaireService) CompleteSubmission(ctx context.Context, submissionID int64) error {
	return s.repo.CompleteSubmission(ctx, submissionID)
}

func (s *QuestionnaireService) GetSubmissionAnswers(ctx context.Context, submissionID int64) ([]models.QuestionnaireAnswer, error) {
	return s.repo.GetSubmissionAnswers(ctx, submissionID)
}

func (s *QuestionnaireService) HasCompleted(ctx context.Context, userID int64, questionnaireID int) (bool, error) {
	subs, err := s.repo.GetUserSubmissions(ctx, userID, questionnaireID)
	if err != nil {
		return false, err
	}
	for _, sub := range subs {
		if sub.CompletedAt != nil {
			return true, nil
		}
	}
	return false, nil
}
