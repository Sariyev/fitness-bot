package service

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
	"fmt"
	"strings"
)

type Recommendation struct {
	ProgramID       *int   `json:"program_id,omitempty"`
	ProgramName     string `json:"program_name,omitempty"`
	RehabCourseID   *int   `json:"rehab_course_id,omitempty"`
	RehabCourseName string `json:"rehab_course_name,omitempty"`
	Calories        int    `json:"calories"`
	Protein         int    `json:"protein"`
	Fat             int    `json:"fat"`
	Carbs           int    `json:"carbs"`
	Message         string `json:"message"`
}

type RecommendationService struct {
	programRepo   repository.ProgramRepository
	rehabRepo     repository.RehabRepository
	nutritionRepo repository.NutritionRepository
}

func NewRecommendationService(
	programRepo repository.ProgramRepository,
	rehabRepo repository.RehabRepository,
	nutritionRepo repository.NutritionRepository,
) *RecommendationService {
	return &RecommendationService{
		programRepo:   programRepo,
		rehabRepo:     rehabRepo,
		nutritionRepo: nutritionRepo,
	}
}

func (s *RecommendationService) GenerateRecommendations(ctx context.Context, profile *models.UserProfile) (*Recommendation, error) {
	goal := primaryGoalFromProfile(profile.Goal)
	format := ""
	if profile.TrainingAccess != nil {
		format = *profile.TrainingAccess
	}
	level := profile.FitnessLevel

	rec := &Recommendation{}

	// Find a matching training program
	programs, err := s.programRepo.ListPrograms(ctx, format, goal, level)
	if err == nil && len(programs) > 0 {
		rec.ProgramID = &programs[0].ID
		rec.ProgramName = programs[0].Name
	}

	// If the user reports pain or has diagnoses, recommend a rehab course
	if profile.HasPain || len(profile.Diagnoses) > 0 {
		category := rehabCategoryFromProfile(profile)
		courses, err := s.rehabRepo.ListCourses(ctx, category)
		if err == nil && len(courses) > 0 {
			rec.RehabCourseID = &courses[0].ID
			rec.RehabCourseName = courses[0].Name
		}
	}

	// Calculate macronutrient targets
	macros := CalculateMacroTargets(profile.Gender, profile.WeightKg, profile.HeightCm, profile.Age, goal, level)
	rec.Calories = macros.Calories
	rec.Protein = macros.Protein
	rec.Fat = macros.Fat
	rec.Carbs = macros.Carbs

	// Build personalized message
	rec.Message = buildRecommendationMessage(rec, goal)

	return rec, nil
}

// rehabCategoryFromProfile determines the rehab category from the user's pain
// locations and diagnoses. It returns the first pain location if available,
// otherwise falls back to the first diagnosis.
func rehabCategoryFromProfile(profile *models.UserProfile) string {
	if len(profile.PainLocations) > 0 {
		return profile.PainLocations[0]
	}
	if len(profile.Diagnoses) > 0 {
		return profile.Diagnoses[0]
	}
	return ""
}

// primaryGoalFromProfile extracts the first goal from a comma-separated goal string.
func primaryGoalFromProfile(goal string) string {
	if goal == "" {
		return ""
	}
	idx := strings.Index(goal, ",")
	if idx == -1 {
		return strings.TrimSpace(goal)
	}
	return strings.TrimSpace(goal[:idx])
}

// buildRecommendationMessage creates a human-readable recommendation summary
// in Russian.
func buildRecommendationMessage(rec *Recommendation, goal string) string {
	var parts []string

	if rec.ProgramName != "" {
		parts = append(parts, fmt.Sprintf("Рекомендуемая программа: %s.", rec.ProgramName))
	}

	if rec.RehabCourseName != "" {
		parts = append(parts, fmt.Sprintf("Курс реабилитации: %s.", rec.RehabCourseName))
	}

	goalLabel := goalDisplayName(goal)
	parts = append(parts, fmt.Sprintf(
		"Для цели \"%s\" рекомендуем %d ккал в день (Б: %dг, Ж: %dг, У: %dг).",
		goalLabel, rec.Calories, rec.Protein, rec.Fat, rec.Carbs,
	))

	return strings.Join(parts, " ")
}

// goalDisplayName returns a user-friendly Russian label for a goal slug.
func goalDisplayName(goal string) string {
	switch goal {
	case "weight_loss":
		return "Похудение"
	case "muscle_gain":
		return "Набор мышечной массы"
	case "maintenance":
		return "Поддержание формы"
	default:
		return "Общая физическая подготовка"
	}
}
