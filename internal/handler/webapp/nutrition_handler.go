package webapp

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
)

type NutritionHandler struct {
	nutritionSvc *service.NutritionService
	accessSvc    *service.AccessService
	mediaSvc     *service.MediaService // optional; nil when R2 not configured
}

func NewNutritionHandler(nutritionSvc *service.NutritionService, accessSvc *service.AccessService, mediaSvc *service.MediaService) *NutritionHandler {
	return &NutritionHandler{nutritionSvc: nutritionSvc, accessSvc: accessSvc, mediaSvc: mediaSvc}
}

// fillImageURL resolves image_media_id → presigned URL and sets ImageURL on the plan.
func (h *NutritionHandler) fillPlanImageURL(r *http.Request, p *models.MealPlan) {
	if p == nil || p.ImageMediaID == nil || h.mediaSvc == nil {
		return
	}
	if u, err := h.mediaSvc.PresignReadURL(r.Context(), *p.ImageMediaID); err == nil {
		p.ImageURL = u
	}
}

func (h *NutritionHandler) fillMealImageURL(r *http.Request, m *models.Meal) {
	if m == nil || m.ImageMediaID == nil || h.mediaSvc == nil {
		return
	}
	if u, err := h.mediaSvc.PresignReadURL(r.Context(), *m.ImageMediaID); err == nil {
		m.ImageURL = u
	}
}

