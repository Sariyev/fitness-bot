BEGIN;

DROP INDEX IF EXISTS idx_rehab_sessions_video_media;
DROP INDEX IF EXISTS idx_workouts_video_media;

ALTER TABLE rehab_sessions DROP COLUMN IF EXISTS video_media_id;
ALTER TABLE workouts DROP COLUMN IF EXISTS video_media_id;

COMMIT;
