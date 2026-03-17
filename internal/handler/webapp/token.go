package webapp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const tokenTTL = 24 * time.Hour

type tokenPayload struct {
	TelegramID int64 `json:"tid"`
	ExpiresAt  int64 `json:"exp"`
}

func GenerateToken(telegramID int64, botToken string) (string, error) {
	payload := tokenPayload{
		TelegramID: telegramID,
		ExpiresAt:  time.Now().Add(tokenTTL).Unix(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal token payload: %w", err)
	}

	encoded := base64.RawURLEncoding.EncodeToString(data)
	sig := signPayload(encoded, botToken)

	return encoded + "." + sig, nil
}

func ValidateToken(token string, botToken string) (int64, error) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid token format")
	}

	encoded, sig := parts[0], parts[1]

	expectedSig := signPayload(encoded, botToken)
	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return 0, fmt.Errorf("invalid token signature")
	}

	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return 0, fmt.Errorf("decode token payload: %w", err)
	}

	var payload tokenPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return 0, fmt.Errorf("unmarshal token payload: %w", err)
	}

	if time.Now().Unix() > payload.ExpiresAt {
		return 0, fmt.Errorf("token expired")
	}

	if payload.TelegramID == 0 {
		return 0, fmt.Errorf("missing telegram id in token")
	}

	return payload.TelegramID, nil
}

func signPayload(payload string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
