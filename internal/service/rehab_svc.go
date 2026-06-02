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

// ListAllCourses returns every course including inactive — admin only.
func (s *RehabService) ListAllCourses(ctx context.Context) ([]models.RehabCourse, error) {
	return s.rehabRepo.ListAllCourses(ctx)
}

func (s *RehabService) CreateCourse(ctx context.Context, c *models.RehabCourse) error {
	return s.rehabRepo.CreateCourse(ctx, c)
}

func (s *RehabService) UpdateCourse(ctx context.Context, c *models.RehabCourse) error {
	return s.rehabRepo.UpdateCourse(ctx, c)
}

func (s *RehabService) DeleteCourse(ctx context.Context, id int) error {
	return s.rehabRepo.DeleteCourse(ctx, id)
}

func (s *RehabService) CreateSession(ctx context.Context, sess *models.RehabSession) error {
	return s.rehabRepo.CreateSession(ctx, sess)
}

func (s *RehabService) UpdateSession(ctx context.Context, sess *models.RehabSession) error {
	return s.rehabRepo.UpdateSession(ctx, sess)
}

func (s *RehabService) DeleteSession(ctx context.Context, id int) error {
	return s.rehabRepo.DeleteSession(ctx, id)
}
