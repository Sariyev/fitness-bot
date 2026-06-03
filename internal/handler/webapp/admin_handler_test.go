package webapp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
	"fitness-bot/internal/service"
)

// errNotFound is a stand-in for "row not found". The admin handler only
// branches on (err != nil), so the exact value doesn't matter here.
var errNotFound = &notFoundErr{}

type notFoundErr struct{}

func (n *notFoundErr) Error() string { return "not found" }

// ============================================================================
//  Fakes (only the methods exercised by these tests have real bodies; the
//  rest exist to satisfy the interface and return zero values).
// ============================================================================

type fakeWorkoutRepo struct {
	store      map[int]*models.Workout
	nextID     int
	lastCreate *models.Workout
}

func newFakeWorkoutRepo() *fakeWorkoutRepo {
	return &fakeWorkoutRepo{store: map[int]*models.Workout{}, nextID: 100}
}

func (r *fakeWorkoutRepo) ListWorkouts(context.Context, string, string, string) ([]models.Workout, error) {
	return nil, nil
}
func (r *fakeWorkoutRepo) ListAllWorkouts(context.Context) ([]models.Workout, error) {
	out := []models.Workout{}
	for _, w := range r.store {
		out = append(out, *w)
	}
	return out, nil
}
func (r *fakeWorkoutRepo) GetWorkoutByID(_ context.Context, id int) (*models.Workout, error) {
	if w, ok := r.store[id]; ok {
		return w, nil
	}
	return nil, errNotFound
}
func (r *fakeWorkoutRepo) CreateWorkout(_ context.Context, w *models.Workout) error {
	w.ID = r.nextID
	r.nextID++
	now := time.Now()
	w.CreatedAt, w.UpdatedAt = now, now
	cp := *w
	r.store[w.ID] = &cp
	r.lastCreate = &cp
	return nil
}
func (r *fakeWorkoutRepo) UpdateWorkout(context.Context, *models.Workout) error { return nil }
func (r *fakeWorkoutRepo) DeleteWorkout(_ context.Context, id int) error {
	delete(r.store, id)
	return nil
}
func (r *fakeWorkoutRepo) ListExercises(context.Context, int) ([]models.WorkoutExercise, error) {
	return nil, nil
}
func (r *fakeWorkoutRepo) AddExercise(context.Context, *models.WorkoutExercise) error { return nil }

type fakeExerciseRepo struct{}

func (r *fakeExerciseRepo) List(context.Context) ([]models.Exercise, error) { return nil, nil }
func (r *fakeExerciseRepo) GetByID(context.Context, int) (*models.Exercise, error) {
	return nil, errNotFound
}
func (r *fakeExerciseRepo) Create(_ context.Context, e *models.Exercise) error {
	e.ID = 1
	return nil
}
func (r *fakeExerciseRepo) Update(context.Context, *models.Exercise) error { return nil }
func (r *fakeExerciseRepo) Delete(context.Context, int) error              { return nil }

type fakeDailyCompletionRepo struct{}

func (r *fakeDailyCompletionRepo) Create(context.Context, *models.DailyCompletion) error {
	return nil
}
func (r *fakeDailyCompletionRepo) ListByDate(context.Context, int64, string) ([]models.DailyCompletion, error) {
	return nil, nil
}
func (r *fakeDailyCompletionRepo) GetStreak(context.Context, int64) (int, int, error) {
	return 0, 0, nil
}
func (r *fakeDailyCompletionRepo) GetCalendar(context.Context, int64, int, int) ([]string, error) {
	return nil, nil
}

type fakeRehabRepo struct {
	courses    map[int]*models.RehabCourse
	sessions   map[int]*models.RehabSession
	nextCourse int
	nextSess   int
	lastCourse *models.RehabCourse
}

