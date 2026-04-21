package webapp

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fitness-bot/internal/payment"
	"fitness-bot/internal/service"
)

type WebAppRouter struct {
	mux          *http.ServeMux
	moduleSvc    *service.ModuleService
	userSvc      *service.UserService
	paymentSvc   *service.PaymentService
	workoutSvc   *service.WorkoutService
	rehabSvc     *service.RehabService
	nutritionSvc *service.NutritionService
	progressSvc  *service.ProgressService
	dashboardSvc *service.DashboardService
	recommendSvc *service.RecommendationService
	scoreSvc     *service.ScoreService
	botToken     string
	staticDir    string
	verifier     payment.CallbackVerifier
	webAppURL    string
}

func NewRouter(
	botToken string,
	userSvc *service.UserService,
	moduleSvc *service.ModuleService,
	paymentSvc *service.PaymentService,
	workoutSvc *service.WorkoutService,
	rehabSvc *service.RehabService,
	nutritionSvc *service.NutritionService,
	progressSvc *service.ProgressService,
	dashboardSvc *service.DashboardService,
	recommendSvc *service.RecommendationService,
	scoreSvc *service.ScoreService,
	staticDir string,
	verifier payment.CallbackVerifier,
	webAppURL string,
) *WebAppRouter {
	r := &WebAppRouter{
		mux:          http.NewServeMux(),
		moduleSvc:    moduleSvc,
		userSvc:      userSvc,
		paymentSvc:   paymentSvc,
		workoutSvc:   workoutSvc,
		rehabSvc:     rehabSvc,
		nutritionSvc: nutritionSvc,
		progressSvc:  progressSvc,
		dashboardSvc: dashboardSvc,
		recommendSvc: recommendSvc,
		scoreSvc:     scoreSvc,
		botToken:     botToken,
		staticDir:    staticDir,
		verifier:     verifier,
		webAppURL:    webAppURL,
	}

	r.setupRoutes()
	return r
}

func (r *WebAppRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Set CORS headers for Telegram Mini App
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Telegram-Init-Data, Authorization")

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	r.mux.ServeHTTP(w, req)
}

