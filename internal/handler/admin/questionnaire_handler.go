package admin

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type QuestionnaireAdminHandler struct {
	repo repository.QuestionnaireRepository
}

func NewQuestionnaireAdminHandler(repo repository.QuestionnaireRepository) *QuestionnaireAdminHandler {
	return &QuestionnaireAdminHandler{repo: repo}
}

func (h *QuestionnaireAdminHandler) HandleQuestionnaires(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	switch r.Method {
	case http.MethodGet:
		quizzes, err := h.repo.ListActive(ctx)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, quizzes)

	case http.MethodPost:
		var q models.Questionnaire
		if err := decodeJSON(r, &q); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if err := h.repo.Create(ctx, &q); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusCreated, q)

	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *QuestionnaireAdminHandler) HandleQuestionnaireByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	path := strings.TrimPrefix(r.URL.Path, "/api/questionnaires/")
	parts := strings.SplitN(path, "/", 2)
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if len(parts) > 1 && parts[1] == "questions" {
		switch r.Method {
		case http.MethodGet:
			questions, err := h.repo.GetQuestionsByQuestionnaireID(ctx, id)
			if err != nil {
				jsonError(w, http.StatusInternalServerError, err.Error())
				return
			}
			jsonResponse(w, http.StatusOK, questions)
			return

		case http.MethodPost:
			var q models.Question
			if err := decodeJSON(r, &q); err != nil {
				jsonError(w, http.StatusBadRequest, "invalid json")
				return
			}
			q.QuestionnaireID = id
			if err := h.repo.CreateQuestion(ctx, &q); err != nil {
				jsonError(w, http.StatusInternalServerError, err.Error())
				return
			}
			jsonResponse(w, http.StatusCreated, q)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		q, err := h.repo.GetByID(ctx, id)
		if err != nil {
			jsonError(w, http.StatusNotFound, "not found")
			return
		}
		jsonResponse(w, http.StatusOK, q)

	case http.MethodPut:
		var q models.Questionnaire
		if err := decodeJSON(r, &q); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid json")
			return
		}
		q.ID = id
		if err := h.repo.Update(ctx, &q); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, q)

	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *QuestionnaireAdminHandler) HandleQuestionByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	path := strings.TrimPrefix(r.URL.Path, "/api/questions/")
	parts := strings.SplitN(path, "/", 2)
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Sub-resource: /api/questions/{id}/options
	if len(parts) > 1 && parts[1] == "options" && r.Method == http.MethodPost {
		var o models.QuestionOption
		if err := decodeJSON(r, &o); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid json")
			return
		}
		o.QuestionID = id
		if err := h.repo.CreateOption(ctx, &o); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusCreated, o)
		return
	}

	switch r.Method {
	case http.MethodGet:
		q, err := h.repo.GetQuestionByID(ctx, id)
		if err != nil {
			jsonError(w, http.StatusNotFound, "not found")
			return
		}
		jsonResponse(w, http.StatusOK, q)

	case http.MethodPut:
		var q models.Question
		if err := decodeJSON(r, &q); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid json")
			return
		}
		q.ID = id
		if err := h.repo.UpdateQuestion(ctx, &q); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, q)

	case http.MethodDelete:
		if err := h.repo.DeleteQuestion(ctx, id); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"status": "deleted"})

	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *QuestionnaireAdminHandler) HandleOptionByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	idStr := strings.TrimPrefix(r.URL.Path, "/api/options/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if r.Method == http.MethodDelete {
		if err := h.repo.DeleteOption(ctx, id); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
		return
	}

	jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
}
