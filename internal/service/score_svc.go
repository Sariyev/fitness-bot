package service

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"

	"github.com/jackc/pgx/v4"
)

type ScoreService struct {
	repo repository.ScoreRepository
}

func NewScoreService(repo repository.ScoreRepository) *ScoreService {
	return &ScoreService{repo: repo}
}

func (s *ScoreService) SaveScore(ctx context.Context, score *models.UserScore) error {
	return s.repo.Create(ctx, score)
}

func (s *ScoreService) HasScored(ctx context.Context, userID int64, refType string, refID int) (bool, error) {
	_, err := s.repo.GetByReference(ctx, userID, refType, refID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *ScoreService) ListByUser(ctx context.Context, userID int64) ([]models.UserScore, error) {
	return s.repo.ListByUser(ctx, userID)
}
