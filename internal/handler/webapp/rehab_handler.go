package webapp

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
)

type RehabHandler struct {
	rehabSvc  *service.RehabService
	accessSvc *service.AccessService
	mediaSvc  *service.MediaService // optional; nil when R2 not configured
}

func NewRehabHandler(rehabSvc *service.RehabService, accessSvc *service.AccessService, mediaSvc *service.MediaService) *RehabHandler {
	return &RehabHandler{rehabSvc: rehabSvc, accessSvc: accessSvc, mediaSvc: mediaSvc}
}

type CompleteRehabRequest struct {
	PainLevel int    `json:"pain_level"`
	Comment   string `json:"comment"`
	DayNumber int    `json:"day_number"`
	CourseID  int    `json:"course_id"`
}

// HandleRehabRoutes dispatches /app/api/rehab/courses/... requests.
//
//	GET /app/api/rehab/courses          -> list courses, optional ?category= filter
//	GET /app/api/rehab/courses/{id}     -> get course with sessions
func (h *RehabHandler) HandleRehabRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/app/api/rehab/courses")
	path = strings.TrimPrefix(path, "/")

	// GET /app/api/rehab/courses
	if path == "" {
		h.ListCourses(w, r)
		return
	}

	// GET /app/api/rehab/courses/{id}
	h.GetCourse(w, r, path)
}

// HandleRehabSessionRoutes dispatches /app/api/rehab/sessions/... requests.
//
//	GET  /app/api/rehab/sessions/{id}          -> get session detail
//	POST /app/api/rehab/sessions/{id}/complete -> complete session with pain level
func (h *RehabHandler) HandleRehabSessionRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/app/api/rehab/sessions/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		jsonError(w, http.StatusNotFound, "not found")
		return
	}

	// GET /app/api/rehab/sessions/{id}
	if len(parts) == 1 {
		if r.Method != http.MethodGet {
			jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.GetSession(w, r, parts[0])
		return
	}

	// POST /app/api/rehab/sessions/{id}/complete
	if len(parts) == 2 && parts[1] == "complete" {
		if r.Method != http.MethodPost {
			jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.CompleteSession(w, r, parts[0])
		return
	}

	jsonError(w, http.StatusNotFound, "not found")
}

// HandleRehabProgressRoutes dispatches /app/api/rehab/progress/... requests.
//
//	GET /app/api/rehab/progress/{courseId} -> get user progress for course
func (h *RehabHandler) HandleRehabProgressRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/app/api/rehab/progress/")
	courseID, err := strconv.Atoi(strings.TrimSuffix(path, "/"))
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid course id")
		return
	}

	progress, err := h.rehabSvc.GetUserProgress(r.Context(), user.ID, courseID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to load progress")
		return
	}

	jsonResponse(w, http.StatusOK, progress)
}

// ListCourses handles GET /app/api/rehab/courses?category=
// Each item is annotated with `locked` based on the viewing user's access.
func (h *RehabHandler) ListCourses(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	category := r.URL.Query().Get("category")

	courses, err := h.rehabSvc.ListCourses(r.Context(), category)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list courses")
		return
	}

	for i := range courses {
		can, _ := h.accessSvc.CanAccess(r.Context(), user, courses[i].AccessTier, models.CategoryLFK)
		courses[i].Locked = !can
	}

	jsonResponse(w, http.StatusOK, courses)
}

// GetCourse handles GET /app/api/rehab/courses/{id}
// Returns 402 if the course is locked.
func (h *RehabHandler) GetCourse(w http.ResponseWriter, r *http.Request, idStr string) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.Atoi(strings.TrimSuffix(idStr, "/"))
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid course id")
		return
	}

	course, err := h.rehabSvc.GetCourse(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "course not found")
		return
	}

	can, err := h.accessSvc.CanAccess(r.Context(), user, course.AccessTier, models.CategoryLFK)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "access check failed")
		return
	}
	if !can {
		price, _ := h.accessSvc.GetPrice(r.Context(), models.CategoryLFK)
		jsonResponse(w, http.StatusPaymentRequired, map[string]interface{}{
			"error":     "locked",
			"category":  models.CategoryLFK,
			"tier":      course.AccessTier,
			"price_kzt": price,
		})
		return
	}
	course.Locked = false

	sessions, err := h.rehabSvc.GetCourseSessions(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to load sessions")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"course":   course,
		"sessions": sessions,
	})
}

// GetSession handles GET /app/api/rehab/sessions/{id}
func (h *RehabHandler) GetSession(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(strings.TrimSuffix(idStr, "/"))
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	session, err := h.rehabSvc.GetSession(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "session not found")
		return
	}

	// If admin uploaded a video, replace video_url with a presigned R2 URL.
	// Uses PresignReadURL (no owner check) because access is gated at the
	// course level via accessSvc.CanAccess above.
	if session.VideoMediaID != nil && h.mediaSvc != nil {
		if u, urlErr := h.mediaSvc.PresignReadURL(r.Context(), *session.VideoMediaID); urlErr == nil {
			session.VideoURL = u
		}
	}

	jsonResponse(w, http.StatusOK, session)
}

// CompleteSession handles POST /app/api/rehab/sessions/{id}/complete
func (h *RehabHandler) CompleteSession(w http.ResponseWriter, r *http.Request, idStr string) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	sessionID, err := strconv.Atoi(idStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	var req CompleteRehabRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.rehabSvc.CompleteSession(r.Context(), user.ID, req.CourseID, sessionID, req.DayNumber, req.PainLevel, req.Comment); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to complete session")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]bool{"success": true})
}
