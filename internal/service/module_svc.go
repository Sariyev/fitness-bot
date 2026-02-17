package service

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
)

type ModuleService struct {
	repo repository.ModuleRepository
}

func NewModuleService(repo repository.ModuleRepository) *ModuleService {
	return &ModuleService{repo: repo}
}

func (s *ModuleService) ListModules(ctx context.Context) ([]models.Module, error) {
	return s.repo.ListActiveModules(ctx)
}

func (s *ModuleService) GetModule(ctx context.Context, id int) (*models.Module, error) {
	return s.repo.GetModuleByID(ctx, id)
}

func (s *ModuleService) ListCategories(ctx context.Context, moduleID int) ([]models.ModuleCategory, error) {
	return s.repo.ListCategoriesByModule(ctx, moduleID)
}

func (s *ModuleService) GetCategory(ctx context.Context, id int) (*models.ModuleCategory, error) {
	return s.repo.GetCategoryByID(ctx, id)
}

func (s *ModuleService) ListLessons(ctx context.Context, categoryID int) ([]models.Lesson, error) {
	return s.repo.ListLessonsByCategory(ctx, categoryID)
}

func (s *ModuleService) GetLesson(ctx context.Context, id int) (*models.Lesson, error) {
	return s.repo.GetLessonByID(ctx, id)
}

func (s *ModuleService) GetLessonContent(ctx context.Context, lessonID int) ([]models.LessonContent, error) {
	return s.repo.ListContentByLesson(ctx, lessonID)
}

func (s *ModuleService) StartLesson(ctx context.Context, userID int64, lessonID int) error {
	return s.repo.UpsertProgress(ctx, userID, lessonID)
}

func (s *ModuleService) CompleteLesson(ctx context.Context, userID int64, lessonID int) error {
	return s.repo.CompleteLesson(ctx, userID, lessonID)
}

func (s *ModuleService) GetCategoryProgress(ctx context.Context, userID int64, categoryID int) (int, int, error) {
	return s.repo.GetUserCategoryProgress(ctx, userID, categoryID)
}

func (s *ModuleService) SelectCategory(ctx context.Context, userID int64, categoryID int) error {
	return s.repo.SelectCategory(ctx, userID, categoryID)
}

func (s *ModuleService) GetLessonProgress(ctx context.Context, userID int64, lessonID int) (*models.UserLessonProgress, error) {
	return s.repo.GetUserProgress(ctx, userID, lessonID)
}

func (s *ModuleService) UpdateTelegramFileID(ctx context.Context, contentID int, fileID string) error {
	return s.repo.UpdateTelegramFileID(ctx, contentID, fileID)
}