func newFakeRehabRepo() *fakeRehabRepo {
	return &fakeRehabRepo{
		courses: map[int]*models.RehabCourse{}, sessions: map[int]*models.RehabSession{},
		nextCourse: 100, nextSess: 100,
	}
}
func (r *fakeRehabRepo) ListCourses(context.Context, string) ([]models.RehabCourse, error) {
	return nil, nil
}
func (r *fakeRehabRepo) ListAllCourses(context.Context) ([]models.RehabCourse, error) {
	out := []models.RehabCourse{}
	for _, c := range r.courses {
		out = append(out, *c)
	}
	return out, nil
}
func (r *fakeRehabRepo) GetCourseByID(_ context.Context, id int) (*models.RehabCourse, error) {
	if c, ok := r.courses[id]; ok {
		return c, nil
	}
	return nil, errNotFound
}
func (r *fakeRehabRepo) CreateCourse(_ context.Context, c *models.RehabCourse) error {
	c.ID = r.nextCourse
	r.nextCourse++
	now := time.Now()
	c.CreatedAt, c.UpdatedAt = now, now
	cp := *c
	r.courses[c.ID] = &cp
	r.lastCourse = &cp
	return nil
}
func (r *fakeRehabRepo) DeleteCourse(_ context.Context, id int) error {
	delete(r.courses, id)
	// Cascade to child sessions — mirrors the real repo's transaction.
	for sid, sess := range r.sessions {
		if sess.CourseID == id {
			delete(r.sessions, sid)
		}
	}
	return nil
}
func (r *fakeRehabRepo) DeleteSession(_ context.Context, id int) error {
	delete(r.sessions, id)
	return nil
}
func (r *fakeRehabRepo) UpdateCourse(_ context.Context, c *models.RehabCourse) error {
	if _, ok := r.courses[c.ID]; !ok {
		return errNotFound
	}
	cp := *c
	r.courses[c.ID] = &cp
	return nil
}
func (r *fakeRehabRepo) ListSessions(context.Context, int) ([]models.RehabSession, error) {
	return nil, nil
}
func (r *fakeRehabRepo) GetSessionByID(_ context.Context, id int) (*models.RehabSession, error) {
	if s, ok := r.sessions[id]; ok {
		return s, nil
	}
	return nil, errNotFound
}
func (r *fakeRehabRepo) CreateSession(_ context.Context, s *models.RehabSession) error {
	s.ID = r.nextSess
	r.nextSess++
	now := time.Now()
	s.CreatedAt, s.UpdatedAt = now, now
	cp := *s
	r.sessions[s.ID] = &cp
	return nil
}
func (r *fakeRehabRepo) UpdateSession(context.Context, *models.RehabSession) error { return nil }
func (r *fakeRehabRepo) CreateProgress(context.Context, *models.UserRehabProgress) error {
	return nil
}
func (r *fakeRehabRepo) ListUserProgress(context.Context, int64, int) ([]models.UserRehabProgress, error) {
	return nil, nil
}

type fakeNutritionRepo struct {
	plans    map[int]*models.MealPlan
	meals    map[int]*models.Meal
	nextP    int
	nextM    int
	lastPlan *models.MealPlan
	lastMeal *models.Meal
}

