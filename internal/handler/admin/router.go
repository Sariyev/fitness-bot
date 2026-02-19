package admin

import (
	"encoding/json"
	"fitness-bot/internal/repository"
	"net/http"
	"strings"
)

type Router struct {
	mux     *http.ServeMux
	apiKey  string
	moduleH *ModuleAdminHandler
	questH  *QuestionnaireAdminHandler
	userH   *UserAdminHandler
}

func NewRouter(
	apiKey string,
	moduleRepo repository.ModuleRepository,
	questRepo repository.QuestionnaireRepository,
	userRepo repository.UserRepository,
	scoreRepo repository.ScoreRepository,
) *Router {
	r := &Router{
		mux:     http.NewServeMux(),
		apiKey:  apiKey,
		moduleH: NewModuleAdminHandler(moduleRepo),
		questH:  NewQuestionnaireAdminHandler(questRepo),
		userH:   NewUserAdminHandler(userRepo, scoreRepo),
	}
	r.routes()
	return r
}

func (r *Router) routes() {
	r.mux.HandleFunc("/api/modules", r.auth(r.moduleH.HandleModules))
	r.mux.HandleFunc("/api/modules/", r.auth(r.moduleH.HandleModuleByID))
	r.mux.HandleFunc("/api/categories/", r.auth(r.moduleH.HandleCategoryByID))
	r.mux.HandleFunc("/api/lessons/", r.auth(r.moduleH.HandleLessonByID))
	r.mux.HandleFunc("/api/contents/", r.auth(r.moduleH.HandleContentByID))

	r.mux.HandleFunc("/api/questionnaires", r.auth(r.questH.HandleQuestionnaires))
	r.mux.HandleFunc("/api/questionnaires/", r.auth(r.questH.HandleQuestionnaireByID))
	r.mux.HandleFunc("/api/questions/", r.auth(r.questH.HandleQuestionByID))
	r.mux.HandleFunc("/api/options/", r.auth(r.questH.HandleOptionByID))

	r.mux.HandleFunc("/api/users", r.auth(r.userH.HandleUsers))
	r.mux.HandleFunc("/api/users/", r.auth(r.userH.HandleUserByID))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if r.apiKey != "" {
			auth := req.Header.Get("Authorization")
			token := strings.TrimPrefix(auth, "Bearer ")
			if token != r.apiKey {
				jsonError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
		}
		next(w, req)
	}
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}

func decodeJSON(req *http.Request, v any) error {
	defer req.Body.Close()
	return json.NewDecoder(req.Body).Decode(v)
}
