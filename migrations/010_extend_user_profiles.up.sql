BEGIN;

ALTER TABLE user_profiles
    ADD COLUMN training_access     VARCHAR(20)  DEFAULT NULL,
    ADD COLUMN training_experience VARCHAR(20)  DEFAULT NULL,
    ADD COLUMN has_pain            BOOLEAN      DEFAULT FALSE,
    ADD COLUMN pain_locations      TEXT[]        DEFAULT '{}',
    ADD COLUMN pain_level          INTEGER      DEFAULT 0 CHECK (pain_level >= 0 AND pain_level <= 10),
    ADD COLUMN diagnoses           TEXT[]        DEFAULT '{}',
    ADD COLUMN contraindications   TEXT          DEFAULT '',
    ADD COLUMN days_per_week       INTEGER      DEFAULT NULL CHECK (days_per_week IS NULL OR (days_per_week >= 2 AND days_per_week <= 7)),
    ADD COLUMN session_duration    INTEGER      DEFAULT NULL,
    ADD COLUMN preferred_time      VARCHAR(20)  DEFAULT NULL,
    ADD COLUMN equipment           TEXT[]        DEFAULT '{}';

COMMIT;
