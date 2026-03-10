BEGIN;

CREATE TABLE achievements (
    id              SERIAL PRIMARY KEY,
    slug            VARCHAR(100) UNIQUE NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    icon            VARCHAR(10),
    criteria        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed initial achievements
INSERT INTO achievements (slug, name, description, icon, criteria) VALUES
    ('first_workout', 'First Workout', 'Complete your first workout', '🏋️', '{"type": "workout_count", "target": 1}'),
    ('10_workouts', '10 Workouts', 'Complete 10 workouts', '💪', '{"type": "workout_count", "target": 10}'),
    ('50_workouts', '50 Workouts', 'Complete 50 workouts', '🔥', '{"type": "workout_count", "target": 50}'),
    ('week_streak', 'Week Streak', 'Train 7 days in a row', '📅', '{"type": "streak_days", "target": 7}'),
    ('month_streak', 'Month Streak', 'Train 30 days in a row', '🏆', '{"type": "streak_days", "target": 30}'),
    ('first_rehab', 'First Rehab', 'Complete a rehabilitation course', '❤️‍🩹', '{"type": "rehab_course_complete", "target": 1}'),
    ('pain_reducer', 'Pain Reducer', 'Reduce pain level by 3+ points during a course', '📉', '{"type": "pain_reduction", "target": 3}'),
    ('food_tracker', 'Food Tracker', 'Log food for 7 days straight', '🥗', '{"type": "food_log_days", "target": 7}'),
    ('weight_tracker', 'Weight Tracker', 'Log weight 10 times', '⚖️', '{"type": "weight_log_count", "target": 10}'),
    ('early_bird', 'Early Bird', 'Complete 5 morning workouts', '🌅', '{"type": "morning_workouts", "target": 5}');

COMMIT;
