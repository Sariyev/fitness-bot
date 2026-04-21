package webapp

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
)

type AdminHandler struct {
	userSvc      *service.UserService
	workoutSvc   *service.WorkoutService
	nutritionSvc *service.NutritionService
	scoreSvc     *service.ScoreService
}

func NewAdminHandler(
	userSvc *service.UserService,
	workoutSvc *service.WorkoutService,
	nutritionSvc *service.NutritionService,
	scoreSvc *service.ScoreService,
) *AdminHandler {
	return &AdminHandler{
		userSvc:      userSvc,
		workoutSvc:   workoutSvc,
		nutritionSvc: nutritionSvc,
		scoreSvc:     scoreSvc,
	}
}

// ==================== USERS ====================

func (h *AdminHandler) HandleUserRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/app/api/admin/users")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		if r.Method == http.MethodGet {
			h.ListUsers(w, r)
			return
		}
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetUser(w, r, id)
	case http.MethodPut:
		h.UpdateUser(w, r, id)
	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	users, total, err := h.userSvc.ListUsers(r.Context(), limit, offset)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"users": users,
		"total": total,
	})
}

func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request, id int64) {
	user, err := h.userSvc.GetUserByID(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "user not found")
		return
	}

	profile, _ := h.userSvc.GetProfile(r.Context(), user.ID)

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"user":    user,
		"profile": profile,
	})
}

type updateUserRequest struct {
	Role   *string `json:"role"`
	IsPaid *bool   `json:"is_paid"`
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request, id int64) {
	user, err := h.userSvc.GetUserByID(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "user not found")
		return
	}

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Role != nil {
		if *req.Role != "client" && *req.Role != "admin" {
			jsonError(w, http.StatusBadRequest, "role must be 'client' or 'admin'")
			return
		}
		user.Role = *req.Role
	}
	if req.IsPaid != nil {
		user.IsPaid = *req.IsPaid
	}

	if err := h.userSvc.UpdateUser(r.Context(), user); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]bool{"success": true})
}

// ==================== PROGRAMS ====================

func (h *AdminHandler) HandleProgramRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/app/api/admin/programs")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.listPrograms(w, r)
		case http.MethodPost:
			h.createProgram(w, r)
		default:
			jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid program id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getProgram(w, r, id)
	case http.MethodPut:
		h.updateProgram(w, r, id)
	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *AdminHandler) listPrograms(w http.ResponseWriter, r *http.Request) {
	programs, err := h.workoutSvc.ListAllPrograms(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list programs")
		return
	}
	jsonResponse(w, http.StatusOK, programs)
}

func (h *AdminHandler) getProgram(w http.ResponseWriter, r *http.Request, id int) {
	p, err := h.workoutSvc.GetProgram(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "program not found")
		return
	}
	jsonResponse(w, http.StatusOK, p)
}

