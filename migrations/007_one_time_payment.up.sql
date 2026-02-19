BEGIN;

-- Add is_paid flag to users
ALTER TABLE users ADD COLUMN is_paid BOOLEAN NOT NULL DEFAULT FALSE;

-- Drop subscription-related tables (payments references them, so drop in order)
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS subscription_plans;

-- Recreate payments table without subscription/plan references
CREATE TABLE payments (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount_kzt      INTEGER NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    provider        VARCHAR(50) NOT NULL DEFAULT 'dummy',
    provider_tx_id  VARCHAR(255),
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_user ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);

COMMIT;
