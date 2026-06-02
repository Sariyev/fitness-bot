BEGIN;

DROP INDEX IF EXISTS idx_meals_image_media;
DROP INDEX IF EXISTS idx_meal_plans_image_media;

ALTER TABLE meals      DROP COLUMN IF EXISTS image_media_id;
ALTER TABLE meal_plans DROP COLUMN IF EXISTS image_media_id;

COMMIT;