type programRequest struct {
	Slug          string `json:"slug"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Goal          string `json:"goal"`
	Format        string `json:"format"`
	Level         string `json:"level"`
	DurationWeeks int    `json:"duration_weeks"`
	IsActive      bool   `json:"is_active"`
	SortOrder     int    `json:"sort_order"`
}

func (h *AdminHandler) createProgram(w http.ResponseWriter, r *http.Request) {
	var req programRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		jsonError(w, http.StatusBadRequest, "name is required")
		return
	}

	p := &models.Program{
		Slug: req.Slug, Name: req.Name, Description: req.Description,
		Goal: req.Goal, Format: req.Format, Level: req.Level,
		DurationWeeks: req.DurationWeeks, IsActive: req.IsActive, SortOrder: req.SortOrder,
	}
	if err := h.workoutSvc.CreateProgram(r.Context(), p); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to create program")
		return
	}
	jsonResponse(w, http.StatusCreated, p)
}

func (h *AdminHandler) updateProgram(w http.ResponseWriter, r *http.Request, id int) {
	p, err := h.workoutSvc.GetProgram(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "program not found")
		return
	}

	var req programRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p.Slug = req.Slug
	p.Name = req.Name
	p.Description = req.Description
	p.Goal = req.Goal
	p.Format = req.Format
	p.Level = req.Level
	p.DurationWeeks = req.DurationWeeks
	p.IsActive = req.IsActive
	p.SortOrder = req.SortOrder

	if err := h.workoutSvc.UpdateProgram(r.Context(), p); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to update program")
		return
	}
	jsonResponse(w, http.StatusOK, p)
}

// ==================== WORKOUTS ====================

func (h *AdminHandler) HandleWorkoutRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/app/api/admin/workouts")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.listWorkouts(w, r)
		case http.MethodPost:
			h.createWorkout(w, r)
		default:
			jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid workout id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getWorkout(w, r, id)
	case http.MethodPut:
		h.updateWorkout(w, r, id)
	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *AdminHandler) listWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.workoutSvc.ListAllWorkouts(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list workouts")
		return
	}
	jsonResponse(w, http.StatusOK, workouts)
}

func (h *AdminHandler) getWorkout(w http.ResponseWriter, r *http.Request, id int) {
	wo, err := h.workoutSvc.GetWorkout(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "workout not found")
		return
	}

	exercises, _ := h.workoutSvc.GetWorkoutExercisesWithDetails(r.Context(), id)

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"workout":   wo,
		"exercises": exercises,
	})
}

type workoutRequest struct {
	ProgramID       *int     `json:"program_id"`
	Slug            string   `json:"slug"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Goal            string   `json:"goal"`
	Format          string   `json:"format"`
	Level           string   `json:"level"`
	DurationMinutes int      `json:"duration_minutes"`
	Equipment       []string `json:"equipment"`
	ExpectedResult  string   `json:"expected_result"`
	VideoURL        string   `json:"video_url"`
	SortOrder       int      `json:"sort_order"`
	WeekNumber      *int     `json:"week_number"`
	DayNumber       *int     `json:"day_number"`
	IsActive        bool     `json:"is_active"`
}

func (h *AdminHandler) createWorkout(w http.ResponseWriter, r *http.Request) {
	var req workoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		jsonError(w, http.StatusBadRequest, "name is required")
		return
	}

	wo := &models.Workout{
		ProgramID: req.ProgramID, Slug: req.Slug, Name: req.Name,
		Description: req.Description, Goal: req.Goal, Format: req.Format,
		Level: req.Level, DurationMinutes: req.DurationMinutes,
		Equipment: req.Equipment, ExpectedResult: req.ExpectedResult,
		VideoURL: req.VideoURL, SortOrder: req.SortOrder,
		WeekNumber: req.WeekNumber, DayNumber: req.DayNumber, IsActive: req.IsActive,
	}
	if err := h.workoutSvc.CreateWorkout(r.Context(), wo); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to create workout")
		return
	}
	jsonResponse(w, http.StatusCreated, wo)
}

func (h *AdminHandler) updateWorkout(w http.ResponseWriter, r *http.Request, id int) {
	wo, err := h.workoutSvc.GetWorkout(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "workout not found")
		return
	}

	var req workoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	wo.ProgramID = req.ProgramID
	wo.Slug = req.Slug
	wo.Name = req.Name
	wo.Description = req.Description
	wo.Goal = req.Goal
	wo.Format = req.Format
	wo.Level = req.Level
	wo.DurationMinutes = req.DurationMinutes
	wo.Equipment = req.Equipment
	wo.ExpectedResult = req.ExpectedResult
	wo.VideoURL = req.VideoURL
	wo.SortOrder = req.SortOrder
	wo.WeekNumber = req.WeekNumber
	wo.DayNumber = req.DayNumber
	wo.IsActive = req.IsActive

	if err := h.workoutSvc.UpdateWorkout(r.Context(), wo); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to update workout")
		return
	}
	jsonResponse(w, http.StatusOK, wo)
}

// ==================== EXERCISES ====================

func (h *AdminHandler) HandleExerciseRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/app/api/admin/exercises")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.listExercises(w, r)
		case http.MethodPost:
			h.createExercise(w, r)
		default:
			jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid exercise id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getExercise(w, r, id)
	case http.MethodPut:
		h.updateExercise(w, r, id)
	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *AdminHandler) listExercises(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.workoutSvc.ListExercises(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list exercises")
		return
	}
	jsonResponse(w, http.StatusOK, exercises)
}

func (h *AdminHandler) getExercise(w http.ResponseWriter, r *http.Request, id int) {
	e, err := h.workoutSvc.GetExercise(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "exercise not found")
		return
	}
	jsonResponse(w, http.StatusOK, e)
}

