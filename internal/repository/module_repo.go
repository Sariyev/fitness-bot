package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type moduleRepo struct {
	pool *pgxpool.Pool
}

func NewModuleRepo(pool *pgxpool.Pool) ModuleRepository {
	return &moduleRepo{pool: pool}
}

func (r *moduleRepo) ListActiveModules(ctx context.Context) ([]models.Module, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, slug, name, description, icon, requires_subscription, is_active, sort_order, created_at, updated_at
		 FROM modules WHERE is_active = TRUE ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []models.Module
	for rows.Next() {
		var m models.Module
		if err := rows.Scan(&m.ID, &m.Slug, &m.Name, &m.Description, &m.Icon,
			&m.RequiresSubscription, &m.IsActive, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}
	return modules, nil
}

func (r *moduleRepo) GetModuleByID(ctx context.Context, id int) (*models.Module, error) {
	m := &models.Module{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, slug, name, description, icon, requires_subscription, is_active, sort_order, created_at, updated_at
		 FROM modules WHERE id = $1`, id,
	).Scan(&m.ID, &m.Slug, &m.Name, &m.Description, &m.Icon,
		&m.RequiresSubscription, &m.IsActive, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *moduleRepo) GetModuleBySlug(ctx context.Context, slug string) (*models.Module, error) {
	m := &models.Module{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, slug, name, description, icon, requires_subscription, is_active, sort_order, created_at, updated_at
		 FROM modules WHERE slug = $1`, slug,
	).Scan(&m.ID, &m.Slug, &m.Name, &m.Description, &m.Icon,
		&m.RequiresSubscription, &m.IsActive, &m.SortOrder, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *moduleRepo) CreateModule(ctx context.Context, m *models.Module) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO modules (slug, name, description, icon, requires_subscription, is_active, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`,
		m.Slug, m.Name, m.Description, m.Icon, m.RequiresSubscription, m.IsActive, m.SortOrder,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *moduleRepo) UpdateModule(ctx context.Context, m *models.Module) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE modules SET slug=$2, name=$3, description=$4, icon=$5, requires_subscription=$6, is_active=$7, sort_order=$8, updated_at=NOW()
		 WHERE id=$1`,
		m.ID, m.Slug, m.Name, m.Description, m.Icon, m.RequiresSubscription, m.IsActive, m.SortOrder)
	return err
}

func (r *moduleRepo) ListCategoriesByModule(ctx context.Context, moduleID int) ([]models.ModuleCategory, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, module_id, slug, name, description, icon, sort_order, is_active, created_at, updated_at
		 FROM module_categories WHERE module_id = $1 AND is_active = TRUE ORDER BY sort_order`, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.ModuleCategory
	for rows.Next() {
		var c models.ModuleCategory
		if err := rows.Scan(&c.ID, &c.ModuleID, &c.Slug, &c.Name, &c.Description, &c.Icon,
			&c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

func (r *moduleRepo) GetCategoryByID(ctx context.Context, id int) (*models.ModuleCategory, error) {
	c := &models.ModuleCategory{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, module_id, slug, name, description, icon, sort_order, is_active, created_at, updated_at
		 FROM module_categories WHERE id = $1`, id,
	).Scan(&c.ID, &c.ModuleID, &c.Slug, &c.Name, &c.Description, &c.Icon,
		&c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *moduleRepo) CreateCategory(ctx context.Context, c *models.ModuleCategory) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO module_categories (module_id, slug, name, description, icon, sort_order, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`,
		c.ModuleID, c.Slug, c.Name, c.Description, c.Icon, c.SortOrder, c.IsActive,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *moduleRepo) UpdateCategory(ctx context.Context, c *models.ModuleCategory) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE module_categories SET slug=$2, name=$3, description=$4, icon=$5, sort_order=$6, is_active=$7, updated_at=NOW()
		 WHERE id=$1`,
		c.ID, c.Slug, c.Name, c.Description, c.Icon, c.SortOrder, c.IsActive)
	return err
}

func (r *moduleRepo) ListLessonsByCategory(ctx context.Context, categoryID int) ([]models.Lesson, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, category_id, slug, title, description, sort_order, is_active, created_at, updated_at
		 FROM lessons WHERE category_id = $1 AND is_active = TRUE ORDER BY sort_order`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []models.Lesson
	for rows.Next() {
		var l models.Lesson
		if err := rows.Scan(&l.ID, &l.CategoryID, &l.Slug, &l.Title, &l.Description,
			&l.SortOrder, &l.IsActive, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

func (r *moduleRepo) GetLessonByID(ctx context.Context, id int) (*models.Lesson, error) {
	l := &models.Lesson{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, category_id, slug, title, description, sort_order, is_active, created_at, updated_at
		 FROM lessons WHERE id = $1`, id,
	).Scan(&l.ID, &l.CategoryID, &l.Slug, &l.Title, &l.Description,
		&l.SortOrder, &l.IsActive, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (r *moduleRepo) CreateLesson(ctx context.Context, l *models.Lesson) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO lessons (category_id, slug, title, description, sort_order, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`,
		l.CategoryID, l.Slug, l.Title, l.Description, l.SortOrder, l.IsActive,
	).Scan(&l.ID, &l.CreatedAt, &l.UpdatedAt)
}

func (r *moduleRepo) UpdateLesson(ctx context.Context, l *models.Lesson) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE lessons SET slug=$2, title=$3, description=$4, sort_order=$5, is_active=$6, updated_at=NOW()
		 WHERE id=$1`,
		l.ID, l.Slug, l.Title, l.Description, l.SortOrder, l.IsActive)
	return err
}