func newFakeNutritionRepo() *fakeNutritionRepo {
	return &fakeNutritionRepo{
		plans: map[int]*models.MealPlan{}, meals: map[int]*models.Meal{},
		nextP: 100, nextM: 100,
	}
}
func (r *fakeNutritionRepo) ListPlans(context.Context, string) ([]models.MealPlan, error) {
	return nil, nil
}
func (r *fakeNutritionRepo) ListAllPlans(context.Context) ([]models.MealPlan, error) {
	out := []models.MealPlan{}
	for _, p := range r.plans {
		out = append(out, *p)
	}
	return out, nil
}
func (r *fakeNutritionRepo) GetPlanByID(_ context.Context, id int) (*models.MealPlan, error) {
	if p, ok := r.plans[id]; ok {
		return p, nil
	}
	return nil, errNotFound
}
func (r *fakeNutritionRepo) CreatePlan(_ context.Context, p *models.MealPlan) error {
	p.ID = r.nextP
	r.nextP++
	cp := *p
	r.plans[p.ID] = &cp
	r.lastPlan = &cp
	return nil
}
func (r *fakeNutritionRepo) UpdatePlan(context.Context, *models.MealPlan) error { return nil }
func (r *fakeNutritionRepo) DeletePlan(_ context.Context, id int) error {
	delete(r.plans, id)
	// Cascade to child meals so test fakes mirror the real repo's transaction.
	for mid, m := range r.meals {
		if m.MealPlanID == id {
			delete(r.meals, mid)
		}
	}
	return nil
}
func (r *fakeNutritionRepo) ListMeals(context.Context, int) ([]models.Meal, error) {
	return nil, nil
}
func (r *fakeNutritionRepo) CreateMeal(_ context.Context, m *models.Meal) error {
	m.ID = r.nextM
	r.nextM++
	cp := *m
	r.meals[m.ID] = &cp
	r.lastMeal = &cp
	return nil
}
func (r *fakeNutritionRepo) UpdateMeal(context.Context, *models.Meal) error { return nil }
func (r *fakeNutritionRepo) DeleteMeal(_ context.Context, id int) error {
	delete(r.meals, id)
	return nil
}
func (r *fakeNutritionRepo) GetMealByID(_ context.Context, id int) (*models.Meal, error) {
	if m, ok := r.meals[id]; ok {
		return m, nil
	}
	return nil, errNotFound
}

type fakeFoodLogRepo struct{}

func (r *fakeFoodLogRepo) Create(context.Context, *models.FoodLogEntry) error { return nil }
func (r *fakeFoodLogRepo) Delete(context.Context, int64, int64) error          { return nil }
func (r *fakeFoodLogRepo) ListByDate(context.Context, int64, string) ([]models.FoodLogEntry, error) {
	return nil, nil
}
func (r *fakeFoodLogRepo) GetDailySummary(context.Context, int64, string) (int, float64, float64, float64, error) {
	return 0, 0, 0, 0, nil
}

type fakeScoreRepo struct{}

func (r *fakeScoreRepo) Create(context.Context, *models.UserScore) error { return nil }
func (r *fakeScoreRepo) GetByReference(context.Context, int64, string, int) (*models.UserScore, error) {
	return nil, errNotFound
}
func (r *fakeScoreRepo) ListByUser(context.Context, int64) ([]models.UserScore, error) {
	return nil, nil
}
func (r *fakeScoreRepo) ListByReference(context.Context, string, int) ([]models.UserScore, error) {
	return nil, nil
}
func (r *fakeScoreRepo) GetSummary(context.Context, string, int) (*models.ReviewSummary, error) {
	return &models.ReviewSummary{}, nil
}
func (r *fakeScoreRepo) GetBotSummary(context.Context) (*models.ReviewSummary, error) {
	return &models.ReviewSummary{}, nil
}

type fakeUserRepo struct{ admin *models.User }

func (r *fakeUserRepo) Create(context.Context, *models.User) error { return nil }
func (r *fakeUserRepo) GetByID(_ context.Context, id int64) (*models.User, error) {
	if r.admin != nil && r.admin.ID == id {
		return r.admin, nil
	}
	return nil, errNotFound
}
func (r *fakeUserRepo) GetByTelegramID(_ context.Context, tid int64) (*models.User, error) {
	if r.admin != nil && r.admin.TelegramID == tid {
		return r.admin, nil
	}
	return nil, errNotFound
}
func (r *fakeUserRepo) Update(context.Context, *models.User) error { return nil }
func (r *fakeUserRepo) ListAll(context.Context, int, int) ([]models.User, int, error) {
	return nil, 0, nil
}
func (r *fakeUserRepo) SetAvatarMediaID(context.Context, int64, *int64) error { return nil }
func (r *fakeUserRepo) CreateProfile(context.Context, *models.UserProfile) error {
	return nil
}
func (r *fakeUserRepo) GetProfileByUserID(context.Context, int64) (*models.UserProfile, error) {
	return nil, errNotFound
}
func (r *fakeUserRepo) UpdateProfile(context.Context, *models.UserProfile) error { return nil }
func (r *fakeUserRepo) ListReminderTargets(context.Context, string) ([]repository.ReminderTarget, error) {
	return nil, nil
}
func (r *fakeUserRepo) MarkReminderSent(context.Context, int64) error { return nil }

