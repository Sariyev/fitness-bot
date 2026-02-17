package webapp

import (
	"net/http"
	"strconv"
	"strings"

	"fitness-bot/internal/service"
)

type ProgressHandler struct {
	moduleSvc *service.ModuleService
}

func NewProgressHandler(moduleSvc *service.ModuleService) *ProgressHandler {
	return &ProgressHandler{moduleSvc: moduleSvc}
}

// POST /app/api/lessons/:id/start
// POST /app/api/lessons/:id/complete
func (h *ProgressHandler) HandleProgressRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse: /app/api/lessons/{id}/{action}
	path := strings.TrimPrefix(r.URL.Path, "/app/api/lessons/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		jsonError(w, http.StatusNotFound, "not found")
		return
	}

	lessonID, err := strconv.Atoi(parts[0])
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid lesson id")
		return
	}

	action := parts[1]

	switch action {
	case "start":
		if err := h.moduleSvc.StartLesson(r.Context(), user.ID, lessonID); err != nil {
			jsonError(w, http.StatusInternalServerError, "failed to start lesson")
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"status": "started"})

	case "complete":
		if err := h.moduleSvc.CompleteLesson(r.Context(), user.ID, lessonID); err != nil {
			jsonError(w, http.StatusInternalServerError, "failed to complete lesson")
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"status": "completed"})

	default:
		jsonError(w, http.StatusNotFound, "not found")
	}
}
