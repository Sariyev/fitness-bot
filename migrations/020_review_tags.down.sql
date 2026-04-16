BEGIN;

ALTER TABLE user_scores DROP CONSTRAINT IF EXISTS user_scores_score_check;
ALTER TABLE user_scores ADD CONSTRAINT user_scores_score_check CHECK (score >= 1 AND score <= 10);
ALTER TABLE user_scores DROP COLUMN IF EXISTS tags;

COMMIT;
