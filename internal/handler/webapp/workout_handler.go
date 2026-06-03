package webapp

import (
	"net/http"
	"strconv"
	"strings"

	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
)

type WorkoutHandler struct {
	workoutSvc *service.WorkoutService
	accessSvc  *service.AccessService
	mediaSvc   *service.MediaService // optional; nil when R2 not configured
}

func NewWorkoutHandler(workoutSvc *service.WorkoutService, accessSvc *service.AccessService, mediaSvc *service.MediaService) *WorkoutHandler {
	return &WorkoutHandler{workoutSvc: workoutSvc, accessSvc: accessSvc, mediaSvc: mediaSvc}
}

// HandleWorkoutRoutes dispatches /app/api/workouts/... requests.
//
//	GET  /app/api/workouts               -> list workouts with query filters
//	GET  /app/api/workouts/{id}          -> get workout by ID with exercises
//	POST /app/api/workouts/{id}/complete -> mark workout completed
func (h *WorkoutHandler) HandleWorkoutRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/app/api/workouts")
	path = strings.TrimPrefix(path, "/")

	// GET /app/api/workouts
	if path == "" {
		h.ListWorkouts(w, r)
		return
	}

	parts := strings.Split(path, "/")

	// GET /app/api/workouts/{id}
	if len(parts) == 1 {
		if r.Method != http.MethodGet {
			jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.GetWorkout(w, r, parts[0])
		return
	}

	// POST /app/api/workouts/{id}/complete
	if len(parts) == 2 && parts[1] == "complete" {
		if r.Method != http.MethodPost {
			jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.CompleteWorkout(w, r, parts[0])
		return
	}

	jsonError(w, http.StatusNotFound, "not found")
}

// ListWorkouts handles GET /app/api/workouts?format=&goal=&level=
// Each item gets a `locked` flag so the frontend can render lock badges.
func (h *WorkoutHandler) ListWorkouts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	format := r.URL.Query().Get("format")
	goal := r.URL.Query().Get("goal")
	level := r.URL.Query().Get("level")

	workouts, err := h.workoutSvc.ListWorkouts(r.Context(), format, goal, level)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list workouts")
		return
	}

	for i := range workouts {
		can, _ := h.accessSvc.CanAccess(r.Context(), user, workouts[i].AccessTier, models.CategoryWorkouts)
		workouts[i].Locked = !can
	}

	jsonResponse(w, http.StatusOK, workouts)
}

// GetWorkout handles GET /app/api/workouts/{id}. Returns 402 with a payment
// CTA payload when the workout is locked.
func (h *WorkoutHandler) GetWorkout(w http.ResponseWriter, r *http.Request, idStr string) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid workout id")
		return
	}

	workout, err := h.workoutSvc.GetWorkout(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "workout not found")
		return
	}

	can, err := h.accessSvc.CanAccess(r.Context(), user, workout.AccessTier, models.CategoryWorkouts)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "access check failed")
		return
	}
	if !can {
		price, _ := h.accessSvc.GetPrice(r.Context(), models.CategoryWorkouts)
		jsonResponse(w, http.StatusPaymentRequired, map[string]interface{}{
			"error":     "locked",
			"category":  models.CategoryWorkouts,
			"tier":      workout.AccessTier,
			"price_kzt": price,
		})
		return
	}
	workout.Locked = false

	// If admin uploaded a video, replace video_url with a presigned R2 URL.
	// Access has already been gated above.
	if workout.VideoMediaID != nil && h.mediaSvc != nil {
		if u, urlErr := h.mediaSvc.PresignReadURL(r.Context(), *workout.VideoMediaID); urlErr == nil {
			workout.VideoURL = u
		}
	}

	exercises, err := h.workoutSvc.GetWorkoutExercisesWithDetails(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to load exercises")
		return
	}

	type workoutWithExercises struct {
		*models.Workout
		Exercises []models.WorkoutExerciseWithDetails `json:"exercises"`
	}

	jsonResponse(w, http.StatusOK, workoutWithExercises{
		Workout:   workout,
		Exercises: exercises,
	})
}

// CompleteWorkout handles POST /app/api/workouts/{id}/complete
func (h *WorkoutHandler) CompleteWorkout(w http.ResponseWriter, r *http.Request, idStr string) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid workout id")
		return
	}

	if err := h.workoutSvc.CompleteWorkout(r.Context(), user.ID, id); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to complete workout")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]bool{"success": true})
}
