package webapp

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fitness-bot/internal/service"
)

type WebAppRouter struct {
	mux        *http.ServeMux
	moduleSvc  *service.ModuleService
	userSvc    *service.UserService
	paymentSvc *service.PaymentService
	botToken   string
	staticDir  string
}

func NewRouter(
	botToken string,
	userSvc *service.UserService,
	moduleSvc *service.ModuleService,
	paymentSvc *service.PaymentService,
	staticDir string,
) *WebAppRouter {
	r := &WebAppRouter{
		mux:        http.NewServeMux(),
		moduleSvc:  moduleSvc,
		userSvc:    userSvc,
		paymentSvc: paymentSvc,
		botToken:   botToken,
		staticDir:  staticDir,
	}

	r.setupRoutes()
	return r
}

func (r *WebAppRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Set CORS headers for Telegram Mini App
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Telegram-Init-Data")

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
	paymentHandler := NewPaymentHandler(r.paymentSvc)
	profileHandler := NewProfileHandler(r.userSvc)
	registrationHandler := NewRegistrationHandler(r.userSvc)

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
	r.mux.Handle("/app/api/profile", auth(http.HandlerFunc(profileHandler.GetProfile)))
	r.mux.Handle("/app/api/register", auth(http.HandlerFunc(registrationHandler.Register)))
	r.mux.Handle("/app/api/registration/status", auth(http.HandlerFunc(registrationHandler.Status)))

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