type fakePricingRepoH struct{ prices map[models.Category]int }

func (r *fakePricingRepoH) GetPrice(_ context.Context, c models.Category) (int, error) {
	if p, ok := r.prices[c]; ok {
		return p, nil
	}
	return 0, errNotFound
}
func (r *fakePricingRepoH) ListPrices(context.Context) (map[models.Category]int, error) {
	out := map[models.Category]int{}
	for k, v := range r.prices {
		out[k] = v
	}
	return out, nil
}
func (r *fakePricingRepoH) SetPrice(_ context.Context, c models.Category, p int) error {
	r.prices[c] = p
	return nil
}

type fakeAccessRepoH struct{ granted map[int64]map[models.Category]bool }

func (r *fakeAccessRepoH) HasAccess(_ context.Context, uid int64, c models.Category) (bool, error) {
	return r.granted[uid][c], nil
}
func (r *fakeAccessRepoH) Grant(_ context.Context, uid int64, c models.Category, _ *int64) error {
	if r.granted[uid] == nil {
		r.granted[uid] = map[models.Category]bool{}
	}
	r.granted[uid][c] = true
	return nil
}
func (r *fakeAccessRepoH) ListGranted(context.Context, int64) ([]models.Category, error) {
	return nil, nil
}

// ============================================================================
//  Test harness
// ============================================================================

type adminTestSetup struct {
	handler  *AdminHandler
	admin    *models.User
	workouts *fakeWorkoutRepo
	rehab    *fakeRehabRepo
	nutr     *fakeNutritionRepo
}

func newAdminTestSetup(t *testing.T) *adminTestSetup {
	t.Helper()
	admin := &models.User{ID: 17, TelegramID: 525578774, Role: "admin", IsRegistered: true}

	workouts := newFakeWorkoutRepo()
	rehab := newFakeRehabRepo()
	nutr := newFakeNutritionRepo()

	userSvc := service.NewUserService(&fakeUserRepo{admin: admin})
	workoutSvc := service.NewWorkoutService(workouts, &fakeExerciseRepo{}, &fakeDailyCompletionRepo{})
	rehabSvc := service.NewRehabService(rehab)
	nutrSvc := service.NewNutritionService(nutr, &fakeFoodLogRepo{})
	scoreSvc := service.NewScoreService(&fakeScoreRepo{})
	accessSvc := service.NewAccessService(
		&fakePricingRepoH{prices: map[models.Category]int{
			models.CategoryWorkouts: 5000, models.CategoryLFK: 5000, models.CategoryNutrition: 5000,
		}},
		&fakeAccessRepoH{granted: map[int64]map[models.Category]bool{}},
	)

	h := NewAdminHandler(userSvc, workoutSvc, rehabSvc, nutrSvc, scoreSvc, accessSvc)
	return &adminTestSetup{handler: h, admin: admin, workouts: workouts, rehab: rehab, nutr: nutr}
}

// doAs invokes the handler method with `admin` in the request context,
// bypassing AuthMiddleware (not under test here).
func (s *adminTestSetup) doAs(method, path string, body any, fn func(http.ResponseWriter, *http.Request)) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req = req.WithContext(context.WithValue(req.Context(), userContextKey, s.admin))
	rec := httptest.NewRecorder()
	fn(rec, req)
	return rec
}

// ============================================================================
//  Tests
// ============================================================================

// The "Program" entity was removed in migration 028 — admin only edits
// workouts directly now. These tests target createWorkout instead, which
// inherited the slug/access_tier/name validation that used to live on
// the program create path.

