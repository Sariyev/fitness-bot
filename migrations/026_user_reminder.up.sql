BEGIN;

-- Used by the daily-reminder cron in cmd/bot to dedupe per-user-per-day:
-- NULL means "never reminded"; a value < DATE_TRUNC('day', NOW()) means
-- "haven't reminded today yet".
ALTER TABLE users ADD COLUMN last_reminder_at TIMESTAMPTZ;

-- Partial index so the reminder-target query (registered + ready-to-remind)
-- doesn't scan the whole users table.
CREATE INDEX idx_users_last_reminder ON users(last_reminder_at)
    WHERE is_registered = TRUE;

COMMIT;
