BEGIN;

-- Three-bucket access model: each content item lives in free / trial / paid.
-- Default 'paid' keeps existing content gated exactly as it was before this
-- migration; admin can demote items to 'free' or 'trial' from the UI later.

ALTER TABLE programs ADD COLUMN access_tier VARCHAR(10) NOT NULL DEFAULT 'paid'
    CHECK (access_tier IN ('free', 'trial', 'paid'));
ALTER TABLE rehab_courses ADD COLUMN access_tier VARCHAR(10) NOT NULL DEFAULT 'paid'
    CHECK (access_tier IN ('free', 'trial', 'paid'));
ALTER TABLE meal_plans ADD COLUMN access_tier VARCHAR(10) NOT NULL DEFAULT 'paid'
    CHECK (access_tier IN ('free', 'trial', 'paid'));

CREATE INDEX idx_programs_access_tier ON programs(access_tier);
CREATE INDEX idx_rehab_courses_access_tier ON rehab_courses(access_tier);
CREATE INDEX idx_meal_plans_access_tier ON meal_plans(access_tier);

-- Admin-editable price per category. One row per category, never deleted.
CREATE TABLE category_pricing (
    category   VARCHAR(20) PRIMARY KEY
        CHECK (category IN ('workouts', 'lfk', 'nutrition')),
    price_kzt  INTEGER NOT NULL CHECK (price_kzt > 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO category_pricing (category, price_kzt) VALUES
    ('workouts',  5000),
    ('lfk',       5000),
    ('nutrition', 5000);

-- Permanent per-user, per-category access granted on successful payment.
-- Composite PK enforces "one paid grant per category per user."
CREATE TABLE user_category_access (
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category    VARCHAR(20) NOT NULL
        CHECK (category IN ('workouts', 'lfk', 'nutrition')),
    payment_id  BIGINT REFERENCES payments(id) ON DELETE SET NULL,
    granted_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, category)
);

CREATE INDEX idx_user_category_access_user ON user_category_access(user_id);

-- Tag each payment with the category it unlocks. NULL = legacy payment
-- (pre-split, when users.is_paid was a single global flag); the access
-- service grandfathers those in via users.is_paid for backwards compat.
ALTER TABLE payments ADD COLUMN category VARCHAR(20)
    CHECK (category IS NULL OR category IN ('workouts', 'lfk', 'nutrition'));

CREATE INDEX idx_payments_category ON payments(category);

COMMIT;
