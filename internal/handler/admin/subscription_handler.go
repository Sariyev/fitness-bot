package admin

import (
	"context"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type SubscriptionAdminHandler struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionAdminHandler(repo repository.SubscriptionRepository) *SubscriptionAdminHandler {
	return &SubscriptionAdminHandler{repo: repo}
}

func (h *SubscriptionAdminHandler) HandlePlans(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	switch r.Method {
	case http.MethodGet:
		plans, err := h.repo.ListActivePlans(ctx)
		if err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, plans)

	case http.MethodPost:
		var p models.SubscriptionPlan
		if err := decodeJSON(r, &p); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid json")
			return
		}
		if err := h.repo.CreatePlan(ctx, &p); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusCreated, p)

	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *SubscriptionAdminHandler) HandlePlanByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	idStr := strings.TrimPrefix(r.URL.Path, "/api/subscription-plans/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		plan, err := h.repo.GetPlanByID(ctx, id)
		if err != nil {
			jsonError(w, http.StatusNotFound, "not found")
			return
		}
		jsonResponse(w, http.StatusOK, plan)

	case http.MethodPut:
		var p models.SubscriptionPlan
		if err := decodeJSON(r, &p); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid json")
			return
		}
		p.ID = id
		if err := h.repo.UpdatePlan(ctx, &p); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, p)

	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
