BEGIN;

CREATE TABLE rehab_courses (
    id              SERIAL PRIMARY KEY,
    slug            VARCHAR(100) UNIQUE NOT NULL,
    category        VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    warnings        TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rehab_courses_category ON rehab_courses(category);

CREATE TABLE rehab_sessions (
    id              SERIAL PRIMARY KEY,
    course_id       INTEGER NOT NULL REFERENCES rehab_courses(id) ON DELETE CASCADE,
    day_number      INTEGER NOT NULL CHECK (day_number >= 1 AND day_number <= 14),
    stage           INTEGER NOT NULL CHECK (stage >= 1 AND stage <= 3),
    video_url       VARCHAR(500),
    duration_minutes INTEGER,
    description     TEXT,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(course_id, day_number)
);

CREATE INDEX idx_rehab_sessions_course ON rehab_sessions(course_id);

COMMIT;
