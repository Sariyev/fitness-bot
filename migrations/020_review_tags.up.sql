BEGIN;

-- Add tags column to user_scores
ALTER TABLE user_scores ADD COLUMN tags TEXT[] DEFAULT '{}';

-- Drop old constraint, migrate data, then add new constraint
ALTER TABLE user_scores DROP CONSTRAINT IF EXISTS user_scores_score_check;
UPDATE user_scores SET score = GREATEST(1, ROUND(score / 2.0)::int) WHERE score > 5;
ALTER TABLE user_scores ADD CONSTRAINT user_scores_score_check CHECK (score >= 1 AND score <= 5);

COMMIT;