func TestCreateWorkout_EmptySlugAutoFills(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{
		"name":             "Test Workout",
		"slug":             "",
		"goal":             "weight_loss",
		"format":           "home",
		"level":            "beginner",
		"duration_minutes": 30,
		"access_tier":      "paid",
		"is_active":        false,
	}
	rec := s.doAs(http.MethodPost, "/app/api/admin/workouts", body, s.handler.createWorkout)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if s.workouts.lastCreate == nil {
		t.Fatal("repo CreateWorkout never called")
	}
	if s.workouts.lastCreate.Slug == "" {
		t.Errorf("expected auto-generated slug, got empty string")
	}
}

func TestCreateWorkout_TwoEmptySlugsDontCollide(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "Same Name", "access_tier": "paid"}

	rec1 := s.doAs(http.MethodPost, "/app/api/admin/workouts", body, s.handler.createWorkout)
	if rec1.Code != http.StatusCreated {
		t.Fatalf("first create: expected 201, got %d", rec1.Code)
	}
	slug1 := s.workouts.lastCreate.Slug

	time.Sleep(2 * time.Millisecond)

	rec2 := s.doAs(http.MethodPost, "/app/api/admin/workouts", body, s.handler.createWorkout)
	if rec2.Code != http.StatusCreated {
		t.Fatalf("second create: expected 201, got %d", rec2.Code)
	}
	slug2 := s.workouts.lastCreate.Slug

	if slug1 == slug2 {
		t.Errorf("auto-generated slugs collided: %q == %q", slug1, slug2)
	}
}

func TestCreateWorkout_BogusAccessTierIs400(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "X", "access_tier": "freemium"}
	rec := s.doAs(http.MethodPost, "/app/api/admin/workouts", body, s.handler.createWorkout)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bogus access_tier, got %d: %s", rec.Code, rec.Body.String())
	}
	if s.workouts.lastCreate != nil {
		t.Error("repo CreateWorkout should NOT have been called for invalid access_tier")
	}
}

func TestCreateWorkout_EmptyAccessTierAccepted(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "X", "access_tier": ""}
	rec := s.doAs(http.MethodPost, "/app/api/admin/workouts", body, s.handler.createWorkout)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201 for empty access_tier, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateWorkout_EmptyNameIs400(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "", "access_tier": "paid"}
	rec := s.doAs(http.MethodPost, "/app/api/admin/workouts", body, s.handler.createWorkout)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty name, got %d", rec.Code)
	}
}

func TestCreateMeal_ZeroPlanIDIs400(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "Test Meal", "meal_plan_id": 0, "meal_type": "breakfast"}
	rec := s.doAs(http.MethodPost, "/app/api/admin/meals", body, s.handler.createMeal)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for meal_plan_id=0, got %d: %s", rec.Code, rec.Body.String())
	}
	if s.nutr.lastMeal != nil {
		t.Error("repo CreateMeal should not have been called for invalid meal_plan_id")
	}
}

func TestCreateRehabCourse_EmptySlugAutoFills(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "Rehab Course", "category": "back", "access_tier": "paid"}
	rec := s.doAs(http.MethodPost, "/app/api/admin/rehab/courses", body, s.handler.createRehabCourse)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if s.rehab.lastCourse == nil || s.rehab.lastCourse.Slug == "" {
		t.Error("expected auto-generated slug")
	}
}

