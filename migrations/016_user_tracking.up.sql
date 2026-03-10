BEGIN;

CREATE TABLE user_program_enrollments (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    program_id      INTEGER NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    started_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    current_week    INTEGER NOT NULL DEFAULT 1,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_program_enrollments_user ON user_program_enrollments(user_id);
CREATE UNIQUE INDEX idx_user_program_enrollments_unique_active
    ON user_program_enrollments(user_id, program_id) WHERE is_active = TRUE;

CREATE TABLE user_rehab_progress (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id       INTEGER NOT NULL REFERENCES rehab_courses(id) ON DELETE CASCADE,
    session_id      INTEGER NOT NULL REFERENCES rehab_sessions(id) ON DELETE CASCADE,
    day_number      INTEGER NOT NULL,
    completed_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    pain_level      INTEGER DEFAULT 0 CHECK (pain_level >= 0 AND pain_level <= 10),
    comment         TEXT,
    UNIQUE(user_id, session_id)
);

CREATE INDEX idx_user_rehab_progress_user ON user_rehab_progress(user_id);
CREATE INDEX idx_user_rehab_progress_course ON user_rehab_progress(user_id, course_id);

CREATE TABLE daily_completions (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entity_type     VARCHAR(30) NOT NULL,
    entity_id       INTEGER NOT NULL,
    date            DATE NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'completed',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, entity_type, entity_id, date)
);

CREATE INDEX idx_daily_completions_user_date ON daily_completions(user_id, date);

CREATE TABLE food_log_entries (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date            DATE NOT NULL,
    meal_type       VARCHAR(20) NOT NULL,
    food_name       VARCHAR(255) NOT NULL,
    calories        INTEGER,
    protein         DECIMAL(6,1),
    fat             DECIMAL(6,1),
    carbs           DECIMAL(6,1),
    photo_url       VARCHAR(500),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_food_log_entries_user_date ON food_log_entries(user_id, date);

CREATE TABLE progress_entries (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date            DATE NOT NULL,
    weight_kg       DECIMAL(5,2),
    measurements    JSONB DEFAULT '{}',
    photo_url       VARCHAR(500),
    wellbeing       VARCHAR(30),
    pain_level      INTEGER DEFAULT 0 CHECK (pain_level >= 0 AND pain_level <= 10),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, date)
);

CREATE INDEX idx_progress_entries_user ON progress_entries(user_id);

CREATE TABLE user_achievements (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id  INTEGER NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
    earned_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, achievement_id)
);

CREATE INDEX idx_user_achievements_user ON user_achievements(user_id);

COMMIT;