func (r *WebAppRouter) setupRoutes() {
	auth := AuthMiddleware(r.botToken, r.userSvc)

	moduleHandler := NewModuleHandler(r.moduleSvc)
	progressHandler := NewProgressHandler(r.moduleSvc)
	paymentHandler := NewPaymentHandler(r.paymentSvc, r.verifier, r.webAppURL)
	profileHandler := NewProfileHandler(r.userSvc)
	registrationHandler := NewRegistrationHandler(r.userSvc, r.recommendSvc)
	dashboardHandler := NewDashboardHandler(r.dashboardSvc)
	workoutHandler := NewWorkoutHandler(r.workoutSvc)
	rehabHandler := NewRehabHandler(r.rehabSvc)
	nutritionHandler := NewNutritionHandler(r.nutritionSvc)
	progressV2Handler := NewProgressV2Handler(r.progressSvc)
	reviewHandler := NewReviewHandler(r.scoreSvc)

	// Auth — exchanges initData for a session token
	r.mux.HandleFunc("/app/api/auth", AuthHandler(r.botToken, r.userSvc))

	// API routes (authenticated)
	r.mux.Handle("/app/api/modules", auth(http.HandlerFunc(moduleHandler.ListModules)))
	r.mux.Handle("/app/api/modules/", auth(http.HandlerFunc(moduleHandler.HandleModuleRoutes)))
	r.mux.Handle("/app/api/categories/", auth(http.HandlerFunc(moduleHandler.HandleCategoryRoutes)))
	r.mux.Handle("/app/api/lessons/", auth(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		path := strings.TrimPrefix(req.URL.Path, "/app/api/lessons/")
		if strings.HasSuffix(path, "/start") || strings.HasSuffix(path, "/complete") {
			progressHandler.HandleProgressRoutes(w, req)
		} else {
			moduleHandler.HandleLessonRoutes(w, req)
		}
	})))
	r.mux.Handle("/app/api/subscription/status", auth(http.HandlerFunc(moduleHandler.PaymentStatus)))
	r.mux.Handle("/app/api/payment/status", auth(http.HandlerFunc(paymentHandler.Status)))
	r.mux.Handle("/app/api/payment/pay", auth(http.HandlerFunc(paymentHandler.Pay)))

	// Robokassa callbacks (unauthenticated — called by Robokassa servers and browser redirects)
	r.mux.HandleFunc("/app/api/payment/robokassa/result", paymentHandler.RobokassaResult)
	r.mux.HandleFunc("/app/api/payment/robokassa/success", paymentHandler.RobokassaSuccess)
	r.mux.HandleFunc("/app/api/payment/robokassa/fail", paymentHandler.RobokassaFail)
	r.mux.Handle("/app/api/profile", auth(http.HandlerFunc(profileHandler.HandleProfile)))
	r.mux.Handle("/app/api/register", auth(http.HandlerFunc(registrationHandler.Register)))
	r.mux.Handle("/app/api/registration/status", auth(http.HandlerFunc(registrationHandler.Status)))

	// Dashboard
	r.mux.Handle("/app/api/dashboard", auth(http.HandlerFunc(dashboardHandler.GetDashboard)))

	// Programs
	r.mux.Handle("/app/api/programs", auth(http.HandlerFunc(workoutHandler.HandleProgramRoutes)))
	r.mux.Handle("/app/api/programs/", auth(http.HandlerFunc(workoutHandler.HandleProgramRoutes)))

	// Workouts
	r.mux.Handle("/app/api/workouts", auth(http.HandlerFunc(workoutHandler.HandleWorkoutRoutes)))
	r.mux.Handle("/app/api/workouts/", auth(http.HandlerFunc(workoutHandler.HandleWorkoutRoutes)))

	// Rehab
	r.mux.Handle("/app/api/rehab/courses", auth(http.HandlerFunc(rehabHandler.HandleRehabRoutes)))
	r.mux.Handle("/app/api/rehab/courses/", auth(http.HandlerFunc(rehabHandler.HandleRehabRoutes)))
	r.mux.Handle("/app/api/rehab/sessions/", auth(http.HandlerFunc(rehabHandler.HandleRehabSessionRoutes)))
	r.mux.Handle("/app/api/rehab/progress/", auth(http.HandlerFunc(rehabHandler.HandleRehabProgressRoutes)))

	// Nutrition
	r.mux.Handle("/app/api/nutrition/plans", auth(http.HandlerFunc(nutritionHandler.HandleNutritionRoutes)))
	r.mux.Handle("/app/api/nutrition/plans/", auth(http.HandlerFunc(nutritionHandler.HandleNutritionRoutes)))
	r.mux.Handle("/app/api/nutrition/calculator", auth(http.HandlerFunc(nutritionHandler.HandleNutritionRoutes)))
	r.mux.Handle("/app/api/food-log", auth(http.HandlerFunc(nutritionHandler.HandleFoodLogRoutes)))
	r.mux.Handle("/app/api/food-log/", auth(http.HandlerFunc(nutritionHandler.HandleFoodLogRoutes)))
	r.mux.Handle("/app/api/food-log/summary", auth(http.HandlerFunc(nutritionHandler.HandleFoodLogSummary)))

	// Progress
	r.mux.Handle("/app/api/progress", auth(http.HandlerFunc(progressV2Handler.HandleProgressRoutes)))
	r.mux.Handle("/app/api/progress/stats", auth(http.HandlerFunc(progressV2Handler.HandleProgressStats)))
	r.mux.Handle("/app/api/progress/achievements", auth(http.HandlerFunc(progressV2Handler.HandleProgressAchievements)))

	// Reviews
	r.mux.Handle("/app/api/reviews", auth(http.HandlerFunc(reviewHandler.HandleReviewRoutes)))
	r.mux.Handle("/app/api/reviews/", auth(http.HandlerFunc(reviewHandler.HandleReviewRoutes)))

	// Admin routes (authenticated + admin role required)
	adminHandler := NewAdminHandler(r.userSvc, r.workoutSvc, r.nutritionSvc, r.scoreSvc)
	adminAuth := func(h http.Handler) http.Handler {
		return auth(AdminMiddleware(h))
	}
	r.mux.Handle("/app/api/admin/users", adminAuth(http.HandlerFunc(adminHandler.HandleUserRoutes)))
	r.mux.Handle("/app/api/admin/users/", adminAuth(http.HandlerFunc(adminHandler.HandleUserRoutes)))
	r.mux.Handle("/app/api/admin/programs", adminAuth(http.HandlerFunc(adminHandler.HandleProgramRoutes)))
	r.mux.Handle("/app/api/admin/programs/", adminAuth(http.HandlerFunc(adminHandler.HandleProgramRoutes)))
	r.mux.Handle("/app/api/admin/workouts", adminAuth(http.HandlerFunc(adminHandler.HandleWorkoutRoutes)))
	r.mux.Handle("/app/api/admin/workouts/", adminAuth(http.HandlerFunc(adminHandler.HandleWorkoutRoutes)))
	r.mux.Handle("/app/api/admin/exercises", adminAuth(http.HandlerFunc(adminHandler.HandleExerciseRoutes)))
	r.mux.Handle("/app/api/admin/exercises/", adminAuth(http.HandlerFunc(adminHandler.HandleExerciseRoutes)))
	r.mux.Handle("/app/api/admin/meal-plans", adminAuth(http.HandlerFunc(adminHandler.HandleMealPlanRoutes)))
	r.mux.Handle("/app/api/admin/meal-plans/", adminAuth(http.HandlerFunc(adminHandler.HandleMealPlanRoutes)))
	r.mux.Handle("/app/api/admin/meals", adminAuth(http.HandlerFunc(adminHandler.HandleMealRoutes)))
	r.mux.Handle("/app/api/admin/meals/", adminAuth(http.HandlerFunc(adminHandler.HandleMealRoutes)))
	r.mux.Handle("/app/api/admin/workout-exercises", adminAuth(http.HandlerFunc(adminHandler.HandleWorkoutExerciseRoutes)))
	r.mux.Handle("/app/api/admin/reviews", adminAuth(http.HandlerFunc(adminHandler.GetReviewsSummary)))
	r.mux.Handle("/app/api/admin/stats", adminAuth(http.HandlerFunc(adminHandler.GetStats)))

	// SPA serving - serve static files, fallback to index.html
	r.mux.HandleFunc("/", r.serveSPA)
}

