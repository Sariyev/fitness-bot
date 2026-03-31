package webapp

import (
	"fmt"
	"log"
	"net/http"

	"fitness-bot/internal/payment"
	"fitness-bot/internal/service"
)

type PaymentHandler struct {
	paymentSvc *service.PaymentService
	verifier   payment.CallbackVerifier // nil for dummy mode
	webAppURL  string
}

func NewPaymentHandler(paymentSvc *service.PaymentService, verifier payment.CallbackVerifier, webAppURL string) *PaymentHandler {
	return &PaymentHandler{
		paymentSvc: paymentSvc,
		verifier:   verifier,
		webAppURL:  webAppURL,
	}
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

	result, err := h.paymentSvc.InitiatePayment(r.Context(), user)
	if err != nil {
		log.Printf("[PAYMENT] initiate failed for user %d: %v", user.ID, err)
		jsonError(w, http.StatusInternalServerError, "payment failed")
		return
	}

	if result.RedirectURL != "" {
		// Async provider (Robokassa) — return URL for frontend to open
		jsonResponse(w, http.StatusOK, map[string]interface{}{
			"success":      true,
			"redirect_url": result.RedirectURL,
		})
	} else {
		// Sync provider (Dummy) — payment already completed
		jsonResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Оплата прошла успешно",
		})
	}
}

// POST /app/api/payment/robokassa/result — server-to-server callback from Robokassa
func (h *PaymentHandler) RobokassaResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.verifier == nil {
		http.Error(w, "robokassa not configured", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	params := map[string]string{
		"OutSum":         r.FormValue("OutSum"),
		"InvId":          r.FormValue("InvId"),
		"SignatureValue": r.FormValue("SignatureValue"),
	}

	paymentID, err := h.verifier.VerifyCallback(params)
	if err != nil {
		log.Printf("[ROBOKASSA] callback verification failed: %v", err)
		http.Error(w, "signature verification failed", http.StatusBadRequest)
		return
	}

	if err := h.paymentSvc.ConfirmPayment(r.Context(), paymentID); err != nil {
		log.Printf("[ROBOKASSA] confirm payment %d failed: %v", paymentID, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	log.Printf("[ROBOKASSA] payment %d confirmed", paymentID)
	fmt.Fprintf(w, "OK%d", paymentID)
}

// GET /app/api/payment/robokassa/success — browser redirect after successful payment
func (h *PaymentHandler) RobokassaSuccess(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, h.webAppURL+"/payment?status=success", http.StatusFound)
}

// GET /app/api/payment/robokassa/fail — browser redirect after failed/cancelled payment
func (h *PaymentHandler) RobokassaFail(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, h.webAppURL+"/payment?status=fail", http.StatusFound)
}
