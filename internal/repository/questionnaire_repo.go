package repository

import (
	"context"
	"encoding/json"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type questionnaireRepo struct {
	pool *pgxpool.Pool
}

func NewQuestionnaireRepo(pool *pgxpool.Pool) QuestionnaireRepository {
	return &questionnaireRepo{pool: pool}
}

func (r *questionnaireRepo) GetBySlug(ctx context.Context, slug string) (*models.Questionnaire, error) {
	q := &models.Questionnaire{}
	var metaBytes []byte
	_ = metaBytes
	err := r.pool.QueryRow(ctx,
		`SELECT id, slug, title, description, is_active, sort_order, created_at, updated_at
		 FROM questionnaires WHERE slug = $1`, slug,
	).Scan(&q.ID, &q.Slug, &q.Title, &q.Description, &q.IsActive, &q.SortOrder, &q.CreatedAt, &q.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (r *questionnaireRepo) GetByID(ctx context.Context, id int) (*models.Questionnaire, error) {
	q := &models.Questionnaire{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, slug, title, description, is_active, sort_order, created_at, updated_at
		 FROM questionnaires WHERE id = $1`, id,
	).Scan(&q.ID, &q.Slug, &q.Title, &q.Description, &q.IsActive, &q.SortOrder, &q.CreatedAt, &q.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (r *questionnaireRepo) ListActive(ctx context.Context) ([]models.Questionnaire, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, slug, title, description, is_active, sort_order, created_at, updated_at
		 FROM questionnaires WHERE is_active = TRUE ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Questionnaire
	for rows.Next() {
		var q models.Questionnaire
		if err := rows.Scan(&q.ID, &q.Slug, &q.Title, &q.Description, &q.IsActive, &q.SortOrder, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, q)
	}
	return result, nil
}

func (r *questionnaireRepo) Create(ctx context.Context, q *models.Questionnaire) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO questionnaires (slug, title, description, is_active, sort_order)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`,
		q.Slug, q.Title, q.Description, q.IsActive, q.SortOrder,
	).Scan(&q.ID, &q.CreatedAt, &q.UpdatedAt)
}

func (r *questionnaireRepo) Update(ctx context.Context, q *models.Questionnaire) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE questionnaires SET slug=$2, title=$3, description=$4, is_active=$5, sort_order=$6, updated_at=NOW()
		 WHERE id=$1`,
		q.ID, q.Slug, q.Title, q.Description, q.IsActive, q.SortOrder)
	return err
}

func (r *questionnaireRepo) GetQuestionsByQuestionnaireID(ctx context.Context, qID int) ([]models.Question, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, questionnaire_id, text, question_type, sort_order, is_required, metadata, created_at
		 FROM questions WHERE questionnaire_id = $1 ORDER BY sort_order`, qID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		var metaBytes []byte
		if err := rows.Scan(&q.ID, &q.QuestionnaireID, &q.Text, &q.QuestionType,
			&q.SortOrder, &q.IsRequired, &metaBytes, &q.CreatedAt); err != nil {
			return nil, err
		}
		if len(metaBytes) > 0 {
			json.Unmarshal(metaBytes, &q.Metadata)
		}
		questions = append(questions, q)
	}

	for i := range questions {
		options, err := r.GetOptionsByQuestionID(ctx, questions[i].ID)
		if err != nil {
			return nil, err
		}
		questions[i].Options = options
	}

	return questions, nil
}

func (r *questionnaireRepo) GetQuestionByID(ctx context.Context, id int) (*models.Question, error) {
	q := &models.Question{}
	var metaBytes []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, questionnaire_id, text, question_type, sort_order, is_required, metadata, created_at
		 FROM questions WHERE id = $1`, id,
	).Scan(&q.ID, &q.QuestionnaireID, &q.Text, &q.QuestionType,
		&q.SortOrder, &q.IsRequired, &metaBytes, &q.CreatedAt)
	if err != nil {
		return nil, err
	}
	if len(metaBytes) > 0 {
		json.Unmarshal(metaBytes, &q.Metadata)
	}
	options, err := r.GetOptionsByQuestionID(ctx, q.ID)
	if err != nil {
		return nil, err
	}
	q.Options = options
	return q, nil
}

func (r *questionnaireRepo) CreateQuestion(ctx context.Context, q *models.Question) error {
	metaBytes, _ := json.Marshal(q.Metadata)
	return r.pool.QueryRow(ctx,
		`INSERT INTO questions (questionnaire_id, text, question_type, sort_order, is_required, metadata)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		q.QuestionnaireID, q.Text, q.QuestionType, q.SortOrder, q.IsRequired, metaBytes,
	).Scan(&q.ID, &q.CreatedAt)
}

