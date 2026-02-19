package webapp

import (
	"net/http"

	"fitness-bot/internal/service"
)

type PaymentHandler struct {
	paymentSvc *service.PaymentService
}

func NewPaymentHandler(paymentSvc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentSvc: paymentSvc}
}

// GET /app/api/payment/status
func (h *PaymentHandler) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"is_paid":   user.IsPaid,
		"price_kzt": 5000,
	})
}

// POST /app/api/payment/pay
func (h *PaymentHandler) Pay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if user.IsPaid {
		jsonResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Доступ уже оплачен",
		})
		return
	}

	if err := h.paymentSvc.ProcessPayment(r.Context(), user); err != nil {
		jsonError(w, http.StatusInternalServerError, "payment failed")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Оплата прошла успешно",
	})
}
