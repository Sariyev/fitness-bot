BEGIN;

CREATE TABLE programs (
    id              SERIAL PRIMARY KEY,
    slug            VARCHAR(100) UNIQUE NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    goal            VARCHAR(50) NOT NULL,
    format          VARCHAR(20) NOT NULL,
    level           VARCHAR(30) NOT NULL,
    duration_weeks  INTEGER NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE workouts (
    id              SERIAL PRIMARY KEY,
    program_id      INTEGER REFERENCES programs(id) ON DELETE SET NULL,
    slug            VARCHAR(100) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    goal            VARCHAR(50),
    format          VARCHAR(20),
    level           VARCHAR(30),
    duration_minutes INTEGER,
    equipment       TEXT[] DEFAULT '{}',
    expected_result TEXT,
    video_url       VARCHAR(500),
    sort_order      INTEGER NOT NULL DEFAULT 0,
    week_number     INTEGER DEFAULT NULL,
    day_number      INTEGER DEFAULT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workouts_program ON workouts(program_id);
CREATE INDEX idx_workouts_goal_format ON workouts(goal, format);

CREATE TABLE workout_exercises (
    id              SERIAL PRIMARY KEY,
    workout_id      INTEGER NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id     INTEGER NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    sets            INTEGER,
    reps            VARCHAR(50),
    duration_seconds INTEGER,
    sort_order      INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_workout_exercises_workout ON workout_exercises(workout_id);

COMMIT;
