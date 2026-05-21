BEGIN;
DROP INDEX IF EXISTS idx_users_last_reminder;
ALTER TABLE users DROP COLUMN IF EXISTS last_reminder_at;
COMMIT;