func TestCreateRehabCourse_BogusAccessTierIs400(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "X", "access_tier": "premium"}
	rec := s.doAs(http.MethodPost, "/app/api/admin/rehab/courses", body, s.handler.createRehabCourse)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestPricing_GetAndSet(t *testing.T) {
	s := newAdminTestSetup(t)

	rec := s.doAs(http.MethodGet, "/app/api/admin/pricing", nil, s.handler.HandlePricingRoutes)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var prices map[string]int
	_ = json.Unmarshal(rec.Body.Bytes(), &prices)
	if prices["workouts"] != 5000 {
		t.Errorf("expected workouts=5000, got %v", prices)
	}

	// PUT updates the price.
	rec2 := s.doAs(http.MethodPut, "/app/api/admin/pricing/workouts",
		map[string]any{"price_kzt": 7500}, s.handler.HandlePricingRoutes)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec2.Code, rec2.Body.String())
	}

	rec3 := s.doAs(http.MethodGet, "/app/api/admin/pricing", nil, s.handler.HandlePricingRoutes)
	var got map[string]int
	_ = json.Unmarshal(rec3.Body.Bytes(), &got)
	if got["workouts"] != 7500 {
		t.Errorf("after PUT, expected workouts=7500, got %v", got)
	}
}

func TestPricing_RejectsNonPositive(t *testing.T) {
	s := newAdminTestSetup(t)
	rec := s.doAs(http.MethodPut, "/app/api/admin/pricing/workouts",
		map[string]any{"price_kzt": 0}, s.handler.HandlePricingRoutes)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for price_kzt=0, got %d", rec.Code)
	}
}

func TestPricing_RejectsInvalidCategory(t *testing.T) {
	s := newAdminTestSetup(t)
	rec := s.doAs(http.MethodPut, "/app/api/admin/pricing/bogus",
		map[string]any{"price_kzt": 1000}, s.handler.HandlePricingRoutes)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid category, got %d", rec.Code)
	}
}

func TestValidAccessTier(t *testing.T) {
	cases := map[models.AccessTier]bool{
		"":                 true, // repo layer defaults empty to 'paid'
		models.AccessFree:  true,
		models.AccessTrial: true,
		models.AccessPaid:  true,
		"freemium":         false,
		"premium":          false,
		"PAID":             false,
	}
	for tier, want := range cases {
		if got := validAccessTier(tier); got != want {
			t.Errorf("validAccessTier(%q) = %v, want %v", tier, got, want)
		}
	}
}

// ============================================================================
//  Range guard tests — both the helper math and the wired handler behavior
// ============================================================================

func TestGuardRangeInt(t *testing.T) {
	cases := []struct {
		name             string
		val, lo, hi      int
		wantErrSubstring string
	}{
		{"in-bounds", 5, 0, 10, ""},
		{"on-lower-bound", 0, 0, 10, ""},
		{"on-upper-bound", 10, 0, 10, ""},
		{"below-min", -1, 0, 10, "out of range"},
		{"above-max", 11, 0, 10, "out of range"},
	}
	for _, c := range cases {
		got := guardRangeInt("field", c.val, c.lo, c.hi)
		if c.wantErrSubstring == "" && got != "" {
			t.Errorf("%s: expected no error, got %q", c.name, got)
		}
		if c.wantErrSubstring != "" && !strings.Contains(got, c.wantErrSubstring) {
			t.Errorf("%s: expected error containing %q, got %q", c.name, c.wantErrSubstring, got)
		}
	}
}

func TestGuardRangeFloat(t *testing.T) {
	if guardRangeFloat("x", 5.0, 0, 10) != "" {
		t.Error("5.0 in [0,10] should be valid")
	}
	if guardRangeFloat("x", -0.01, 0, 10) == "" {
		t.Error("-0.01 should fail [0,10] guard")
	}
	if guardRangeFloat("x", 10.0001, 0, 10) == "" {
		t.Error("10.0001 should fail [0,10] guard")
	}
}

func TestCreateWorkout_DurationMinutesOutOfRange(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "Wk", "access_tier": "paid", "duration_minutes": 9999}
	rec := s.doAs(http.MethodPost, "/app/api/admin/workouts", body, s.handler.createWorkout)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for duration_minutes=9999, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateWorkout_SortOrderNegativeIs400(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "Wk", "access_tier": "paid", "sort_order": -1}
	rec := s.doAs(http.MethodPost, "/app/api/admin/workouts", body, s.handler.createWorkout)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for sort_order=-1, got %d", rec.Code)
	}
}

