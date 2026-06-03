BEGIN;

-- Irreversible: this migration drops tables and merges program names into
-- workout names. The down migration restores the schema shape but cannot
-- recover the original (program, workout) split — historical workout names
-- will keep the "<program> — <workout>" form.

CREATE TABLE IF NOT EXISTS programs (
    id             SERIAL PRIMARY KEY,
    slug           TEXT UNIQUE NOT NULL,
    name           TEXT NOT NULL,
    description    TEXT,
    goal           TEXT,
    format         TEXT,
    level          TEXT,
    duration_weeks INT,
    access_tier    TEXT,
    is_active      BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order     INT NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_program_enrollments (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL,
    program_id   INT NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    started_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    current_week INT NOT NULL DEFAULT 1,
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE workouts
    ADD COLUMN IF NOT EXISTS program_id  INT REFERENCES programs(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS week_number INT,
    ADD COLUMN IF NOT EXISTS day_number  INT;

DROP INDEX IF EXISTS idx_workouts_access_tier;
ALTER TABLE workouts DROP COLUMN IF EXISTS access_tier;

COMMIT;
