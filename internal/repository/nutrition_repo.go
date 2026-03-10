package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type nutritionRepo struct {
	pool *pgxpool.Pool
}

func NewNutritionRepo(pool *pgxpool.Pool) NutritionRepository {
	return &nutritionRepo{pool: pool}
}

func (r *nutritionRepo) ListPlans(ctx context.Context, goal string) ([]models.MealPlan, error) {
	var query string
	var args []interface{}

	if goal != "" {
		query = `SELECT id, slug, name, goal, day_number, calories, protein, fat, carbs,
				 is_active, sort_order, created_at, updated_at
				 FROM meal_plans WHERE is_active = TRUE AND goal = $1
				 ORDER BY sort_order`
		args = append(args, goal)
	} else {
		query = `SELECT id, slug, name, goal, day_number, calories, protein, fat, carbs,
				 is_active, sort_order, created_at, updated_at
				 FROM meal_plans WHERE is_active = TRUE
				 ORDER BY sort_order`
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := []models.MealPlan{}
	for rows.Next() {
		var p models.MealPlan
		if err := rows.Scan(&p.ID, &p.Slug, &p.Name, &p.Goal, &p.DayNumber,
			&p.Calories, &p.Protein, &p.Fat, &p.Carbs,
			&p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}

func (r *nutritionRepo) GetPlanByID(ctx context.Context, id int) (*models.MealPlan, error) {
	p := &models.MealPlan{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, slug, name, goal, day_number, calories, protein, fat, carbs,
			is_active, sort_order, created_at, updated_at
		 FROM meal_plans WHERE id = $1`, id,
	).Scan(&p.ID, &p.Slug, &p.Name, &p.Goal, &p.DayNumber,
		&p.Calories, &p.Protein, &p.Fat, &p.Carbs,
		&p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *nutritionRepo) CreatePlan(ctx context.Context, p *models.MealPlan) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO meal_plans (slug, name, goal, day_number, calories, protein, fat, carbs, is_active, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, created_at, updated_at`,
		p.Slug, p.Name, p.Goal, p.DayNumber, p.Calories, p.Protein, p.Fat, p.Carbs, p.IsActive, p.SortOrder,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *nutritionRepo) UpdatePlan(ctx context.Context, p *models.MealPlan) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE meal_plans SET slug=$2, name=$3, goal=$4, day_number=$5, calories=$6,
			protein=$7, fat=$8, carbs=$9, is_active=$10, sort_order=$11, updated_at=NOW()
		 WHERE id=$1`,
		p.ID, p.Slug, p.Name, p.Goal, p.DayNumber, p.Calories,
		p.Protein, p.Fat, p.Carbs, p.IsActive, p.SortOrder)
	return err
}

func (r *nutritionRepo) ListMeals(ctx context.Context, planID int) ([]models.Meal, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, meal_plan_id, meal_type, name, recipe, calories, protein, fat, carbs,
			alternatives, sort_order, created_at, updated_at
		 FROM meals WHERE meal_plan_id = $1 ORDER BY sort_order`, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meals := []models.Meal{}
	for rows.Next() {
		var m models.Meal
		if err := rows.Scan(&m.ID, &m.MealPlanID, &m.MealType, &m.Name, &m.Recipe,
			&m.Calories, &m.Protein, &m.Fat, &m.Carbs,
			&m.Alternatives, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		meals = append(meals, m)
	}
	return meals, nil
}

func (r *nutritionRepo) CreateMeal(ctx context.Context, m *models.Meal) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO meals (meal_plan_id, meal_type, name, recipe, calories, protein, fat, carbs,
			alternatives, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, created_at, updated_at`,
		m.MealPlanID, m.MealType, m.Name, m.Recipe, m.Calories, m.Protein, m.Fat, m.Carbs,
		m.Alternatives, m.SortOrder,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *nutritionRepo) UpdateMeal(ctx context.Context, m *models.Meal) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE meals SET meal_plan_id=$2, meal_type=$3, name=$4, recipe=$5, calories=$6,
			protein=$7, fat=$8, carbs=$9, alternatives=$10, sort_order=$11, updated_at=NOW()
		 WHERE id=$1`,
		m.ID, m.MealPlanID, m.MealType, m.Name, m.Recipe, m.Calories,
		m.Protein, m.Fat, m.Carbs, m.Alternatives, m.SortOrder)
	return err
}