func TestCreateMealPlan_CaloriesOutOfRange(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "Plan", "access_tier": "paid", "calories": 99999}
	rec := s.doAs(http.MethodPost, "/app/api/admin/meal-plans", body, s.handler.createMealPlan)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for calories=99999, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateMealPlan_MacrosNegative(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "Plan", "access_tier": "paid", "protein": -10}
	rec := s.doAs(http.MethodPost, "/app/api/admin/meal-plans", body, s.handler.createMealPlan)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for protein=-10, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateMeal_CaloriesOutOfRange(t *testing.T) {
	s := newAdminTestSetup(t)
	// Seed a valid plan so meal_plan_id=N doesn't 400 first.
	_ = s.nutr.CreatePlan(context.Background(), &models.MealPlan{Name: "P"})
	planID := s.nutr.lastPlan.ID

	body := map[string]any{"name": "M", "meal_plan_id": planID, "calories": 50000}
	rec := s.doAs(http.MethodPost, "/app/api/admin/meals", body, s.handler.createMeal)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for calories=50000, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateRehabSession_StageOutOfRange(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"course_id": 1, "stage": 9}
	rec := s.doAs(http.MethodPost, "/app/api/admin/rehab/sessions", body, s.handler.createRehabSession)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for stage=9, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateRehabSession_DurationOutOfRange(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"course_id": 1, "duration_minutes": 9999}
	rec := s.doAs(http.MethodPost, "/app/api/admin/rehab/sessions", body, s.handler.createRehabSession)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for duration_minutes=9999, got %d: %s", rec.Code, rec.Body.String())
	}
}

// ============================================================================
//  image_media_id wiring — meal plan + meal round-trip
// ============================================================================

func TestCreateMealPlan_ImageMediaIDRoundTrip(t *testing.T) {
	s := newAdminTestSetup(t)
	var imageID int64 = 42
	body := map[string]any{
		"name":           "P",
		"access_tier":    "paid",
		"image_media_id": imageID,
	}
	rec := s.doAs(http.MethodPost, "/app/api/admin/meal-plans", body, s.handler.createMealPlan)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if s.nutr.lastPlan == nil {
		t.Fatal("CreatePlan never called")
	}
	if s.nutr.lastPlan.ImageMediaID == nil {
		t.Fatal("ImageMediaID dropped on the way to repo")
	}
	if *s.nutr.lastPlan.ImageMediaID != imageID {
		t.Errorf("ImageMediaID = %d, want %d", *s.nutr.lastPlan.ImageMediaID, imageID)
	}
}

func TestCreateMealPlan_NoImageMediaID(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "P", "access_tier": "paid"}
	rec := s.doAs(http.MethodPost, "/app/api/admin/meal-plans", body, s.handler.createMealPlan)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
	if s.nutr.lastPlan.ImageMediaID != nil {
		t.Errorf("expected nil ImageMediaID, got %v", *s.nutr.lastPlan.ImageMediaID)
	}
}

// ============================================================================
//  Delete behavior — including cascade from parent → child for the
//  meal-plan and rehab-course transactional deletes.
// ============================================================================

func TestDeleteWorkout_RemovesFromStore(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "W", "access_tier": "paid"}
	rec := s.doAs(http.MethodPost, "/app/api/admin/workouts", body, s.handler.createWorkout)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: %d", rec.Code)
	}
	id := s.workouts.lastCreate.ID

	rec2 := s.doAs(http.MethodDelete, "/app/api/admin/workouts", nil, func(w http.ResponseWriter, r *http.Request) {
		s.handler.deleteWorkout(w, r, id)
	})
	if rec2.Code != http.StatusOK {
		t.Errorf("delete: %d: %s", rec2.Code, rec2.Body.String())
	}
	if _, ok := s.workouts.store[id]; ok {
		t.Error("workout still in store after delete")
	}
}