type AddFoodLogRequest struct {
	Date     string  `json:"date"`
	MealType string  `json:"meal_type"`
	FoodName string  `json:"food_name"`
	Calories int     `json:"calories"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
}

// HandleNutritionRoutes dispatches /app/api/nutrition/... requests.
//
//	GET /app/api/nutrition/plans          -> list plans, optional ?goal= filter
//	GET /app/api/nutrition/plans/{id}     -> get plan with meals
//	GET /app/api/nutrition/calculator     -> calculate macros from query params
func (h *NutritionHandler) HandleNutritionRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Handle /app/api/nutrition/calculator
	if strings.HasPrefix(r.URL.Path, "/app/api/nutrition/calculator") {
		h.CalculateMacros(w, r)
		return
	}

	// Handle /app/api/nutrition/plans/...
	path := strings.TrimPrefix(r.URL.Path, "/app/api/nutrition/plans")
	path = strings.TrimPrefix(path, "/")

	// GET /app/api/nutrition/plans
	if path == "" {
		h.ListPlans(w, r)
		return
	}

	// GET /app/api/nutrition/plans/{id}
	h.GetPlan(w, r, path)
}

// HandleFoodLogRoutes dispatches /app/api/food-log/... requests.
//
//	GET    /app/api/food-log?date=YYYY-MM-DD -> list entries for date
//	POST   /app/api/food-log                 -> add entry
//	DELETE /app/api/food-log/{id}            -> delete entry
func (h *NutritionHandler) HandleFoodLogRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/app/api/food-log")
	path = strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		// Do not handle /app/api/food-log/summary here
		if strings.HasPrefix(path, "summary") {
			jsonError(w, http.StatusNotFound, "not found")
			return
		}
		h.ListFoodLog(w, r)

	case http.MethodPost:
		if path != "" {
			jsonError(w, http.StatusNotFound, "not found")
			return
		}
		h.AddFoodLogEntry(w, r)

	case http.MethodDelete:
		if path == "" {
			jsonError(w, http.StatusBadRequest, "missing entry id")
			return
		}
		h.DeleteFoodLogEntry(w, r, path)

	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// HandleFoodLogSummary handles GET /app/api/food-log/summary?date=YYYY-MM-DD
func (h *NutritionHandler) HandleFoodLogSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		jsonError(w, http.StatusBadRequest, "date parameter required")
		return
	}

	calories, protein, fat, carbs, err := h.nutritionSvc.GetDailySummary(r.Context(), user.ID, date)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to load summary")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"date":     date,
		"calories": calories,
		"protein":  protein,
		"fat":      fat,
		"carbs":    carbs,
	})
}

// ListPlans handles GET /app/api/nutrition/plans?goal=
// Each plan is annotated with `locked` based on viewing user's access.
func (h *NutritionHandler) ListPlans(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	goal := r.URL.Query().Get("goal")

	plans, err := h.nutritionSvc.ListPlans(r.Context(), goal)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list plans")
		return
	}

	for i := range plans {
		can, _ := h.accessSvc.CanAccess(r.Context(), user, plans[i].AccessTier, models.CategoryNutrition)
		plans[i].Locked = !can
		h.fillPlanImageURL(r, &plans[i])
	}

	jsonResponse(w, http.StatusOK, plans)
}

// GetPlan handles GET /app/api/nutrition/plans/{id}
// Returns 402 if the plan is locked.
func (h *NutritionHandler) GetPlan(w http.ResponseWriter, r *http.Request, idStr string) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.Atoi(strings.TrimSuffix(idStr, "/"))
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid plan id")
		return
	}

	plan, err := h.nutritionSvc.GetPlan(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "plan not found")
		return
	}

	can, err := h.accessSvc.CanAccess(r.Context(), user, plan.AccessTier, models.CategoryNutrition)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "access check failed")
		return
	}
	if !can {
		price, _ := h.accessSvc.GetPrice(r.Context(), models.CategoryNutrition)
		jsonResponse(w, http.StatusPaymentRequired, map[string]interface{}{
			"error":     "locked",
			"category":  models.CategoryNutrition,
			"tier":      plan.AccessTier,
			"price_kzt": price,
		})
		return
	}
	plan.Locked = false
	h.fillPlanImageURL(r, plan)

	meals, err := h.nutritionSvc.GetPlanMeals(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to load meals")
		return
	}
	for i := range meals {
		h.fillMealImageURL(r, &meals[i])
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"plan":  plan,
		"meals": meals,
	})
}

// CalculateMacros handles GET /app/api/nutrition/calculator?gender=&weight_kg=&height_cm=&age=&goal=
func (h *NutritionHandler) CalculateMacros(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	gender := q.Get("gender")
	if gender == "" {
		jsonError(w, http.StatusBadRequest, "gender parameter required")
		return
	}

	weightKg, err := strconv.ParseFloat(q.Get("weight_kg"), 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid weight_kg")
		return
	}

	heightCm, err := strconv.Atoi(q.Get("height_cm"))
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid height_cm")
		return
	}

	age, err := strconv.Atoi(q.Get("age"))
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid age")
		return
	}

	goal := q.Get("goal")
	level := q.Get("level")

	result := h.nutritionSvc.CalculateTargets(gender, weightKg, heightCm, age, goal, level)
	jsonResponse(w, http.StatusOK, result)
}

// ListFoodLog handles GET /app/api/food-log?date=YYYY-MM-DD
func (h *NutritionHandler) ListFoodLog(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		jsonError(w, http.StatusBadRequest, "date parameter required")
		return
	}

	entries, err := h.nutritionSvc.GetFoodLog(r.Context(), user.ID, date)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list food log")
		return
	}

	jsonResponse(w, http.StatusOK, entries)
}

// AddFoodLogEntry handles POST /app/api/food-log
func (h *NutritionHandler) AddFoodLogEntry(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req AddFoodLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Date == "" {
		jsonError(w, http.StatusBadRequest, "date is required")
		return
	}
	if req.FoodName == "" {
		jsonError(w, http.StatusBadRequest, "food_name is required")
		return
	}

	entry := &models.FoodLogEntry{
		UserID:   user.ID,
		Date:     req.Date,
		MealType: req.MealType,
		FoodName: req.FoodName,
		Calories: req.Calories,
		Protein:  req.Protein,
		Fat:      req.Fat,
		Carbs:    req.Carbs,
	}

	if err := h.nutritionSvc.AddFoodLog(r.Context(), entry); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to add food log entry")
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]bool{"success": true})
}

// DeleteFoodLogEntry handles DELETE /app/api/food-log/{id}
func (h *NutritionHandler) DeleteFoodLogEntry(w http.ResponseWriter, r *http.Request, idStr string) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseInt(strings.TrimSuffix(idStr, "/"), 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid entry id")
		return
	}

	if err := h.nutritionSvc.DeleteFoodLog(r.Context(), user.ID, id); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to delete food log entry")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]bool{"success": true})
}
