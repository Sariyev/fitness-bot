package webapp

import (
	"net/http"

	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
)

type AccessHandler struct {
	accessSvc *service.AccessService
}

func NewAccessHandler(accessSvc *service.AccessService) *AccessHandler {
	return &AccessHandler{accessSvc: accessSvc}
}

// GetStatus answers GET /app/api/access/status — used by the frontend to
// drive lock badges and "Unlock" CTAs without checking each item individually.
//
// Response shape:
//
//	{
//	  "granted":           ["workouts"],          // categories the user has paid for
//	  "trial_remaining":   86400,                 // seconds left in the signup trial window (0 if expired)
//	  "trial_active":      true,                  // convenience flag (= trial_remaining > 0)
//	  "legacy_paid":       false,                 // grandfathered via users.is_paid
//	  "prices": {                                 // current category prices in KZT
//	    "workouts":  5000,
//	    "lfk":       5000,
//	    "nutrition": 5000
//	  }
//	}
func (h *AccessHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	granted, err := h.accessSvc.ListGranted(r.Context(), user.ID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to read access")
		return
	}
	if granted == nil {
		granted = []models.Category{}
	}

	prices, err := h.accessSvc.ListPrices(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to read prices")
		return
	}

	trialLeft := h.accessSvc.TrialRemaining(user)
	trialSec := int(trialLeft.Seconds())
	if trialSec < 0 {
		trialSec = 0
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"granted":         granted,
		"trial_remaining": trialSec,
		"trial_active":    trialSec > 0,
		"legacy_paid":     user.IsPaid,
		"prices":          prices,
	})
}
