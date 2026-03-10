package webapp

import (
	"encoding/json"
	"net/http"
	"strings"

	"fitness-bot/internal/service"
)

type ProfileHandler struct {
	userSvc *service.UserService
}

func NewProfileHandler(userSvc *service.UserService) *ProfileHandler {
	return &ProfileHandler{userSvc: userSvc}
}

// HandleProfile dispatches GET and PUT requests for /app/api/profile.
func (h *ProfileHandler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetProfile(w, r)
	case http.MethodPut:
		h.UpdateProfile(w, r)
	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// GET /app/api/profile
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.userSvc.GetProfile(r.Context(), user.ID)
	if err != nil {
		jsonError(w, http.StatusNotFound, "profile not found")
		return
	}

	var goals []string
	if profile.Goal != "" {
		goals = strings.Split(profile.Goal, ",")
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"first_name":          user.FirstName,
		"last_name":           user.LastName,
		"username":            user.Username,
		"age":                 profile.Age,
		"height_cm":           profile.HeightCm,
		"weight_kg":           profile.WeightKg,
		"gender":              profile.Gender,
		"fitness_level":       profile.FitnessLevel,
		"goals":               goals,
		"is_paid":             user.IsPaid,
		"training_access":     profile.TrainingAccess,
		"training_experience": profile.TrainingExperience,
		"has_pain":            profile.HasPain,
		"pain_locations":      profile.PainLocations,
		"pain_level":          profile.PainLevel,
		"diagnoses":           profile.Diagnoses,
		"contraindications":   profile.Contraindications,
		"days_per_week":       profile.DaysPerWeek,
		"session_duration":    profile.SessionDuration,
		"preferred_time":      profile.PreferredTime,
		"equipment":           profile.Equipment,
	})
}

type UpdateProfileRequest struct {
	Age                *int      `json:"age"`
	HeightCm           *int      `json:"height_cm"`
	WeightKg           *float64  `json:"weight_kg"`
	Gender             *string   `json:"gender"`
	FitnessLevel       *string   `json:"fitness_level"`
	Goals              []string  `json:"goals"`
	TrainingAccess     *string   `json:"training_access"`
	TrainingExperience *string   `json:"training_experience"`
	HasPain            *bool     `json:"has_pain"`
	PainLocations      []string  `json:"pain_locations"`
	PainLevel          *int      `json:"pain_level"`
	Diagnoses          []string  `json:"diagnoses"`
	Contraindications  *string   `json:"contraindications"`
	DaysPerWeek        *int      `json:"days_per_week"`
	SessionDuration    *int      `json:"session_duration"`
	PreferredTime      *string   `json:"preferred_time"`
	Equipment          []string  `json:"equipment"`
}

// PUT /app/api/profile
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.userSvc.GetProfile(r.Context(), user.ID)
	if err != nil {
		jsonError(w, http.StatusNotFound, "profile not found")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Apply partial updates
	if req.Age != nil {
		profile.Age = *req.Age
	}
	if req.HeightCm != nil {
		profile.HeightCm = *req.HeightCm
	}
	if req.WeightKg != nil {
		profile.WeightKg = *req.WeightKg
	}
	if req.Gender != nil {
		profile.Gender = *req.Gender
	}
	if req.FitnessLevel != nil {
		profile.FitnessLevel = *req.FitnessLevel
	}
	if req.Goals != nil {
		profile.Goal = strings.Join(req.Goals, ",")
	}
	if req.TrainingAccess != nil {
		profile.TrainingAccess = req.TrainingAccess
	}
	if req.TrainingExperience != nil {
		profile.TrainingExperience = req.TrainingExperience
	}
	if req.HasPain != nil {
		profile.HasPain = *req.HasPain
	}
	if req.PainLocations != nil {
		profile.PainLocations = req.PainLocations
	}
	if req.PainLevel != nil {
		profile.PainLevel = *req.PainLevel
	}
	if req.Diagnoses != nil {
		profile.Diagnoses = req.Diagnoses
	}
	if req.Contraindications != nil {
		profile.Contraindications = *req.Contraindications
	}
	if req.DaysPerWeek != nil {
		profile.DaysPerWeek = req.DaysPerWeek
	}
	if req.SessionDuration != nil {
		profile.SessionDuration = req.SessionDuration
	}
	if req.PreferredTime != nil {
		profile.PreferredTime = req.PreferredTime
	}
	if req.Equipment != nil {
		profile.Equipment = req.Equipment
	}

	if err := h.userSvc.UpdateProfile(r.Context(), profile); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]bool{"success": true})
}