func TestDeleteMealPlan_CascadesToMeals(t *testing.T) {
	s := newAdminTestSetup(t)
	_ = s.nutr.CreatePlan(context.Background(), &models.MealPlan{Name: "P"})
	planID := s.nutr.lastPlan.ID
	_ = s.nutr.CreateMeal(context.Background(), &models.Meal{Name: "M1", MealPlanID: planID})
	_ = s.nutr.CreateMeal(context.Background(), &models.Meal{Name: "M2", MealPlanID: planID})
	if len(s.nutr.meals) != 2 {
		t.Fatalf("expected 2 meals seeded, got %d", len(s.nutr.meals))
	}

	rec := s.doAs(http.MethodDelete, "/app/api/admin/meal-plans", nil, func(w http.ResponseWriter, r *http.Request) {
		s.handler.deleteMealPlan(w, r, planID)
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("delete: %d", rec.Code)
	}
	if _, ok := s.nutr.plans[planID]; ok {
		t.Error("plan still present")
	}
	if len(s.nutr.meals) != 0 {
		t.Errorf("expected cascade-delete of child meals, %d remain", len(s.nutr.meals))
	}
}

func TestDeleteRehabCourse_CascadesToSessions(t *testing.T) {
	s := newAdminTestSetup(t)
	body := map[string]any{"name": "C", "category": "back", "access_tier": "paid"}
	rec := s.doAs(http.MethodPost, "/app/api/admin/rehab/courses", body, s.handler.createRehabCourse)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create course: %d", rec.Code)
	}
	courseID := s.rehab.lastCourse.ID
	// Seed two sessions belonging to the course.
	_ = s.rehab.CreateSession(context.Background(), &models.RehabSession{CourseID: courseID, Stage: 1, DayNumber: 1})
	_ = s.rehab.CreateSession(context.Background(), &models.RehabSession{CourseID: courseID, Stage: 1, DayNumber: 2})

	rec2 := s.doAs(http.MethodDelete, "/app/api/admin/rehab/courses", nil, func(w http.ResponseWriter, r *http.Request) {
		s.handler.deleteRehabCourse(w, r, courseID)
	})
	if rec2.Code != http.StatusOK {
		t.Fatalf("delete: %d", rec2.Code)
	}
	if _, ok := s.rehab.courses[courseID]; ok {
		t.Error("course still present")
	}
	for _, sess := range s.rehab.sessions {
		if sess.CourseID == courseID {
			t.Error("child session not cascade-deleted")
		}
	}
}

func TestDeleteWorkout_OK(t *testing.T) {
	s := newAdminTestSetup(t)
	rec := s.doAs(http.MethodDelete, "/app/api/admin/workouts", nil, func(w http.ResponseWriter, r *http.Request) {
		s.handler.deleteWorkout(w, r, 999)
	})
	// Fake repo's DeleteWorkout returns nil unconditionally; the route still
	// must respond 200. Real repo deletes workout_exercises + the workout
	// in a transaction.
	if rec.Code != http.StatusOK {
		t.Errorf("delete: %d", rec.Code)
	}
}

func TestCreateMeal_ImageMediaIDRoundTrip(t *testing.T) {
	s := newAdminTestSetup(t)
	_ = s.nutr.CreatePlan(context.Background(), &models.MealPlan{Name: "P"})
	planID := s.nutr.lastPlan.ID

	var imageID int64 = 77
	body := map[string]any{
		"name":           "M",
		"meal_plan_id":   planID,
		"image_media_id": imageID,
	}
	rec := s.doAs(http.MethodPost, "/app/api/admin/meals", body, s.handler.createMeal)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if s.nutr.lastMeal == nil || s.nutr.lastMeal.ImageMediaID == nil {
		t.Fatal("ImageMediaID dropped on the way to repo")
	}
	if *s.nutr.lastMeal.ImageMediaID != imageID {
		t.Errorf("ImageMediaID = %d, want %d", *s.nutr.lastMeal.ImageMediaID, imageID)
	}
}
