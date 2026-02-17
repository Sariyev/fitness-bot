package admin

import (
	"context"
	"fitness-bot/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type UserAdminHandler struct {
	userRepo  repository.UserRepository
	scoreRepo repository.ScoreRepository
}

func NewUserAdminHandler(userRepo repository.UserRepository, scoreRepo repository.ScoreRepository) *UserAdminHandler {
	return &UserAdminHandler{userRepo: userRepo, scoreRepo: scoreRepo}
}

func (h *UserAdminHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	// For now, return a placeholder — full listing requires a new repo method
	jsonResponse(w, http.StatusOK, map[string]string{"message": "user listing not yet implemented"})
}

func (h *UserAdminHandler) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	parts := strings.SplitN(path, "/", 2)
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Check for sub-resource: /api/users/{id}/scores
	if len(parts) > 1 && parts[1] == "scores" {
		scores, err := h.scoreRepo.ListByUser(ctx, id)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, scores)
		return
	}

	user, err := h.userRepo.GetByTelegramID(ctx, id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "user not found")
		return
	}

	profile, _ := h.userRepo.GetProfileByUserID(ctx, user.ID)

	result := map[string]any{
		"user":    user,
		"profile": profile,
	}
	jsonResponse(w, http.StatusOK, result)
}
