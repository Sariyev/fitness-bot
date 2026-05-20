BEGIN;

DROP INDEX IF EXISTS idx_payments_category;
ALTER TABLE payments DROP COLUMN IF EXISTS category;

DROP INDEX IF EXISTS idx_user_category_access_user;
DROP TABLE IF EXISTS user_category_access;
DROP TABLE IF EXISTS category_pricing;

DROP INDEX IF EXISTS idx_meal_plans_access_tier;
DROP INDEX IF EXISTS idx_rehab_courses_access_tier;
DROP INDEX IF EXISTS idx_programs_access_tier;

ALTER TABLE meal_plans DROP COLUMN IF EXISTS access_tier;
ALTER TABLE rehab_courses DROP COLUMN IF EXISTS access_tier;
ALTER TABLE programs DROP COLUMN IF EXISTS access_tier;

COMMIT;
