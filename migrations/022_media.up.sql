BEGIN;

CREATE TABLE media (
    id BIGSERIAL PRIMARY KEY,
    storage_key TEXT NOT NULL UNIQUE,
    bucket VARCHAR(20) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL,
    owner_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    reference_type VARCHAR(50),
    reference_id BIGINT,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_media_reference ON media(reference_type, reference_id);
CREATE INDEX idx_media_owner ON media(owner_user_id);
CREATE INDEX idx_media_confirmed ON media(confirmed) WHERE confirmed = TRUE;

COMMIT;
