BEGIN;

CREATE TABLE exercises (
    id                    SERIAL PRIMARY KEY,
    name                  VARCHAR(255) NOT NULL,
    technique             TEXT,
    common_mistakes       TEXT,
    easier_modification   TEXT,
    harder_modification   TEXT,
    rest_seconds          INTEGER DEFAULT 60,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
