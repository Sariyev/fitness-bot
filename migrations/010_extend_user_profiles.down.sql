BEGIN;

ALTER TABLE user_profiles
    DROP COLUMN IF EXISTS training_access,
    DROP COLUMN IF EXISTS training_experience,
    DROP COLUMN IF EXISTS has_pain,
    DROP COLUMN IF EXISTS pain_locations,
    DROP COLUMN IF EXISTS pain_level,
    DROP COLUMN IF EXISTS diagnoses,
    DROP COLUMN IF EXISTS contraindications,
    DROP COLUMN IF EXISTS days_per_week,
    DROP COLUMN IF EXISTS session_duration,
    DROP COLUMN IF EXISTS preferred_time,
    DROP COLUMN IF EXISTS equipment;

COMMIT;
