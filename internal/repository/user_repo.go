package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) UserRepository {
	return &userRepo{pool: pool}
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO users (telegram_id, username, first_name, last_name, language_code, is_registered)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (telegram_id) DO UPDATE SET username = EXCLUDED.username
		 RETURNING id, created_at, updated_at`,
		user.TelegramID, user.Username, user.FirstName, user.LastName, user.LanguageCode, user.IsRegistered,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepo) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	u := &models.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, telegram_id, username, first_name, last_name, language_code, is_registered, created_at, updated_at
		 FROM users WHERE telegram_id = $1`, telegramID,
	).Scan(&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName,
		&u.LanguageCode, &u.IsRegistered, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET username=$2, first_name=$3, last_name=$4, is_registered=$5, updated_at=NOW()
		 WHERE id=$1`,
		user.ID, user.Username, user.FirstName, user.LastName, user.IsRegistered,
	)
	return err
}

func (r *userRepo) CreateProfile(ctx context.Context, p *models.UserProfile) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO user_profiles (user_id, gender, age, weight_kg, height_cm, fitness_level, goal)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at, updated_at`,
		p.UserID, p.Gender, p.Age, p.WeightKg, p.HeightCm, p.FitnessLevel, p.Goal,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *userRepo) GetProfileByUserID(ctx context.Context, userID int64) (*models.UserProfile, error) {
	p := &models.UserProfile{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, gender, age, weight_kg, height_cm, fitness_level, goal, created_at, updated_at
		 FROM user_profiles WHERE user_id = $1`, userID,
	).Scan(&p.ID, &p.UserID, &p.Gender, &p.Age, &p.WeightKg, &p.HeightCm,
		&p.FitnessLevel, &p.Goal, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}
