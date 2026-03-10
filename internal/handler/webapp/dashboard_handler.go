package webapp

import (
	"net/http"

	"fitness-bot/internal/service"
)

type DashboardHandler struct {
	dashboardSvc *service.DashboardService
}

func NewDashboardHandler(dashboardSvc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardSvc: dashboardSvc}
}

// GET /app/api/dashboard
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	data, err := h.dashboardSvc.GetDashboard(r.Context(), user)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to load dashboard")
		return
	}

	jsonResponse(w, http.StatusOK, data)
}
