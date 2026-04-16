package webapp

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"fitness-bot/internal/service"
)

type AdminHandler struct {
	userSvc    *service.UserService
	workoutSvc *service.WorkoutService
	scoreSvc   *service.ScoreService
}

func NewAdminHandler(
	userSvc *service.UserService,
	workoutSvc *service.WorkoutService,
	scoreSvc *service.ScoreService,
) *AdminHandler {
	return &AdminHandler{
		userSvc:    userSvc,
		workoutSvc: workoutSvc,
		scoreSvc:   scoreSvc,
	}
}

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

func (h *AdminHandler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	programs, err := h.workoutSvc.ListPrograms(r.Context(), "", "", "")
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list programs")
		return
	}

	jsonResponse(w, http.StatusOK, programs)
}

func (h *AdminHandler) ListWorkouts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	workouts, err := h.workoutSvc.ListWorkouts(r.Context(), "", "", "")
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list workouts")
		return
	}

	jsonResponse(w, http.StatusOK, workouts)
}

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
