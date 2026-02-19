package webapp

import (
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

// GET /app/api/profile
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

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
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"username":      user.Username,
		"age":           profile.Age,
		"height_cm":     profile.HeightCm,
		"weight_kg":     profile.WeightKg,
		"gender":        profile.Gender,
		"fitness_level": profile.FitnessLevel,
		"goals":         goals,
		"is_paid":       user.IsPaid,
	})
}
