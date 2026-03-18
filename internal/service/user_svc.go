package service

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
	"strings"

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
	trainingAccess := &data.TrainingAccess
	if data.TrainingAccess == "" {
		trainingAccess = nil
	}
	trainingExperience := &data.TrainingExperience
	if data.TrainingExperience == "" {
		trainingExperience = nil
	}
	preferredTime := &data.PreferredTime
	if data.PreferredTime == "" {
		preferredTime = nil
	}

	profile := &models.UserProfile{
		UserID:             userID,
		Gender:             data.Gender,
		Age:                data.Age,
		WeightKg:           data.WeightKg,
		HeightCm:           data.HeightCm,
		FitnessLevel:       data.FitnessLevel,
		Goal:               strings.Join(data.Goals, ","),
		TrainingAccess:     trainingAccess,
		TrainingExperience: trainingExperience,
		HasPain:            data.HasPain,
		PainLocations:      data.PainLocations,
		PainLevel:          data.PainLevel,
		Diagnoses:          data.Diagnoses,
		Contraindications:  data.Contraindications,
		DaysPerWeek:        data.DaysPerWeek,
		SessionDuration:    data.SessionDuration,
		PreferredTime:      preferredTime,
		Equipment:          data.Equipment,
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

func (s *UserService) UpdateProfile(ctx context.Context, profile *models.UserProfile) error {
	return s.repo.UpdateProfile(ctx, profile)
}

func (s *UserService) UpdateProfileFromData(ctx context.Context, userID int64, data models.RegistrationData) error {
	profile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}

	trainingAccess := &data.TrainingAccess
	if data.TrainingAccess == "" {
		trainingAccess = nil
	}
	trainingExperience := &data.TrainingExperience
	if data.TrainingExperience == "" {
		trainingExperience = nil
	}
	preferredTime := &data.PreferredTime
	if data.PreferredTime == "" {
		preferredTime = nil
	}

	profile.Gender = data.Gender
	profile.Age = data.Age
	profile.WeightKg = data.WeightKg
	profile.HeightCm = data.HeightCm
	profile.FitnessLevel = data.FitnessLevel
	profile.Goal = strings.Join(data.Goals, ",")
	profile.TrainingAccess = trainingAccess
	profile.TrainingExperience = trainingExperience
	profile.HasPain = data.HasPain
	profile.PainLocations = data.PainLocations
	profile.PainLevel = data.PainLevel
	profile.Diagnoses = data.Diagnoses
	profile.Contraindications = data.Contraindications
	profile.DaysPerWeek = data.DaysPerWeek
	profile.SessionDuration = data.SessionDuration
	profile.PreferredTime = preferredTime
	profile.Equipment = data.Equipment

	return s.repo.UpdateProfile(ctx, profile)
}
