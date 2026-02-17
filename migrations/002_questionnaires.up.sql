BEGIN;

CREATE TABLE questionnaires (
    id              SERIAL PRIMARY KEY,
    slug            VARCHAR(100) UNIQUE NOT NULL,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE questions (
    id                  SERIAL PRIMARY KEY,
    questionnaire_id    INTEGER NOT NULL REFERENCES questionnaires(id) ON DELETE CASCADE,
    text                TEXT NOT NULL,
    question_type       VARCHAR(30) NOT NULL,
    sort_order          INTEGER NOT NULL DEFAULT 0,
    is_required         BOOLEAN NOT NULL DEFAULT TRUE,
    metadata            JSONB DEFAULT '{}',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_questions_questionnaire ON questions(questionnaire_id);
CREATE INDEX idx_questions_sort ON questions(questionnaire_id, sort_order);

CREATE TABLE question_options (
    id              SERIAL PRIMARY KEY,
    question_id     INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    text            VARCHAR(255) NOT NULL,
    value           VARCHAR(100) NOT NULL,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_question_options_question ON question_options(question_id);

CREATE TABLE questionnaire_submissions (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    questionnaire_id    INTEGER NOT NULL REFERENCES questionnaires(id) ON DELETE CASCADE,
    completed_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_questionnaire_submissions_user ON questionnaire_submissions(user_id);
CREATE INDEX idx_questionnaire_submissions_quiz ON questionnaire_submissions(questionnaire_id);

CREATE TABLE questionnaire_answers (
    id              BIGSERIAL PRIMARY KEY,
    submission_id   BIGINT NOT NULL REFERENCES questionnaire_submissions(id) ON DELETE CASCADE,
    question_id     INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    answer_text     TEXT,
    answer_value    VARCHAR(100),
    answer_values   TEXT[],
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_questionnaire_answers_submission ON questionnaire_answers(submission_id);

COMMIT;
