package webapp

import (
	"encoding/json"
	"net/http"
	"time"

	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
)

type ProgressV2Handler struct {
	progressSvc *service.ProgressService
}

func NewProgressV2Handler(progressSvc *service.ProgressService) *ProgressV2Handler {
	return &ProgressV2Handler{progressSvc: progressSvc}
}

type CreateProgressRequest struct {
	Date         string         `json:"date"`
	WeightKg     *float64       `json:"weight_kg"`
	Measurements map[string]any `json:"measurements"`
	PhotoURL     string         `json:"photo_url"`
	Wellbeing    string         `json:"wellbeing"`
	PainLevel    int            `json:"pain_level"`
}

type StatsResponse struct {
	CurrentStreak int      `json:"current_streak"`
	LongestStreak int      `json:"longest_streak"`
	Calendar      []string `json:"calendar"`
}

// HandleProgressRoutes dispatches /app/api/progress requests.
//
//	GET  /app/api/progress -> list user progress entries
//	POST /app/api/progress -> add progress entry
func (h *ProgressV2Handler) HandleProgressRoutes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ListProgress(w, r)
	case http.MethodPost:
		h.AddProgress(w, r)
	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// HandleProgressStats handles GET /app/api/progress/stats
func (h *ProgressV2Handler) HandleProgressStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	current, longest, err := h.progressSvc.GetStreak(r.Context(), user.ID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to load streaks")
		return
	}

	now := time.Now()
	calendar, err := h.progressSvc.GetCalendar(r.Context(), user.ID, now.Year(), int(now.Month()))
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to load calendar")
		return
	}

	jsonResponse(w, http.StatusOK, StatsResponse{
		CurrentStreak: current,
		LongestStreak: longest,
		Calendar:      calendar,
	})
}

// HandleProgressAchievements handles GET /app/api/progress/achievements
func (h *ProgressV2Handler) HandleProgressAchievements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	achievements, err := h.progressSvc.GetAchievements(r.Context(), user.ID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to load achievements")
		return
	}

	jsonResponse(w, http.StatusOK, achievements)
}

// ListProgress handles GET /app/api/progress
func (h *ProgressV2Handler) ListProgress(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	entries, err := h.progressSvc.ListEntries(r.Context(), user.ID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list progress")
		return
	}

	jsonResponse(w, http.StatusOK, entries)
}

// AddProgress handles POST /app/api/progress
func (h *ProgressV2Handler) AddProgress(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Date == "" {
		jsonError(w, http.StatusBadRequest, "date is required")
		return
	}

	entry := &models.ProgressEntry{
		UserID:       user.ID,
		Date:         req.Date,
		WeightKg:     req.WeightKg,
		Measurements: req.Measurements,
		PhotoURL:     req.PhotoURL,
		Wellbeing:    req.Wellbeing,
		PainLevel:    req.PainLevel,
	}

	if err := h.progressSvc.AddEntry(r.Context(), entry); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to add progress entry")
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]bool{"success": true})
}
