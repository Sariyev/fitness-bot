package admin

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type ModuleAdminHandler struct {
	repo repository.ModuleRepository
}

func NewModuleAdminHandler(repo repository.ModuleRepository) *ModuleAdminHandler {
	return &ModuleAdminHandler{repo: repo}
}

func (h *ModuleAdminHandler) HandleModules(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	switch r.Method {
	case http.MethodGet:
		modules, err := h.repo.ListActiveModules(ctx)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, modules)

	case http.MethodPost:
		var m models.Module
		if err := decodeJSON(r, &m); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if err := h.repo.CreateModule(ctx, &m); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusCreated, m)

	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *ModuleAdminHandler) HandleModuleByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	path := strings.TrimPrefix(r.URL.Path, "/api/modules/")
	parts := strings.SplitN(path, "/", 2)
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Check for sub-resources: /api/modules/{id}/categories, /api/modules/{id}/categories/{id}/lessons
	if len(parts) > 1 {
		subPath := parts[1]
		switch {
		case subPath == "categories" && r.Method == http.MethodGet:
			cats, err := h.repo.ListCategoriesByModule(ctx, id)
			if err != nil {
				jsonError(w, http.StatusInternalServerError, err.Error())
				return
			}
			jsonResponse(w, http.StatusOK, cats)
			return

		case subPath == "categories" && r.Method == http.MethodPost:
			var c models.ModuleCategory
			if err := decodeJSON(r, &c); err != nil {
				jsonError(w, http.StatusBadRequest, "invalid json")
				return
			}
			c.ModuleID = id
			if err := h.repo.CreateCategory(ctx, &c); err != nil {
				jsonError(w, http.StatusInternalServerError, err.Error())
				return
			}
			jsonResponse(w, http.StatusCreated, c)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		m, err := h.repo.GetModuleByID(ctx, id)
		if err != nil {
			jsonError(w, http.StatusNotFound, "not found")
			return
		}
		jsonResponse(w, http.StatusOK, m)

	case http.MethodPut:
		var m models.Module
		if err := decodeJSON(r, &m); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid json")
			return
		}
		m.ID = id
		if err := h.repo.UpdateModule(ctx, &m); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, m)

	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *ModuleAdminHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	path := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	parts := strings.SplitN(path, "/", 2)
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Sub-resource: /api/categories/{id}/lessons
	if len(parts) > 1 {
		subPath := parts[1]
		switch {
		case subPath == "lessons" && r.Method == http.MethodGet:
			lessons, err := h.repo.ListLessonsByCategory(ctx, id)
			if err != nil {
				jsonError(w, http.StatusInternalServerError, err.Error())
				return
			}
			jsonResponse(w, http.StatusOK, lessons)
			return

		case subPath == "lessons" && r.Method == http.MethodPost:
			var l models.Lesson
			if err := decodeJSON(r, &l); err != nil {
				jsonError(w, http.StatusBadRequest, "invalid json")
				return
			}
			l.CategoryID = id
			if err := h.repo.CreateLesson(ctx, &l); err != nil {
				jsonError(w, http.StatusInternalServerError, err.Error())
				return
			}
			jsonResponse(w, http.StatusCreated, l)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		c, err := h.repo.GetCategoryByID(ctx, id)
		if err != nil {
			jsonError(w, http.StatusNotFound, "not found")
			return
		}
		jsonResponse(w, http.StatusOK, c)

	case http.MethodPut:
		var c models.ModuleCategory
		if err := decodeJSON(r, &c); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid json")
			return
		}
		c.ID = id
		if err := h.repo.UpdateCategory(ctx, &c); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, c)

	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *ModuleAdminHandler) HandleLessonByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	path := strings.TrimPrefix(r.URL.Path, "/api/lessons/")
	parts := strings.SplitN(path, "/", 2)
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Sub-resource: /api/lessons/{id}/contents
	if len(parts) > 1 {
		subPath := parts[1]
		switch {
		case subPath == "contents" && r.Method == http.MethodGet:
			contents, err := h.repo.ListContentByLesson(ctx, id)
			if err != nil {
				jsonError(w, http.StatusInternalServerError, err.Error())
				return
			}
			jsonResponse(w, http.StatusOK, contents)
			return

		case subPath == "contents" && r.Method == http.MethodPost:
			var c models.LessonContent
			if err := decodeJSON(r, &c); err != nil {
				jsonError(w, http.StatusBadRequest, "invalid json")
				return
			}
			c.LessonID = id
			if err := h.repo.CreateContent(ctx, &c); err != nil {
				jsonError(w, http.StatusInternalServerError, err.Error())
				return
			}
			jsonResponse(w, http.StatusCreated, c)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		l, err := h.repo.GetLessonByID(ctx, id)
		if err != nil {
			jsonError(w, http.StatusNotFound, "not found")
			return
		}
		jsonResponse(w, http.StatusOK, l)

	case http.MethodPut:
		var l models.Lesson
		if err := decodeJSON(r, &l); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid json")
			return
		}
		l.ID = id
		if err := h.repo.UpdateLesson(ctx, &l); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, l)

	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *ModuleAdminHandler) HandleContentByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	idStr := strings.TrimPrefix(r.URL.Path, "/api/contents/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		c, err := h.repo.GetContentByID(ctx, id)
		if err != nil {
			jsonError(w, http.StatusNotFound, "not found")
			return
		}
		jsonResponse(w, http.StatusOK, c)

	case http.MethodPut:
		var c models.LessonContent
		if err := decodeJSON(r, &c); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid json")
			return
		}
		c.ID = id
		if err := h.repo.UpdateContent(ctx, &c); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, c)

	case http.MethodDelete:
		if err := h.repo.DeleteContent(ctx, id); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"status": "deleted"})

	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
