BEGIN;

CREATE TABLE user_scores (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    score_type      VARCHAR(50) NOT NULL,
    reference_type  VARCHAR(50) NOT NULL,
    reference_id    INTEGER NOT NULL,
    score           INTEGER NOT NULL CHECK (score >= 1 AND score <= 10),
    comment         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_scores_user ON user_scores(user_id);
CREATE INDEX idx_user_scores_reference ON user_scores(reference_type, reference_id);

COMMIT;