type exerciseRequest struct {
	Name               string `json:"name"`
	Technique          string `json:"technique"`
	CommonMistakes     string `json:"common_mistakes"`
	EasierModification string `json:"easier_modification"`
	HarderModification string `json:"harder_modification"`
	RestSeconds        int    `json:"rest_seconds"`
}

func (h *AdminHandler) createExercise(w http.ResponseWriter, r *http.Request) {
	var req exerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		jsonError(w, http.StatusBadRequest, "name is required")
		return
	}

	e := &models.Exercise{
		Name: req.Name, Technique: req.Technique,
		CommonMistakes: req.CommonMistakes, EasierModification: req.EasierModification,
		HarderModification: req.HarderModification, RestSeconds: req.RestSeconds,
	}
	if err := h.workoutSvc.CreateExercise(r.Context(), e); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to create exercise")
		return
	}
	jsonResponse(w, http.StatusCreated, e)
}

func (h *AdminHandler) updateExercise(w http.ResponseWriter, r *http.Request, id int) {
	e, err := h.workoutSvc.GetExercise(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "exercise not found")
		return
	}

	var req exerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	e.Name = req.Name
	e.Technique = req.Technique
	e.CommonMistakes = req.CommonMistakes
	e.EasierModification = req.EasierModification
	e.HarderModification = req.HarderModification
	e.RestSeconds = req.RestSeconds

	if err := h.workoutSvc.UpdateExercise(r.Context(), e); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to update exercise")
		return
	}
	jsonResponse(w, http.StatusOK, e)
}

// ==================== WORKOUT EXERCISES ====================

func (h *AdminHandler) HandleWorkoutExerciseRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		WorkoutID       int    `json:"workout_id"`
		ExerciseID      int    `json:"exercise_id"`
		Sets            int    `json:"sets"`
		Reps            string `json:"reps"`
		DurationSeconds int    `json:"duration_seconds"`
		SortOrder       int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	we := &models.WorkoutExercise{
		WorkoutID: req.WorkoutID, ExerciseID: req.ExerciseID,
		Sets: req.Sets, Reps: req.Reps,
		DurationSeconds: req.DurationSeconds, SortOrder: req.SortOrder,
	}
	if err := h.workoutSvc.AddExerciseToWorkout(r.Context(), we); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to add exercise")
		return
	}
	jsonResponse(w, http.StatusCreated, we)
}

// ==================== MEAL PLANS ====================

func (h *AdminHandler) HandleMealPlanRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/app/api/admin/meal-plans")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.listMealPlans(w, r)
		case http.MethodPost:
			h.createMealPlan(w, r)
		default:
			jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid meal plan id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getMealPlan(w, r, id)
	case http.MethodPut:
		h.updateMealPlan(w, r, id)
	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *AdminHandler) listMealPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.nutritionSvc.ListAllPlans(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list meal plans")
		return
	}
	jsonResponse(w, http.StatusOK, plans)
}

func (h *AdminHandler) getMealPlan(w http.ResponseWriter, r *http.Request, id int) {
	plan, err := h.nutritionSvc.GetPlan(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "meal plan not found")
		return
	}

	meals, _ := h.nutritionSvc.GetPlanMeals(r.Context(), id)

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"plan":  plan,
		"meals": meals,
	})
}

type mealPlanRequest struct {
	Slug      string  `json:"slug"`
	Name      string  `json:"name"`
	Goal      string  `json:"goal"`
	DayNumber int     `json:"day_number"`
	Calories  int     `json:"calories"`
	Protein   float64 `json:"protein"`
	Fat       float64 `json:"fat"`
	Carbs     float64 `json:"carbs"`
	IsActive  bool    `json:"is_active"`
	SortOrder int     `json:"sort_order"`
}

func (h *AdminHandler) createMealPlan(w http.ResponseWriter, r *http.Request) {
	var req mealPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		jsonError(w, http.StatusBadRequest, "name is required")
		return
	}

	p := &models.MealPlan{
		Slug: req.Slug, Name: req.Name, Goal: req.Goal,
		DayNumber: req.DayNumber, Calories: req.Calories,
		Protein: req.Protein, Fat: req.Fat, Carbs: req.Carbs,
		IsActive: req.IsActive, SortOrder: req.SortOrder,
	}
	if err := h.nutritionSvc.CreatePlan(r.Context(), p); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to create meal plan")
		return
	}
	jsonResponse(w, http.StatusCreated, p)
}

