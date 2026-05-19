package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Role identifies which binary is loading the config, for Validate().
type Role string

const (
	RoleBot     Role = "bot"
	RoleAdmin   Role = "admin"
	RoleWebApp  Role = "webapp"
)

type Config struct {
	TelegramToken string
	Debug         bool
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	AdminPort     string
	AdminAPIKey   string
	WebAppPort    string
	WebAppURL     string

	// Robokassa
	RobokassaMerchantLogin string
	RobokassaPassword1     string
	RobokassaPassword2     string
	RobokassaIsTest        bool

	// Cloudflare R2 object storage. Empty AccessKeyID disables the media
	// feature server-side (handler routes are skipped at boot).
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2BucketPrivate   string
	R2BucketPublic    string
	R2PublicURL       string
	// MediaQuotaBytes is a hard cap on the total size of confirmed media we
	// store. Defaults below R2's 10 GB free tier so an over-quota upload
	// never tips us into paid usage. Override via MEDIA_QUOTA_BYTES.
	MediaQuotaBytes int64
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		TelegramToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		Debug:         getEnv("DEBUG", "false") == "true",
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", ""),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", ""),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		AdminPort:     getEnv("ADMIN_PORT", "8080"),
		AdminAPIKey:   getEnv("ADMIN_API_KEY", ""),
		WebAppPort:    getEnv("WEBAPP_PORT", "8081"),
		WebAppURL:     getEnv("WEBAPP_URL", ""),

		RobokassaMerchantLogin: getEnv("ROBOKASSA_MERCHANT_LOGIN", ""),
		RobokassaPassword1:     getEnv("ROBOKASSA_PASSWORD1", ""),
		RobokassaPassword2:     getEnv("ROBOKASSA_PASSWORD2", ""),
		RobokassaIsTest:        getEnv("ROBOKASSA_IS_TEST", "true") == "true",

		R2AccountID:       getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2BucketPrivate:   getEnv("R2_BUCKET_PRIVATE", "fitness-bot"),
		R2BucketPublic:    getEnv("R2_BUCKET_PUBLIC", "fitness-bot-public"),
		R2PublicURL:       getEnv("R2_PUBLIC_URL", ""),
		// 9.5 GB default — half-GB headroom under R2's 10 GB free tier.
		MediaQuotaBytes: getEnvInt64("MEDIA_QUOTA_BYTES", 9_500_000_000),
	}
}

func getEnvInt64(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		if n, err := strconv.ParseInt(value, 10, 64); err == nil {
			return n
		}
	}
	return fallback
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

// Validate checks that all required env vars are set for the given role.
// Returns a multi-line error listing every missing var so the operator fixes
// them in one go instead of one-at-a-time.
func (c *Config) Validate(role Role) error {
	var missing []string
	require := func(name, value string) {
		if value == "" {
			missing = append(missing, name)
		}
	}

	require("DB_USER", c.DBUser)
	require("DB_PASSWORD", c.DBPassword)
	require("DB_NAME", c.DBName)
	require("DB_HOST", c.DBHost)

	switch role {
	case RoleBot:
		require("TELEGRAM_BOT_TOKEN", c.TelegramToken)
	case RoleAdmin:
		require("ADMIN_API_KEY", c.AdminAPIKey)
	case RoleWebApp:
		require("TELEGRAM_BOT_TOKEN", c.TelegramToken)
	}

	if len(missing) > 0 {
		return errors.New("missing required env vars: " + strings.Join(missing, ", "))
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
