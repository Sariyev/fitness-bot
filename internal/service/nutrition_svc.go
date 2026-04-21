package service

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
	"math"
)

type MacroTargets struct {
	Calories int `json:"calories"`
	Protein  int `json:"protein"`
	Fat      int `json:"fat"`
	Carbs    int `json:"carbs"`
}

type NutritionService struct {
	nutritionRepo repository.NutritionRepository
	foodLogRepo   repository.FoodLogRepository
}

func NewNutritionService(nutritionRepo repository.NutritionRepository, foodLogRepo repository.FoodLogRepository) *NutritionService {
	return &NutritionService{
		nutritionRepo: nutritionRepo,
		foodLogRepo:   foodLogRepo,
	}
}

func (s *NutritionService) ListPlans(ctx context.Context, goal string) ([]models.MealPlan, error) {
	return s.nutritionRepo.ListPlans(ctx, goal)
}

func (s *NutritionService) GetPlan(ctx context.Context, id int) (*models.MealPlan, error) {
	return s.nutritionRepo.GetPlanByID(ctx, id)
}

func (s *NutritionService) GetPlanMeals(ctx context.Context, planID int) ([]models.Meal, error) {
	return s.nutritionRepo.ListMeals(ctx, planID)
}

// CalculateTargets computes daily calorie and macronutrient targets using the
// Mifflin-St Jeor equation. The level parameter determines the activity
// multiplier: "beginner" -> 1.375, "intermediate" -> 1.55, "advanced" -> 1.725.
func (s *NutritionService) CalculateTargets(gender string, weightKg float64, heightCm, age int, goal, level string) MacroTargets {
	return CalculateMacroTargets(gender, weightKg, heightCm, age, goal, level)
}

// CalculateMacroTargets is a package-level helper so other services can reuse
// the Mifflin-St Jeor calculation without needing a NutritionService instance.
func CalculateMacroTargets(gender string, weightKg float64, heightCm, age int, goal, level string) MacroTargets {
	// Mifflin-St Jeor BMR
	bmr := 10*weightKg + 6.25*float64(heightCm) - 5*float64(age)
	if gender == "male" {
		bmr += 5
	} else {
		bmr -= 161
	}

	// Activity multiplier based on fitness level
	var multiplier float64
	switch level {
	case "beginner":
		multiplier = 1.375
	case "intermediate":
		multiplier = 1.55
	case "advanced":
		multiplier = 1.725
	default:
		multiplier = 1.55
	}

	tdee := bmr * multiplier

	// Adjust calories based on goal
	var calories float64
	switch goal {
	case "weight_loss":
		calories = tdee - 500
	case "muscle_gain":
		calories = tdee + 300
	default:
		calories = tdee
	}

	// Protein per kg based on goal
	var proteinPerKg float64
	switch goal {
	case "weight_loss":
		proteinPerKg = 1.8
	case "muscle_gain":
		proteinPerKg = 2.0
	default:
		proteinPerKg = 1.4
	}

	proteinGrams := proteinPerKg * weightKg
	fatGrams := (calories * 0.25) / 9
	carbGrams := (calories - proteinGrams*4 - fatGrams*9) / 4

	return MacroTargets{
		Calories: int(math.Round(calories)),
		Protein:  int(math.Round(proteinGrams)),
		Fat:      int(math.Round(fatGrams)),
		Carbs:    int(math.Round(carbGrams)),
	}
}

func (s *NutritionService) ListAllPlans(ctx context.Context) ([]models.MealPlan, error) {
	return s.nutritionRepo.ListAllPlans(ctx)
}

func (s *NutritionService) CreatePlan(ctx context.Context, p *models.MealPlan) error {
	return s.nutritionRepo.CreatePlan(ctx, p)
}

func (s *NutritionService) UpdatePlan(ctx context.Context, p *models.MealPlan) error {
	return s.nutritionRepo.UpdatePlan(ctx, p)
}

func (s *NutritionService) CreateMeal(ctx context.Context, m *models.Meal) error {
	return s.nutritionRepo.CreateMeal(ctx, m)
}

func (s *NutritionService) UpdateMeal(ctx context.Context, m *models.Meal) error {
	return s.nutritionRepo.UpdateMeal(ctx, m)
}

func (s *NutritionService) GetMeal(ctx context.Context, id int) (*models.Meal, error) {
	return s.nutritionRepo.GetMealByID(ctx, id)
}

func (s *NutritionService) GetFoodLog(ctx context.Context, userID int64, date string) ([]models.FoodLogEntry, error) {
	return s.foodLogRepo.ListByDate(ctx, userID, date)
}

func (s *NutritionService) AddFoodLog(ctx context.Context, entry *models.FoodLogEntry) error {
	return s.foodLogRepo.Create(ctx, entry)
}

func (s *NutritionService) DeleteFoodLog(ctx context.Context, userID int64, id int64) error {
	return s.foodLogRepo.Delete(ctx, userID, id)
}

func (s *NutritionService) GetDailySummary(ctx context.Context, userID int64, date string) (calories int, protein, fat, carbs float64, err error) {
	return s.foodLogRepo.GetDailySummary(ctx, userID, date)
}
