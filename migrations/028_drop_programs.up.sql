BEGIN;

-- Flatten the Programs → Workouts hierarchy to a single Workouts layer.
-- We copy the parent program's metadata down onto each child workout so we
-- don't lose access tier / goal / format / level, and we bake the
-- (week_number, day_number) ordering into sort_order.

-- 1) Workouts need their own access_tier (only programs/rehab_courses/meal_plans
--    got it in migration 024). Default to 'paid' so any straggler that isn't
--    backfilled from a program still has a sane lock state.
ALTER TABLE workouts ADD COLUMN access_tier VARCHAR(10) NOT NULL DEFAULT 'paid'
    CHECK (access_tier IN ('free', 'trial', 'paid'));
CREATE INDEX idx_workouts_access_tier ON workouts(access_tier);

-- 2) Backfill workouts that belong to a program. Prepend the program name
--    so we don't lose the curriculum context, and bake (week, day) into
--    sort_order so workouts stay in the right order in the flat list.
UPDATE workouts w
SET
    name        = CASE
                    WHEN p.name IS NOT NULL AND p.name <> ''
                    THEN p.name || ' — ' || w.name
                    ELSE w.name
                  END,
    goal        = CASE WHEN COALESCE(w.goal, '')   = '' THEN COALESCE(p.goal, '')   ELSE w.goal   END,
    format      = CASE WHEN COALESCE(w.format, '') = '' THEN COALESCE(p.format, '') ELSE w.format END,
    level       = CASE WHEN COALESCE(w.level, '')  = '' THEN COALESCE(p.level, '')  ELSE w.level  END,
    access_tier = p.access_tier,
    sort_order  = COALESCE(w.week_number, 1) * 100 + COALESCE(w.day_number, 1)
FROM programs p
WHERE w.program_id = p.id;

-- 3) Drop the FK columns we no longer need.
DROP INDEX IF EXISTS idx_workouts_video_media;
ALTER TABLE workouts
    DROP COLUMN IF EXISTS program_id,
    DROP COLUMN IF EXISTS week_number,
    DROP COLUMN IF EXISTS day_number;
CREATE INDEX idx_workouts_video_media
    ON workouts(video_media_id) WHERE video_media_id IS NOT NULL;

-- 4) Drop dependent tables.
DROP TABLE IF EXISTS user_program_enrollments;
DROP TABLE IF EXISTS programs;

COMMIT;
