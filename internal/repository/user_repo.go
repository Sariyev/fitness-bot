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
	if user.Role == "" {
		user.Role = models.RoleClient
	}
	return r.pool.QueryRow(ctx,
		`INSERT INTO users (telegram_id, username, first_name, last_name, language_code, is_registered, role)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (telegram_id) DO UPDATE SET username = EXCLUDED.username
		 RETURNING id, is_paid, role, created_at, updated_at`,
		user.TelegramID, user.Username, user.FirstName, user.LastName, user.LanguageCode, user.IsRegistered, user.Role,
	).Scan(&user.ID, &user.IsPaid, &user.Role, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	u := &models.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, telegram_id, username, first_name, last_name, language_code, is_registered, is_paid, role, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName,
		&u.LanguageCode, &u.IsRegistered, &u.IsPaid, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	u := &models.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, telegram_id, username, first_name, last_name, language_code, is_registered, is_paid, role, created_at, updated_at
		 FROM users WHERE telegram_id = $1`, telegramID,
	).Scan(&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName,
		&u.LanguageCode, &u.IsRegistered, &u.IsPaid, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET username=$2, first_name=$3, last_name=$4, is_registered=$5, is_paid=$6, role=$7, updated_at=NOW()
		 WHERE id=$1`,
		user.ID, user.Username, user.FirstName, user.LastName, user.IsRegistered, user.IsPaid, user.Role,
	)
	return err
}

func (r *userRepo) ListAll(ctx context.Context, limit, offset int) ([]models.User, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, telegram_id, username, first_name, last_name, language_code,
				is_registered, is_paid, role, created_at, updated_at
		 FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.TelegramID, &u.Username, &u.FirstName, &u.LastName,
			&u.LanguageCode, &u.IsRegistered, &u.IsPaid, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *userRepo) CreateProfile(ctx context.Context, p *models.UserProfile) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO user_profiles (
			user_id, gender, age, weight_kg, height_cm, fitness_level, goal,
			training_access, training_experience, has_pain, pain_locations, pain_level,
			diagnoses, contraindications, days_per_week, session_duration, preferred_time, equipment
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		 RETURNING id, created_at, updated_at`,
		p.UserID, p.Gender, p.Age, p.WeightKg, p.HeightCm, p.FitnessLevel, p.Goal,
		p.TrainingAccess, p.TrainingExperience, p.HasPain, p.PainLocations, p.PainLevel,
		p.Diagnoses, p.Contraindications, p.DaysPerWeek, p.SessionDuration, p.PreferredTime, p.Equipment,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *userRepo) GetProfileByUserID(ctx context.Context, userID int64) (*models.UserProfile, error) {
	p := &models.UserProfile{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, gender, age, weight_kg, height_cm, fitness_level, goal,
			training_access, training_experience, has_pain, pain_locations, pain_level,
			diagnoses, contraindications, days_per_week, session_duration, preferred_time, equipment,
			created_at, updated_at
		 FROM user_profiles WHERE user_id = $1`, userID,
	).Scan(&p.ID, &p.UserID, &p.Gender, &p.Age, &p.WeightKg, &p.HeightCm,
		&p.FitnessLevel, &p.Goal,
		&p.TrainingAccess, &p.TrainingExperience, &p.HasPain, &p.PainLocations, &p.PainLevel,
		&p.Diagnoses, &p.Contraindications, &p.DaysPerWeek, &p.SessionDuration, &p.PreferredTime, &p.Equipment,
		&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *userRepo) UpdateProfile(ctx context.Context, p *models.UserProfile) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE user_profiles SET
			gender=$2, age=$3, weight_kg=$4, height_cm=$5, fitness_level=$6, goal=$7,
			training_access=$8, training_experience=$9, has_pain=$10, pain_locations=$11, pain_level=$12,
			diagnoses=$13, contraindications=$14, days_per_week=$15, session_duration=$16,
			preferred_time=$17, equipment=$18, updated_at=NOW()
		 WHERE user_id=$1`,
		p.UserID, p.Gender, p.Age, p.WeightKg, p.HeightCm, p.FitnessLevel, p.Goal,
		p.TrainingAccess, p.TrainingExperience, p.HasPain, p.PainLocations, p.PainLevel,
		p.Diagnoses, p.Contraindications, p.DaysPerWeek, p.SessionDuration, p.PreferredTime, p.Equipment)
	return err
}
