BEGIN;

-- Optional link from a workout/rehab session to a media row (R2-hosted video).
-- ON DELETE SET NULL so deleting the media doesn't cascade-delete the workout.
-- Existing video_url TEXT stays as the fallback for external (YouTube) URLs.

ALTER TABLE workouts
    ADD COLUMN video_media_id BIGINT REFERENCES media(id) ON DELETE SET NULL;

ALTER TABLE rehab_sessions
    ADD COLUMN video_media_id BIGINT REFERENCES media(id) ON DELETE SET NULL;

-- Partial indexes — most rows will be NULL until admin uploads videos.
CREATE INDEX idx_workouts_video_media
    ON workouts(video_media_id) WHERE video_media_id IS NOT NULL;
CREATE INDEX idx_rehab_sessions_video_media
    ON rehab_sessions(video_media_id) WHERE video_media_id IS NOT NULL;

COMMIT;
