package webapp

import (
	"net/http"
	"strconv"
	"strings"

	"fitness-bot/internal/service"
)

type ModuleHandler struct {
	moduleSvc *service.ModuleService
}

func NewModuleHandler(moduleSvc *service.ModuleService) *ModuleHandler {
	return &ModuleHandler{moduleSvc: moduleSvc}
}

// GET /app/api/modules
func (h *ModuleHandler) ListModules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	modules, err := h.moduleSvc.ListModules(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list modules")
		return
	}

	jsonResponse(w, http.StatusOK, modules)
}

// GET /app/api/modules/:id/categories
func (h *ModuleHandler) HandleModuleRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse: /app/api/modules/{id}/categories
	path := strings.TrimPrefix(r.URL.Path, "/app/api/modules/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[1] != "categories" {
		jsonError(w, http.StatusNotFound, "not found")
		return
	}

	moduleID, err := strconv.Atoi(parts[0])
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid module id")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	categories, err := h.moduleSvc.ListCategories(r.Context(), moduleID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list categories")
		return
	}

	// Enrich with progress data
	type CategoryWithProgress struct {
		ID          int    `json:"id"`
		ModuleID    int    `json:"module_id"`
		Slug        string `json:"slug"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		SortOrder   int    `json:"sort_order"`
		Completed   int    `json:"completed"`
		Total       int    `json:"total"`
	}

	var result []CategoryWithProgress
	for _, cat := range categories {
		completed, total, _ := h.moduleSvc.GetCategoryProgress(r.Context(), user.ID, cat.ID)
		result = append(result, CategoryWithProgress{
			ID:          cat.ID,
			ModuleID:    cat.ModuleID,
			Slug:        cat.Slug,
			Name:        cat.Name,
			Description: cat.Description,
			Icon:        cat.Icon,
			SortOrder:   cat.SortOrder,
			Completed:   completed,
			Total:       total,
		})
	}

	jsonResponse(w, http.StatusOK, result)
}

// GET /app/api/categories/:id/lessons
func (h *ModuleHandler) HandleCategoryRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse: /app/api/categories/{id}/lessons
	path := strings.TrimPrefix(r.URL.Path, "/app/api/categories/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[1] != "lessons" {
		jsonError(w, http.StatusNotFound, "not found")
		return
	}

	categoryID, err := strconv.Atoi(parts[0])
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	lessons, err := h.moduleSvc.ListLessons(r.Context(), categoryID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list lessons")
		return
	}

	// Enrich with progress
	type LessonWithProgress struct {
		ID          int    `json:"id"`
		CategoryID  int    `json:"category_id"`
		Slug        string `json:"slug"`
		Title       string `json:"title"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
		Status      string `json:"status"` // "not_started", "in_progress", "completed"
	}

	var result []LessonWithProgress
	for _, les := range lessons {
		status := "not_started"
		progress, err := h.moduleSvc.GetLessonProgress(r.Context(), user.ID, les.ID)
		if err == nil && progress != nil {
			status = progress.Status
		}
		result = append(result, LessonWithProgress{
			ID:          les.ID,
			CategoryID:  les.CategoryID,
			Slug:        les.Slug,
			Title:       les.Title,
			Description: les.Description,
			SortOrder:   les.SortOrder,
			Status:      status,
		})
	}

	jsonResponse(w, http.StatusOK, result)
}

// GET /app/api/lessons/:id
func (h *ModuleHandler) HandleLessonRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse: /app/api/lessons/{id}
	path := strings.TrimPrefix(r.URL.Path, "/app/api/lessons/")
	lessonID, err := strconv.Atoi(strings.TrimSuffix(path, "/"))
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid lesson id")
		return
	}

	lesson, err := h.moduleSvc.GetLesson(r.Context(), lessonID)
	if err != nil {
		jsonError(w, http.StatusNotFound, "lesson not found")
		return
	}

	content, err := h.moduleSvc.GetLessonContent(r.Context(), lessonID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to load content")
		return
	}

	user := UserFromContext(r.Context())
	status := "not_started"
	if user != nil {
		progress, err := h.moduleSvc.GetLessonProgress(r.Context(), user.ID, lessonID)
		if err == nil && progress != nil {
			status = progress.Status
		}
	}

	type LessonDetail struct {
		ID          int         `json:"id"`
		CategoryID  int         `json:"category_id"`
		Title       string      `json:"title"`
		Description string      `json:"description"`
		Status      string      `json:"status"`
		Content     interface{} `json:"content"`
	}

	jsonResponse(w, http.StatusOK, LessonDetail{
		ID:          lesson.ID,
		CategoryID:  lesson.CategoryID,
		Title:       lesson.Title,
		Description: lesson.Description,
		Status:      status,
		Content:     content,
	})
}

// GET /app/api/subscription/status
func (h *ModuleHandler) PaymentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]bool{"active": user.IsPaid})
}
