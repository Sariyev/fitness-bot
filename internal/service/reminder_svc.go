package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"fitness-bot/internal/repository"
)

// ReminderService sends a daily "time to train" message to each registered
// user at their preferred_time bucket, with a deep-link button into the
// Mini App. It runs as a goroutine in the bot binary.
//
// Dedup: each user is reminded at most once per UTC day (users.last_reminder_at).
// Throttle: small sleep between sends to stay under Telegram's ~30 msg/sec
// global limit (which only matters at scale, but cheap to do).
type ReminderService struct {
	userRepo  repository.UserRepository
	botToken  string
	webAppURL string
	tzOffset  time.Duration
	tickEvery time.Duration
}

func NewReminderService(userRepo repository.UserRepository, botToken, webAppURL string, tzOffsetHours int) *ReminderService {
	return &ReminderService{
		userRepo:  userRepo,
		botToken:  botToken,
		webAppURL: webAppURL,
		tzOffset:  time.Duration(tzOffsetHours) * time.Hour,
		tickEvery: 15 * time.Minute,
	}
}

// Run blocks until ctx is cancelled. Tick once on startup and then on the
// configured interval.
func (s *ReminderService) Run(ctx context.Context) {
	if s.webAppURL == "" {
		log.Println("[REMINDER] disabled — WEBAPP_URL not set")
		return
	}
	if s.botToken == "" {
		log.Println("[REMINDER] disabled — TELEGRAM_BOT_TOKEN not set")
		return
	}

	log.Printf("[REMINDER] started — tick=%s, tz_offset=%s", s.tickEvery, s.tzOffset)

	s.tick(ctx)

	t := time.NewTicker(s.tickEvery)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("[REMINDER] stopped")
			return
		case <-t.C:
			s.tick(ctx)
		}
	}
}

// tick decides whether the current wall-clock hour corresponds to any
// preferred_time bucket, and if so, sends reminders to matching users.
func (s *ReminderService) tick(ctx context.Context) {
	now := time.Now().UTC().Add(s.tzOffset)
	bucket := bucketForLocalHour(now.Hour())
	if bucket == "" {
		return
	}

	// Morning slot also reaches users with NULL/'any' preferred_time
	// (default reminder time for users who didn't pick one).
	buckets := []string{bucket}
	if bucket == "morning" {
		buckets = append(buckets, "any")
	}

	var targets []repository.ReminderTarget
	for _, b := range buckets {
		ts, err := s.userRepo.ListReminderTargets(ctx, b)
		if err != nil {
			log.Printf("[REMINDER] list error bucket=%s: %v", b, err)
			continue
		}
		targets = append(targets, ts...)
	}
	if len(targets) == 0 {
		return
	}

	log.Printf("[REMINDER] bucket=%s local_hour=%d targets=%d", bucket, now.Hour(), len(targets))
	for _, t := range targets {
		if ctx.Err() != nil {
			return
		}
		if err := s.send(t); err != nil {
			log.Printf("[REMINDER] send err telegram_id=%d: %v", t.TelegramID, err)
			continue
		}
		if err := s.userRepo.MarkReminderSent(ctx, t.UserID); err != nil {
			log.Printf("[REMINDER] mark err user_id=%d: %v", t.UserID, err)
		}
		time.Sleep(50 * time.Millisecond) // ~20 msgs/sec — well below Telegram's 30/sec
	}
}

// bucketForLocalHour maps wall-clock hour → preferred_time bucket.
// Returns "" when the hour doesn't correspond to any reminder window.
// Buckets:
//   - morning   = 8:00 local
//   - afternoon = 13:00 local
//   - evening   = 19:00 local
//   - any       = 8:00 local (same as morning — default reminder time)
//
// Targets are queried with exact bucket-match against the user's
// preferred_time (NULL → 'any'), so 'any' users only get pinged at the
// morning slot.
func bucketForLocalHour(hour int) string {
	switch hour {
	case 8:
		// 8 AM serves both 'morning' and 'any' (NULL preferred_time).
		// Caller will run both queries.
		return "morning"
	case 13:
		return "afternoon"
	case 19:
		return "evening"
	}
	return ""
}

// send issues the reminder via Telegram's sendMessage API with an inline
// WebApp button. Raw HTTP because tgbotapi v5.5.1 doesn't have native WebApp
// button support (same pattern as handleWebApp in internal/handler/bot/router.go).
func (s *ReminderService) send(t repository.ReminderTarget) error {
	greeting := "Привет"
	if t.FirstName != "" {
		greeting = "Привет, " + t.FirstName
	}
	text := greeting + "! 💪\nДавно не виделись — время вернуться к тренировкам:"

	payload := map[string]interface{}{
		"chat_id": t.TelegramID,
		"text":    text,
		"reply_markup": map[string]interface{}{
			"inline_keyboard": [][]map[string]interface{}{
				{
					{
						"text":    "🏋️ Открыть приложение",
						"web_app": map[string]string{"url": s.webAppURL},
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram api %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
