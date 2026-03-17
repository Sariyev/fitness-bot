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
}

func NewWorkoutHandler(workoutSvc *service.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{workoutSvc: workoutSvc}
}

// HandleProgramRoutes dispatches /app/api/programs/... requests.
//
//	GET  /app/api/programs              -> list programs with query filters
//	GET  /app/api/programs/{id}         -> get program by ID
//	POST /app/api/programs/{id}/enroll  -> enroll user in program
func (h *WorkoutHandler) HandleProgramRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/app/api/programs")
	path = strings.TrimPrefix(path, "/")

	// GET /app/api/programs
	if path == "" {
		h.ListPrograms(w, r)
		return
	}

	parts := strings.Split(path, "/")

	// GET /app/api/programs/{id}
	if len(parts) == 1 {
		if r.Method != http.MethodGet {
			jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.GetProgram(w, r, parts[0])
		return
	}

	// POST /app/api/programs/{id}/enroll
	if len(parts) == 2 && parts[1] == "enroll" {
		if r.Method != http.MethodPost {
			jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.EnrollProgram(w, r, parts[0])
		return
	}

	jsonError(w, http.StatusNotFound, "not found")
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

// ListPrograms handles GET /app/api/programs?format=&goal=&level=
func (h *WorkoutHandler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	format := r.URL.Query().Get("format")
	goal := r.URL.Query().Get("goal")
	level := r.URL.Query().Get("level")

	programs, err := h.workoutSvc.ListPrograms(r.Context(), format, goal, level)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list programs")
		return
	}

	jsonResponse(w, http.StatusOK, programs)
}

// GetProgram handles GET /app/api/programs/{id}
func (h *WorkoutHandler) GetProgram(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid program id")
		return
	}

	program, err := h.workoutSvc.GetProgram(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "program not found")
		return
	}

	jsonResponse(w, http.StatusOK, program)
}

// EnrollProgram handles POST /app/api/programs/{id}/enroll
func (h *WorkoutHandler) EnrollProgram(w http.ResponseWriter, r *http.Request, idStr string) {
	user := UserFromContext(r.Context())
	if user == nil {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "invalid program id")
		return
	}

	if err := h.workoutSvc.EnrollInProgram(r.Context(), user.ID, id); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to enroll in program")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]bool{"success": true})
}

// ListWorkouts handles GET /app/api/workouts?format=&goal=&level=
func (h *WorkoutHandler) ListWorkouts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
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

	jsonResponse(w, http.StatusOK, workouts)
}

// GetWorkout handles GET /app/api/workouts/{id}
func (h *WorkoutHandler) GetWorkout(w http.ResponseWriter, r *http.Request, idStr string) {
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