func (r *WebAppRouter) serveSPA(w http.ResponseWriter, req *http.Request) {
	// Don't serve SPA for API routes
	if strings.HasPrefix(req.URL.Path, "/app/api/") {
		http.NotFound(w, req)
		return
	}

	// Try to serve static file
	filePath := filepath.Join(r.staticDir, req.URL.Path)
	info, err := os.Stat(filePath)
	if err == nil && !info.IsDir() {
		http.ServeFile(w, req, filePath)
		return
	}

	// Check if the static directory has the assets subdirectory
	if strings.HasPrefix(req.URL.Path, "/assets/") {
		assetsPath := filepath.Join(r.staticDir, req.URL.Path)
		if _, err := os.Stat(assetsPath); err == nil {
			http.ServeFile(w, req, assetsPath)
			return
		}
	}

	// Fallback to index.html for Vue Router HTML5 history mode
	indexPath := filepath.Join(r.staticDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		// No static files available - serve a placeholder
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<!DOCTYPE html><html><body><p>Mini App not built yet. Run: cd web && npm run build</p></body></html>"))
		return
	}

	// Prevent stale caching of index.html in Telegram WebView
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	http.ServeFile(w, req, indexPath)
}

// StaticFileServer returns a handler that serves files from the static directory
func StaticFileServer(dir string) http.Handler {
	return http.FileServer(http.FS(os.DirFS(dir).(fs.FS)))
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}
