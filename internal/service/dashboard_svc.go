package service

import (
	"context"
	"fitness-bot/internal/models"
	"fmt"
	"time"
)

type DashboardItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Done  bool   `json:"done"`
}

type DashboardData struct {
	Greeting       string         `json:"greeting"`
	TrainerMessage string         `json:"trainer_message"`
	TodayWorkout   *DashboardItem `json:"today_workout"`
	TodayMeal      *DashboardItem `json:"today_meal"`
	TodayRehab     *DashboardItem `json:"today_rehab"`
	Goals          []string       `json:"goals"`
	CurrentStreak  int            `json:"current_streak"`
}

type DashboardService struct {
	userSvc      *UserService
	workoutSvc   *WorkoutService
	rehabSvc     *RehabService
	nutritionSvc *NutritionService
}

func NewDashboardService(
	userSvc *UserService,
	workoutSvc *WorkoutService,
	rehabSvc *RehabService,
	nutritionSvc *NutritionService,
) *DashboardService {
	return &DashboardService{
		userSvc:      userSvc,
		workoutSvc:   workoutSvc,
		rehabSvc:     rehabSvc,
		nutritionSvc: nutritionSvc,
	}
}

func (s *DashboardService) GetDashboard(ctx context.Context, user *models.User) (*DashboardData, error) {
	dashboard := &DashboardData{
		Greeting:       fmt.Sprintf("Привет, %s!", user.FirstName),
		TrainerMessage: buildTrainerMessage(),
	}

	// Load user profile for goals
	profile, err := s.userSvc.GetProfile(ctx, user.TelegramID)
	if err == nil && profile != nil {
		dashboard.Goals = parseGoals(profile.Goal)
	}

	// Determine today's workout from active program enrollment
	enrollment, err := s.workoutSvc.GetActiveEnrollment(ctx, user.ID)
	if err == nil && enrollment != nil {
		todayWorkout := findTodayWorkout(ctx, s.workoutSvc, enrollment)
		if todayWorkout != nil {
			dashboard.TodayWorkout = todayWorkout
		}
	}

	// Check for rehab sessions if the user has pain
	if profile != nil && profile.HasPain && len(profile.PainLocations) > 0 {
		courses, err := s.rehabSvc.ListCourses(ctx, profile.PainLocations[0])
		if err == nil && len(courses) > 0 {
			dashboard.TodayRehab = &DashboardItem{
				ID:    courses[0].ID,
				Title: courses[0].Name,
				Type:  "rehab",
				Done:  false,
			}
		}
	}

	// Find a meal plan matching the user's goal
	if profile != nil {
		goal := primaryGoal(profile.Goal)
		plans, err := s.nutritionSvc.ListPlans(ctx, goal)
		if err == nil && len(plans) > 0 {
			dashboard.TodayMeal = &DashboardItem{
				ID:    plans[0].ID,
				Title: plans[0].Name,
				Type:  "meal",
				Done:  false,
			}
		}
	}

	// Current streak from the workout service's completion repo
	current, _, err := s.workoutSvc.completionRepo.GetStreak(ctx, user.ID)
	if err == nil {
		dashboard.CurrentStreak = current
	}

	return dashboard, nil
}

// findTodayWorkout locates the workout scheduled for today based on the
// enrollment start date, current week, and day of week.
func findTodayWorkout(ctx context.Context, svc *WorkoutService, enrollment *models.UserProgramEnrollment) *DashboardItem {
	workouts, err := svc.workoutRepo.ListByProgram(ctx, enrollment.ProgramID)
	if err != nil || len(workouts) == 0 {
		return nil
	}

	daysSinceStart := int(time.Since(enrollment.StartedAt).Hours() / 24)
	currentWeek := daysSinceStart/7 + 1
	currentDay := daysSinceStart%7 + 1

	for _, w := range workouts {
		if w.WeekNumber != nil && w.DayNumber != nil {
			if *w.WeekNumber == currentWeek && *w.DayNumber == currentDay {
				return &DashboardItem{
					ID:    w.ID,
					Title: w.Name,
					Type:  "workout",
					Done:  false,
				}
			}
		}
	}

	// Fall back to the first workout if no schedule match
	return &DashboardItem{
		ID:    workouts[0].ID,
		Title: workouts[0].Name,
		Type:  "workout",
		Done:  false,
	}
}

// buildTrainerMessage returns a motivational message from the virtual coach.
func buildTrainerMessage() string {
	hour := time.Now().Hour()
	switch {
	case hour < 12:
		return "Доброе утро! Отличный день для тренировки. Давайте начнём!"
	case hour < 18:
		return "Добрый день! Не забудьте про тренировку сегодня. Вы на верном пути!"
	default:
		return "Добрый вечер! Ещё не поздно позаниматься. Каждый шаг приближает вас к цели!"
	}
}

// parseGoals splits a comma-separated goal string into a slice.
func parseGoals(goal string) []string {
	if goal == "" {
		return nil
	}
	var goals []string
	for _, g := range splitComma(goal) {
		if g != "" {
			goals = append(goals, g)
		}
	}
	return goals
}

// splitComma is a minimal comma splitter that trims whitespace.
func splitComma(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			part := trimSpace(s[start:i])
			parts = append(parts, part)
			start = i + 1
		}
	}
	parts = append(parts, trimSpace(s[start:]))
	return parts
}

// trimSpace removes leading and trailing whitespace from s.
func trimSpace(s string) string {
	i, j := 0, len(s)
	for i < j && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	for j > i && (s[j-1] == ' ' || s[j-1] == '\t') {
		j--
	}
	return s[i:j]
}

// primaryGoal returns the first goal from a comma-separated goal string.
func primaryGoal(goal string) string {
	goals := parseGoals(goal)
	if len(goals) > 0 {
		return goals[0]
	}
	return ""
}
