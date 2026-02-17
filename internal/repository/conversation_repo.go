package repository

import (
	"context"
	"fitness-bot/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type conversationRepo struct {
	pool *pgxpool.Pool
}

func NewConversationRepo(pool *pgxpool.Pool) ConversationRepository {
	return &conversationRepo{pool: pool}
}

func (r *conversationRepo) GetState(ctx context.Context, telegramID int64) (*models.ConversationState, error) {
	s := &models.ConversationState{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, telegram_id, state, data, expires_at, created_at, updated_at
		 FROM conversation_states
		 WHERE telegram_id = $1 AND expires_at > NOW()`, telegramID,
	).Scan(&s.ID, &s.TelegramID, &s.State, &s.Data, &s.ExpiresAt, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *conversationRepo) UpsertState(ctx context.Context, state *models.ConversationState) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO conversation_states (telegram_id, state, data, expires_at)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (telegram_id) DO UPDATE SET
		   state = EXCLUDED.state,
		   data = EXCLUDED.data,
		   expires_at = EXCLUDED.expires_at,
		   updated_at = NOW()
		 RETURNING id, created_at, updated_at`,
		state.TelegramID, state.State, state.Data, state.ExpiresAt,
	).Scan(&state.ID, &state.CreatedAt, &state.UpdatedAt)
}

func (r *conversationRepo) ClearState(ctx context.Context, telegramID int64) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM conversation_states WHERE telegram_id = $1`, telegramID)
	return err
}

func (r *conversationRepo) CleanupExpired(ctx context.Context) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM conversation_states WHERE expires_at < NOW()`)
	return err
}
