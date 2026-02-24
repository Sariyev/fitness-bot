package webapp

import (
	"encoding/json"
	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
	"net/http"
)

type RegistrationHandler struct {
	userSvc *service.UserService
}

func NewRegistrationHandler(userSvc *service.UserService) *RegistrationHandler {
	return &RegistrationHandler{userSvc: userSvc}
}

type RegisterRequest struct {
	Age          int      `json:"age"`
	HeightCm     int      `json:"height_cm"`
	WeightKg     float64  `json:"weight_kg"`
	Gender       string   `json:"gender"`
	FitnessLevel string   `json:"fitness_level"`
	Goals        []string `json:"goals"`
}

var validGoals = map[string]bool{
	"weight_loss": true, "muscle_gain": true, "strength": true,
	"endurance": true, "maintenance": true, "hernia": true,
	"protrusion": true, "scoliosis": true, "kyphosis": true, "lordosis": true,
}

var validFitness = map[string]bool{
	"beginner": true, "intermediate": true, "advanced": true,
}

var validGenders = map[string]bool{"male": true, "female": true}

// POST /app/api/register
func (h *RegistrationHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if user.IsRegistered {
		jsonError(w, http.StatusConflict, "already registered")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate
	if req.Age < 10 || req.Age > 120 {
		jsonError(w, http.StatusBadRequest, "age must be between 10 and 120")
		return
	}
	if req.HeightCm < 100 || req.HeightCm > 250 {
		jsonError(w, http.StatusBadRequest, "height must be between 100 and 250")
		return
	}
	if req.WeightKg < 30 || req.WeightKg > 300 {
		jsonError(w, http.StatusBadRequest, "weight must be between 30 and 300")
		return
	}
	if !validGenders[req.Gender] {
		jsonError(w, http.StatusBadRequest, "invalid gender")
		return
	}
	if !validFitness[req.FitnessLevel] {
		jsonError(w, http.StatusBadRequest, "invalid fitness level")
		return
	}
	for _, g := range req.Goals {
		if !validGoals[g] {
			jsonError(w, http.StatusBadRequest, "invalid goal: "+g)
			return
		}
	}

	data := models.RegistrationData{
		Age:          req.Age,
		HeightCm:     req.HeightCm,
		WeightKg:     req.WeightKg,
		Gender:       req.Gender,
		FitnessLevel: req.FitnessLevel,
		Goals:        req.Goals,
	}

	if err := h.userSvc.CreateProfile(r.Context(), user.ID, data); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to create profile")
		return
	}

	if err := h.userSvc.MarkRegistered(r.Context(), user); err != nil {
		// Profile created, registration flag failed — not critical
	}

	jsonResponse(w, http.StatusOK, map[string]bool{"success": true})
}

// GET /app/api/registration/status
func (h *RegistrationHandler) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]bool{"is_registered": user.IsRegistered})
}
