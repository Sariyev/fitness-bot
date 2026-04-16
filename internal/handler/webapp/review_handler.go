package webapp

import (
	"encoding/json"
	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type ReviewHandler struct {
	scoreSvc *service.ScoreService
}

func NewReviewHandler(scoreSvc *service.ScoreService) *ReviewHandler {
	return &ReviewHandler{scoreSvc: scoreSvc}
}

type createReviewRequest struct {
	ReferenceType string   `json:"reference_type"`
	ReferenceID   int      `json:"reference_id"`
	Score         int      `json:"score"`
	Comment       string   `json:"comment"`
	Tags          []string `json:"tags"`
}

// HandleReviewRoutes dispatches /app/api/reviews and /app/api/reviews/...
func (h *ReviewHandler) HandleReviewRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/app/api/reviews")
	path = strings.TrimPrefix(path, "/")

	switch {
	case path == "" && r.Method == http.MethodPost:
		h.CreateReview(w, r)
	case path == "tags":
		h.GetTags(w, r)
	case path == "bot/summary":
		h.GetBotSummary(w, r)
	case path == "my":
		h.GetMyReviews(w, r)
	case strings.HasPrefix(path, "summary"):
		h.GetSummary(w, r)
	case r.Method == http.MethodGet:
		h.GetReviews(w, r)
	default:
		jsonError(w, http.StatusNotFound, "not found")
	}
}

// POST /app/api/reviews — create a review
func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Score < 1 || req.Score > 5 {
		jsonError(w, http.StatusBadRequest, "score must be between 1 and 5")
		return
	}
	if req.ReferenceType == "" {
		jsonError(w, http.StatusBadRequest, "reference_type is required")
		return
	}

	if req.Tags == nil {
		req.Tags = []string{}
	}

	score := &models.UserScore{
		UserID:        user.ID,
		ScoreType:     "rating",
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
		Score:         req.Score,
		Comment:       req.Comment,
		Tags:          req.Tags,
	}

	if err := h.scoreSvc.SaveScore(r.Context(), score); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to save review")
		return
	}

	jsonResponse(w, http.StatusCreated, score)
}

// GET /app/api/reviews?reference_type=X&reference_id=Y — list reviews for content
func (h *ReviewHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	refType := r.URL.Query().Get("reference_type")
	refIDStr := r.URL.Query().Get("reference_id")

	if refType == "" || refIDStr == "" {
		jsonError(w, http.StatusBadRequest, "reference_type and reference_id are required")
		return
	}

	refID, err := strconv.Atoi(refIDStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid reference_id")
		return
	}

	reviews, err := h.scoreSvc.ListByReference(r.Context(), refType, refID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to get reviews")
		return
	}

	if reviews == nil {
		reviews = []models.UserScore{}
	}

	jsonResponse(w, http.StatusOK, reviews)
}

// GET /app/api/reviews/my — list current user's reviews
func (h *ReviewHandler) GetMyReviews(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	reviews, err := h.scoreSvc.ListByUser(r.Context(), user.ID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to get reviews")
		return
	}

	if reviews == nil {
		reviews = []models.UserScore{}
	}

	jsonResponse(w, http.StatusOK, reviews)
}

// GET /app/api/reviews/summary?reference_type=X&reference_id=Y — get average score + count
func (h *ReviewHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	refType := r.URL.Query().Get("reference_type")
	refIDStr := r.URL.Query().Get("reference_id")

	if refType == "" || refIDStr == "" {
		jsonError(w, http.StatusBadRequest, "reference_type and reference_id are required")
		return
	}

	refID, err := strconv.Atoi(refIDStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid reference_id")
		return
	}

	summary, err := h.scoreSvc.GetSummary(r.Context(), refType, refID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to get summary")
		return
	}

	jsonResponse(w, http.StatusOK, summary)
}

// GET /app/api/reviews/bot/summary — get bot review summary
func (h *ReviewHandler) GetBotSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.scoreSvc.GetBotSummary(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to get bot summary")
		return
	}

	jsonResponse(w, http.StatusOK, summary)
}

// GET /app/api/reviews/tags?reference_type=X — get preset tags for a reference type
func (h *ReviewHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	refType := r.URL.Query().Get("reference_type")
	if refType == "" {
		refType = "bot"
	}

	tags, ok := models.ReviewTags[refType]
	if !ok {
		tags = models.ReviewTags["bot"]
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"reference_type": refType,
		"tags":           tags,
	})
}
