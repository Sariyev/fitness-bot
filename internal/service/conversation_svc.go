package service

import (
	"context"
	"encoding/json"
	"fitness-bot/internal/models"
	"fitness-bot/internal/repository"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

type ConversationService struct {
	repo repository.ConversationRepository
}

func NewConversationService(repo repository.ConversationRepository) *ConversationService {
	return &ConversationService{repo: repo}
}

func (s *ConversationService) GetState(ctx context.Context, telegramID int64) (*models.ConversationState, error) {
	state, err := s.repo.GetState(ctx, telegramID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return state, nil
}

func (s *ConversationService) SetState(ctx context.Context, telegramID int64, state string, data any, ttl time.Duration) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	cs := &models.ConversationState{
		TelegramID: telegramID,
		State:      state,
		Data:       dataBytes,
		ExpiresAt:  time.Now().Add(ttl),
	}
	return s.repo.UpsertState(ctx, cs)
}

func (s *ConversationService) ClearState(ctx context.Context, telegramID int64) error {
	return s.repo.ClearState(ctx, telegramID)
}

func (s *ConversationService) GetData(state *models.ConversationState, dest any) error {
	if state == nil || len(state.Data) == 0 {
		return nil
	}
	return json.Unmarshal(state.Data, dest)
}

func IsFlowActive(state string, prefix string) bool {
	return strings.HasPrefix(state, prefix)
}
