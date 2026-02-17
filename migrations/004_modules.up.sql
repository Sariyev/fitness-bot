BEGIN;

CREATE TABLE modules (
    id                      SERIAL PRIMARY KEY,
    slug                    VARCHAR(100) UNIQUE NOT NULL,
    name                    VARCHAR(255) NOT NULL,
    description             TEXT,
    icon                    VARCHAR(10),
    requires_subscription   BOOLEAN NOT NULL DEFAULT TRUE,
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order              INTEGER NOT NULL DEFAULT 0,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE module_categories (
    id              SERIAL PRIMARY KEY,
    module_id       INTEGER NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    slug            VARCHAR(100) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    icon            VARCHAR(10),
    sort_order      INTEGER NOT NULL DEFAULT 0,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(module_id, slug)
);

CREATE INDEX idx_module_categories_module ON module_categories(module_id);

CREATE TABLE lessons (
    id              SERIAL PRIMARY KEY,
    category_id     INTEGER NOT NULL REFERENCES module_categories(id) ON DELETE CASCADE,
    slug            VARCHAR(100) NOT NULL,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(category_id, slug)
);

CREATE INDEX idx_lessons_category ON lessons(category_id);

CREATE TABLE lesson_contents (
    id                  SERIAL PRIMARY KEY,
    lesson_id           INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    content_type        VARCHAR(30) NOT NULL,
    title               VARCHAR(255),
    body                TEXT,
    video_url           VARCHAR(500),
    telegram_file_id    VARCHAR(255),
    file_url            VARCHAR(500),
    sort_order          INTEGER NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_lesson_contents_lesson ON lesson_contents(lesson_id);

CREATE TABLE user_lesson_progress (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id       INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    status          VARCHAR(20) NOT NULL DEFAULT 'started',
    started_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at    TIMESTAMPTZ,
    UNIQUE(user_id, lesson_id)
);

CREATE INDEX idx_user_lesson_progress_user ON user_lesson_progress(user_id);
CREATE INDEX idx_user_lesson_progress_lesson ON user_lesson_progress(lesson_id);

CREATE TABLE user_module_selections (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id     INTEGER NOT NULL REFERENCES module_categories(id) ON DELETE CASCADE,
    selected_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, category_id)
);

CREATE INDEX idx_user_module_selections_user ON user_module_selections(user_id);

COMMIT;