func (r *questionnaireRepo) UpdateQuestion(ctx context.Context, q *models.Question) error {
	metaBytes, _ := json.Marshal(q.Metadata)
	_, err := r.pool.Exec(ctx,
		`UPDATE questions SET text=$2, question_type=$3, sort_order=$4, is_required=$5, metadata=$6
		 WHERE id=$1`,
		q.ID, q.Text, q.QuestionType, q.SortOrder, q.IsRequired, metaBytes)
	return err
}

func (r *questionnaireRepo) DeleteQuestion(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM questions WHERE id = $1`, id)
	return err
}

func (r *questionnaireRepo) GetOptionsByQuestionID(ctx context.Context, questionID int) ([]models.QuestionOption, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, question_id, text, value, sort_order, created_at
		 FROM question_options WHERE question_id = $1 ORDER BY sort_order`, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []models.QuestionOption
	for rows.Next() {
		var o models.QuestionOption
		if err := rows.Scan(&o.ID, &o.QuestionID, &o.Text, &o.Value, &o.SortOrder, &o.CreatedAt); err != nil {
			return nil, err
		}
		options = append(options, o)
	}
	return options, nil
}

func (r *questionnaireRepo) CreateOption(ctx context.Context, o *models.QuestionOption) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO question_options (question_id, text, value, sort_order)
		 VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		o.QuestionID, o.Text, o.Value, o.SortOrder,
	).Scan(&o.ID, &o.CreatedAt)
}

func (r *questionnaireRepo) DeleteOption(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM question_options WHERE id = $1`, id)
	return err
}

func (r *questionnaireRepo) CreateSubmission(ctx context.Context, s *models.QuestionnaireSubmission) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO questionnaire_submissions (user_id, questionnaire_id)
		 VALUES ($1, $2) RETURNING id, created_at`,
		s.UserID, s.QuestionnaireID,
	).Scan(&s.ID, &s.CreatedAt)
}

func (r *questionnaireRepo) CompleteSubmission(ctx context.Context, submissionID int64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE questionnaire_submissions SET completed_at = NOW() WHERE id = $1`, submissionID)
	return err
}

func (r *questionnaireRepo) SaveAnswer(ctx context.Context, a *models.QuestionnaireAnswer) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO questionnaire_answers (submission_id, question_id, answer_text, answer_value, answer_values)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`,
		a.SubmissionID, a.QuestionID, a.AnswerText, a.AnswerValue, a.AnswerValues,
	).Scan(&a.ID, &a.CreatedAt)
}

func (r *questionnaireRepo) GetSubmissionAnswers(ctx context.Context, submissionID int64) ([]models.QuestionnaireAnswer, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, submission_id, question_id, answer_text, answer_value, answer_values, created_at
		 FROM questionnaire_answers WHERE submission_id = $1 ORDER BY id`, submissionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []models.QuestionnaireAnswer
	for rows.Next() {
		var a models.QuestionnaireAnswer
		if err := rows.Scan(&a.ID, &a.SubmissionID, &a.QuestionID,
			&a.AnswerText, &a.AnswerValue, &a.AnswerValues, &a.CreatedAt); err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}
	return answers, nil
}

func (r *questionnaireRepo) GetUserSubmissions(ctx context.Context, userID int64, questionnaireID int) ([]models.QuestionnaireSubmission, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, questionnaire_id, completed_at, created_at
		 FROM questionnaire_submissions
		 WHERE user_id = $1 AND questionnaire_id = $2 ORDER BY created_at DESC`, userID, questionnaireID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.QuestionnaireSubmission
	for rows.Next() {
		var s models.QuestionnaireSubmission
		if err := rows.Scan(&s.ID, &s.UserID, &s.QuestionnaireID, &s.CompletedAt, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}

// Ensure pgx is used (for error handling in callers)
var _ = pgx.ErrNoRows
