package webapp

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"fitness-bot/internal/models"
	"fitness-bot/internal/service"
)

type contextKey string

const userContextKey contextKey = "user"

func UserFromContext(ctx context.Context) *models.User {
	u, _ := ctx.Value(userContextKey).(*models.User)
	return u
}

type TelegramInitData struct {
	QueryID      string       `json:"query_id"`
	User         TelegramUser `json:"user"`
	AuthDate     string       `json:"auth_date"`
	Hash         string       `json:"hash"`
}

type TelegramUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

func ValidateInitData(initData string, botToken string) (*TelegramInitData, error) {
	params, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("invalid init data format: %w", err)
	}

	hash := params.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("missing hash")
	}

	// Build data check string: sort all key=value pairs except hash, join with \n
	var pairs []string
	for key, values := range params {
		if key == "hash" {
			continue
		}
		for _, val := range values {
			pairs = append(pairs, key+"="+val)
		}
	}
	sort.Strings(pairs)
	dataCheckString := strings.Join(pairs, "\n")

	// secret_key = HMAC-SHA256("WebAppData", bot_token)
	secretKeyHMAC := hmac.New(sha256.New, []byte("WebAppData"))
	secretKeyHMAC.Write([]byte(botToken))
	secretKey := secretKeyHMAC.Sum(nil)

	// calculated_hash = HMAC-SHA256(secret_key, data_check_string)
	calculatedHMAC := hmac.New(sha256.New, secretKey)
	calculatedHMAC.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(calculatedHMAC.Sum(nil))

	if !hmac.Equal([]byte(calculatedHash), []byte(hash)) {
		return nil, fmt.Errorf("invalid hash")
	}

	// Parse user data
	result := &TelegramInitData{
		AuthDate: params.Get("auth_date"),
		Hash:     hash,
		QueryID:  params.Get("query_id"),
	}

	userJSON := params.Get("user")
	if userJSON != "" {
		if err := json.Unmarshal([]byte(userJSON), &result.User); err != nil {
			return nil, fmt.Errorf("invalid user data: %w", err)
		}
	}

	if result.User.ID == 0 {
		return nil, fmt.Errorf("missing user id")
	}

	return result, nil
}

func AuthMiddleware(botToken string, userSvc *service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			initData := r.Header.Get("X-Telegram-Init-Data")
			if initData == "" {
				http.Error(w, `{"error":"missing init data"}`, http.StatusUnauthorized)
				return
			}

			parsed, err := ValidateInitData(initData, botToken)
			if err != nil {
				http.Error(w, `{"error":"invalid init data"}`, http.StatusUnauthorized)
				return
			}

			user, err := userSvc.GetOrCreateUser(
				r.Context(),
				parsed.User.ID,
				parsed.User.Username,
				parsed.User.FirstName,
				parsed.User.LastName,
			)
			if err != nil {
				http.Error(w, `{"error":"user lookup failed"}`, http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
