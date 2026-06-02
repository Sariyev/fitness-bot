BEGIN;

-- Optional thumbnail/photo for a meal plan or an individual meal.
-- ON DELETE SET NULL mirrors the video_media_id pattern in 025 — deleting
-- the underlying media row doesn't cascade and nuke the meal.

ALTER TABLE meal_plans
    ADD COLUMN image_media_id BIGINT REFERENCES media(id) ON DELETE SET NULL;

ALTER TABLE meals
    ADD COLUMN image_media_id BIGINT REFERENCES media(id) ON DELETE SET NULL;

CREATE INDEX idx_meal_plans_image_media
    ON meal_plans(image_media_id) WHERE image_media_id IS NOT NULL;
CREATE INDEX idx_meals_image_media
    ON meals(image_media_id) WHERE image_media_id IS NOT NULL;

COMMIT;
