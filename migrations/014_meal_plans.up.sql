BEGIN;

CREATE TABLE meal_plans (
    id              SERIAL PRIMARY KEY,
    slug            VARCHAR(100) UNIQUE NOT NULL,
    name            VARCHAR(255) NOT NULL,
    goal            VARCHAR(30) NOT NULL,
    day_number      INTEGER NOT NULL,
    calories        INTEGER,
    protein         DECIMAL(6,1),
    fat             DECIMAL(6,1),
    carbs           DECIMAL(6,1),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_meal_plans_goal ON meal_plans(goal);

CREATE TABLE meals (
    id              SERIAL PRIMARY KEY,
    meal_plan_id    INTEGER NOT NULL REFERENCES meal_plans(id) ON DELETE CASCADE,
    meal_type       VARCHAR(20) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    recipe          TEXT,
    calories        INTEGER,
    protein         DECIMAL(6,1),
    fat             DECIMAL(6,1),
    carbs           DECIMAL(6,1),
    alternatives    TEXT,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_meals_plan ON meals(meal_plan_id);
CREATE INDEX idx_meals_type ON meals(meal_type);

COMMIT;