func (h *AdminHandler) updateMealPlan(w http.ResponseWriter, r *http.Request, id int) {
	p, err := h.nutritionSvc.GetPlan(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "meal plan not found")
		return
	}

	var req mealPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p.Slug = req.Slug
	p.Name = req.Name
	p.Goal = req.Goal
	p.DayNumber = req.DayNumber
	p.Calories = req.Calories
	p.Protein = req.Protein
	p.Fat = req.Fat
	p.Carbs = req.Carbs
	p.IsActive = req.IsActive
	p.SortOrder = req.SortOrder

	if err := h.nutritionSvc.UpdatePlan(r.Context(), p); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to update meal plan")
		return
	}
	jsonResponse(w, http.StatusOK, p)
}

// ==================== MEALS ====================

func (h *AdminHandler) HandleMealRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/app/api/admin/meals")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			planID, _ := strconv.Atoi(r.URL.Query().Get("plan_id"))
			if planID == 0 {
				jsonError(w, http.StatusBadRequest, "plan_id is required")
				return
			}
			h.listMeals(w, r, planID)
		case http.MethodPost:
			h.createMeal(w, r)
		default:
			jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid meal id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getMeal(w, r, id)
	case http.MethodPut:
		h.updateMeal(w, r, id)
	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *AdminHandler) listMeals(w http.ResponseWriter, r *http.Request, planID int) {
	meals, err := h.nutritionSvc.GetPlanMeals(r.Context(), planID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list meals")
		return
	}
	jsonResponse(w, http.StatusOK, meals)
}

func (h *AdminHandler) getMeal(w http.ResponseWriter, r *http.Request, id int) {
	m, err := h.nutritionSvc.GetMeal(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "meal not found")
		return
	}
	jsonResponse(w, http.StatusOK, m)
}

type mealRequest struct {
	MealPlanID   int     `json:"meal_plan_id"`
	MealType     string  `json:"meal_type"`
	Name         string  `json:"name"`
	Recipe       string  `json:"recipe"`
	Calories     int     `json:"calories"`
	Protein      float64 `json:"protein"`
	Fat          float64 `json:"fat"`
	Carbs        float64 `json:"carbs"`
	Alternatives string  `json:"alternatives"`
	SortOrder    int     `json:"sort_order"`
}

func (h *AdminHandler) createMeal(w http.ResponseWriter, r *http.Request) {
	var req mealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		jsonError(w, http.StatusBadRequest, "name is required")
		return
	}

	m := &models.Meal{
		MealPlanID: req.MealPlanID, MealType: req.MealType, Name: req.Name,
		Recipe: req.Recipe, Calories: req.Calories,
		Protein: req.Protein, Fat: req.Fat, Carbs: req.Carbs,
		Alternatives: req.Alternatives, SortOrder: req.SortOrder,
	}
	if err := h.nutritionSvc.CreateMeal(r.Context(), m); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to create meal")
		return
	}
	jsonResponse(w, http.StatusCreated, m)
}

func (h *AdminHandler) updateMeal(w http.ResponseWriter, r *http.Request, id int) {
	m, err := h.nutritionSvc.GetMeal(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "meal not found")
		return
	}

	var req mealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	m.MealPlanID = req.MealPlanID
	m.MealType = req.MealType
	m.Name = req.Name
	m.Recipe = req.Recipe
	m.Calories = req.Calories
	m.Protein = req.Protein
	m.Fat = req.Fat
	m.Carbs = req.Carbs
	m.Alternatives = req.Alternatives
	m.SortOrder = req.SortOrder

	if err := h.nutritionSvc.UpdateMeal(r.Context(), m); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to update meal")
		return
	}
	jsonResponse(w, http.StatusOK, m)
}

// ==================== REVIEWS & STATS ====================

func (h *AdminHandler) GetReviewsSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	summary, err := h.scoreSvc.GetBotSummary(r.Context())
	if err != nil {
		jsonResponse(w, http.StatusOK, map[string]interface{}{
			"average_score": 0,
			"total_reviews": 0,
		})
		return
	}

	jsonResponse(w, http.StatusOK, summary)
}

func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	_, total, err := h.userSvc.ListUsers(r.Context(), 1, 0)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to get stats")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"total_users": total,
	})
}
