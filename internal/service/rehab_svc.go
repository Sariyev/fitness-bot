package service

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
	"time"
)

type RehabService struct {
	rehabRepo repository.RehabRepository
}

func NewRehabService(rehabRepo repository.RehabRepository) *RehabService {
	return &RehabService{rehabRepo: rehabRepo}
}

func (s *RehabService) ListCourses(ctx context.Context, category string) ([]models.RehabCourse, error) {
	return s.rehabRepo.ListCourses(ctx, category)
}

func (s *RehabService) GetCourse(ctx context.Context, id int) (*models.RehabCourse, error) {
	return s.rehabRepo.GetCourseByID(ctx, id)
}

func (s *RehabService) GetCourseSessions(ctx context.Context, courseID int) ([]models.RehabSession, error) {
	return s.rehabRepo.ListSessions(ctx, courseID)
}

func (s *RehabService) GetSession(ctx context.Context, id int) (*models.RehabSession, error) {
	return s.rehabRepo.GetSessionByID(ctx, id)
}

func (s *RehabService) CompleteSession(ctx context.Context, userID int64, courseID, sessionID, dayNumber, painLevel int, comment string) error {
	progress := &models.UserRehabProgress{
		UserID:      userID,
		CourseID:    courseID,
		SessionID:   sessionID,
		DayNumber:   dayNumber,
		CompletedAt: time.Now(),
		PainLevel:   painLevel,
		Comment:     comment,
	}
	return s.rehabRepo.CreateProgress(ctx, progress)
}

func (s *RehabService) GetUserProgress(ctx context.Context, userID int64, courseID int) ([]models.UserRehabProgress, error) {
	return s.rehabRepo.ListUserProgress(ctx, userID, courseID)
}