func (r *moduleRepo) ListContentByLesson(ctx context.Context, lessonID int) ([]models.LessonContent, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, lesson_id, content_type, title, body, video_url, telegram_file_id, file_url, sort_order, created_at, updated_at
		 FROM lesson_contents WHERE lesson_id = $1 ORDER BY sort_order`, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []models.LessonContent
	for rows.Next() {
		var c models.LessonContent
		if err := rows.Scan(&c.ID, &c.LessonID, &c.ContentType, &c.Title, &c.Body,
			&c.VideoURL, &c.TelegramFileID, &c.FileURL, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}
	return contents, nil
}

func (r *moduleRepo) GetContentByID(ctx context.Context, id int) (*models.LessonContent, error) {
	c := &models.LessonContent{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, lesson_id, content_type, title, body, video_url, telegram_file_id, file_url, sort_order, created_at, updated_at
		 FROM lesson_contents WHERE id = $1`, id,
	).Scan(&c.ID, &c.LessonID, &c.ContentType, &c.Title, &c.Body,
		&c.VideoURL, &c.TelegramFileID, &c.FileURL, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *moduleRepo) CreateContent(ctx context.Context, c *models.LessonContent) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO lesson_contents (lesson_id, content_type, title, body, video_url, telegram_file_id, file_url, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`,
		c.LessonID, c.ContentType, c.Title, c.Body, c.VideoURL, c.TelegramFileID, c.FileURL, c.SortOrder,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *moduleRepo) UpdateContent(ctx context.Context, c *models.LessonContent) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE lesson_contents SET content_type=$2, title=$3, body=$4, video_url=$5, telegram_file_id=$6, file_url=$7, sort_order=$8, updated_at=NOW()
		 WHERE id=$1`,
		c.ID, c.ContentType, c.Title, c.Body, c.VideoURL, c.TelegramFileID, c.FileURL, c.SortOrder)
	return err
}

func (r *moduleRepo) DeleteContent(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM lesson_contents WHERE id = $1`, id)
	return err
}

func (r *moduleRepo) UpdateTelegramFileID(ctx context.Context, contentID int, fileID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE lesson_contents SET telegram_file_id = $2, updated_at = NOW() WHERE id = $1`,
		contentID, fileID)
	return err
}

func (r *moduleRepo) UpsertProgress(ctx context.Context, userID int64, lessonID int) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_lesson_progress (user_id, lesson_id, status)
		 VALUES ($1, $2, 'started')
		 ON CONFLICT (user_id, lesson_id) DO NOTHING`,
		userID, lessonID)
	return err
}

func (r *moduleRepo) CompleteLesson(ctx context.Context, userID int64, lessonID int) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_lesson_progress (user_id, lesson_id, status, completed_at)
		 VALUES ($1, $2, 'completed', NOW())
		 ON CONFLICT (user_id, lesson_id) DO UPDATE SET status = 'completed', completed_at = NOW()`,
		userID, lessonID)
	return err
}

func (r *moduleRepo) GetUserProgress(ctx context.Context, userID int64, lessonID int) (*models.UserLessonProgress, error) {
	p := &models.UserLessonProgress{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, lesson_id, status, started_at, completed_at
		 FROM user_lesson_progress WHERE user_id = $1 AND lesson_id = $2`,
		userID, lessonID,
	).Scan(&p.ID, &p.UserID, &p.LessonID, &p.Status, &p.StartedAt, &p.CompletedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *moduleRepo) GetUserCategoryProgress(ctx context.Context, userID int64, categoryID int) (int, int, error) {
	var completed, total int
	err := r.pool.QueryRow(ctx,
		`SELECT
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE ulp.status = 'completed') AS completed
		 FROM lessons l
		 LEFT JOIN user_lesson_progress ulp ON ulp.lesson_id = l.id AND ulp.user_id = $1
		 WHERE l.category_id = $2 AND l.is_active = TRUE`,
		userID, categoryID,
	).Scan(&total, &completed)
	return completed, total, err
}

func (r *moduleRepo) SelectCategory(ctx context.Context, userID int64, categoryID int) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_module_selections (user_id, category_id)
		 VALUES ($1, $2)
		 ON CONFLICT (user_id, category_id) DO NOTHING`,
		userID, categoryID)
	return err
}

func (r *moduleRepo) GetUserSelections(ctx context.Context, userID int64) ([]models.UserModuleSelection, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, category_id, selected_at
		 FROM user_module_selections WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sels []models.UserModuleSelection
	for rows.Next() {
		var s models.UserModuleSelection
		if err := rows.Scan(&s.ID, &s.UserID, &s.CategoryID, &s.SelectedAt); err != nil {
			return nil, err
		}
		sels = append(sels, s)
	}
	return sels, nil
}
