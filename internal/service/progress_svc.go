package service

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
)

type ProgressService struct {
	progressRepo    repository.ProgressRepository
	completionRepo  repository.DailyCompletionRepository
	achievementRepo repository.AchievementRepository
}

func NewProgressService(
	progressRepo repository.ProgressRepository,
	completionRepo repository.DailyCompletionRepository,
	achievementRepo repository.AchievementRepository,
) *ProgressService {
	return &ProgressService{
		progressRepo:    progressRepo,
		completionRepo:  completionRepo,
		achievementRepo: achievementRepo,
	}
}

func (s *ProgressService) AddEntry(ctx context.Context, entry *models.ProgressEntry) error {
	return s.progressRepo.Create(ctx, entry)
}

func (s *ProgressService) ListEntries(ctx context.Context, userID int64) ([]models.ProgressEntry, error) {
	return s.progressRepo.ListByUser(ctx, userID)
}

func (s *ProgressService) GetWeightHistory(ctx context.Context, userID int64) ([]models.WeightPoint, error) {
	return s.progressRepo.GetWeightHistory(ctx, userID)
}

func (s *ProgressService) GetStreak(ctx context.Context, userID int64) (current int, longest int, err error) {
	return s.completionRepo.GetStreak(ctx, userID)
}

func (s *ProgressService) GetCalendar(ctx context.Context, userID int64, year, month int) ([]string, error) {
	return s.completionRepo.GetCalendar(ctx, userID, year, month)
}

func (s *ProgressService) GetAchievements(ctx context.Context, userID int64) ([]models.UserAchievement, error) {
	return s.achievementRepo.ListByUser(ctx, userID)
}
