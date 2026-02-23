-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Admin users
CREATE TABLE admins (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP DEFAULT NOW()
);

-- Quiz versions (only one active at a time)
CREATE TABLE questionnaires (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version    SERIAL,
    is_active  BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Questions belonging to a questionnaire
CREATE TABLE questions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    questionnaire_id UUID NOT NULL REFERENCES questionnaires(id) ON DELETE CASCADE,
    text             TEXT NOT NULL,
    sort_order       INT NOT NULL
);

-- Possible answers for each question
CREATE TABLE answers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    text        TEXT NOT NULL,
    score       INT NOT NULL DEFAULT 0,
    sort_order  INT NOT NULL
);

-- Candidates (no login required)
CREATE TABLE candidates (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name       VARCHAR(100) NOT NULL,
    last_name        VARCHAR(100) NOT NULL,
    email            VARCHAR(255) UNIQUE NOT NULL,
    questionnaire_id UUID REFERENCES questionnaires(id),
    total_score      INT DEFAULT 0,
    audio_path       VARCHAR(500),
    completed        BOOLEAN DEFAULT FALSE,
    transcript       TEXT,
    ai_analysis      TEXT,
    ai_score         INT,
    analysis_status  VARCHAR(20) DEFAULT 'pending',
    created_at       TIMESTAMP DEFAULT NOW()
);

-- Candidate selected answers
CREATE TABLE candidate_answers (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    question_id  UUID NOT NULL REFERENCES questions(id),
    answer_id    UUID NOT NULL REFERENCES answers(id),
    score        INT NOT NULL
);
