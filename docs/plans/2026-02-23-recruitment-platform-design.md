# Recruitment Platform - Design Document

## Overview

Platform for personnel selection. Candidates register, complete a scored quiz, and record an audio presentation. Admins manage quiz questions and review candidates ranked by score.

## Architecture

```
project/
  backend/                        # Go API server (Air hot reload)
  frontend/                       # Vue 3 SPA (Vite dev server)
  docker-entrypoint-initdb.d/     # Postgres init scripts
    01_schema.sql                 # Table creation
    02_seed.sql                   # Example data
  docker-compose.yml              # Go + Vue + Postgres
```

Development only. No production build, no Nginx.

- **Backend**: Go + GORM + Postgres, serves REST API and handles audio upload/download
- **Frontend**: Vue 3 + Composition API + Pinia Colada + DaisyUI + Tailwind
- **Database**: Postgres in Docker
- **Docker Compose**: 3 services, hot reload on both backend and frontend

## Database Schema

```sql
-- Admin users
admins (
  id            UUID PK DEFAULT gen_random_uuid(),
  email         VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at    TIMESTAMP DEFAULT NOW()
)

-- Quiz versions (only one active at a time)
questionnaires (
  id         UUID PK DEFAULT gen_random_uuid(),
  version    SERIAL,
  is_active  BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW()
)

-- Questions belonging to a questionnaire
questions (
  id               UUID PK DEFAULT gen_random_uuid(),
  questionnaire_id UUID FK -> questionnaires ON DELETE CASCADE,
  text             TEXT NOT NULL,
  sort_order       INT NOT NULL
)

-- Possible answers for each question
answers (
  id          UUID PK DEFAULT gen_random_uuid(),
  question_id UUID FK -> questions ON DELETE CASCADE,
  text        TEXT NOT NULL,
  score       INT NOT NULL DEFAULT 0,
  sort_order  INT NOT NULL
)

-- Candidates (no login required)
candidates (
  id               UUID PK DEFAULT gen_random_uuid(),
  first_name       VARCHAR(100) NOT NULL,
  last_name        VARCHAR(100) NOT NULL,
  email            VARCHAR(255) UNIQUE NOT NULL,
  questionnaire_id UUID FK -> questionnaires,
  total_score      INT DEFAULT 0,
  audio_path       VARCHAR(500),
  completed        BOOLEAN DEFAULT FALSE,
  created_at       TIMESTAMP DEFAULT NOW()
)

-- Candidate selected answers
candidate_answers (
  id           UUID PK DEFAULT gen_random_uuid(),
  candidate_id UUID FK -> candidates ON DELETE CASCADE,
  question_id  UUID FK -> questions,
  answer_id    UUID FK -> answers,
  score        INT NOT NULL
)
```

### Questionnaire versioning

- Admin creates a new questionnaire version, the old one gets `is_active = false`
- New candidates always get the questionnaire with `is_active = true`
- Old candidates keep their link to the version they completed
- Scores are copied to `candidate_answers.score` at submission time (immutable)

## API Endpoints

### Public (Candidates)

```
POST   /api/candidates              # Step 1: Register (first_name, last_name, email)
                                     # Returns candidate_id

GET    /api/questionnaire            # Get active questionnaire with questions and answers

POST   /api/candidates/:id/answers   # Step 2: Submit quiz answers
                                     # Body: { answers: [{ question_id, answer_id }] }
                                     # Calculates and stores total_score

POST   /api/candidates/:id/audio     # Step 3: Upload audio (multipart/form-data)
                                     # Saves file to disk, marks candidate as completed
```

### Admin (JWT Protected)

```
POST   /api/admin/login              # Login (email, password) -> JWT token

GET    /api/admin/candidates         # List all candidates ordered by total_score DESC
GET    /api/admin/candidates/:id     # Candidate detail with answers
GET    /api/admin/candidates/:id/audio  # Stream/download candidate audio

GET    /api/admin/questionnaires     # List all questionnaire versions
GET    /api/admin/questionnaires/:id # Questionnaire detail with questions/answers
POST   /api/admin/questionnaires     # Create new questionnaire version (auto-activates)
                                     # Body: { questions: [{ text, sort_order, answers: [{ text, score, sort_order }] }] }
```

## Frontend Routes

### Candidate Flow (multi-step wizard)

```
/                    # Landing / Step 1: Registration form (name, surname, email)
/quiz/:candidateId   # Step 2: Quiz questions (one page, all questions)
/audio/:candidateId  # Step 3: Audio recording (browser MediaRecorder API, max 2 min)
/done                # Thank you page
```

### Admin Area

```
/admin/login                    # Login form
/admin/candidates               # Candidates list (sorted by score, with audio player)
/admin/candidates/:id           # Candidate detail (answers + audio)
/admin/questionnaires           # Questionnaire versions list
/admin/questionnaires/new       # Create new questionnaire (add questions + answers + scores)
/admin/questionnaires/:id       # View questionnaire detail
```

## Audio Recording

- Browser MediaRecorder API for recording
- Max duration: 2 minutes (enforced client-side with timer)
- Format: WebM/Opus (native browser format, widely supported)
- Stored on filesystem at `backend/uploads/audio/{candidate_id}.webm`
- Docker volume mounts `backend/uploads` for persistence

## Auth

- Candidates: no auth, they receive a candidate_id after step 1
- Admins: JWT token, issued on login, sent as `Authorization: Bearer <token>` header
- JWT secret configured via environment variable

## Docker Compose Services

1. **postgres**: Postgres 16, mounts `./docker-entrypoint-initdb.d` to `/docker-entrypoint-initdb.d`
2. **backend**: Go with Air, mounts `./backend` source, exposes port 8080, depends on postgres
3. **frontend**: Node with Vite dev server, mounts `./frontend` source, exposes port 5173

## Seed Data (02_seed.sql)

- 1 admin: `admin@test.com` / `password` (bcrypt hashed)
- 1 active questionnaire with 5 example questions, each with 4 answers (1 correct per question)
- 2 example candidates with answers and scores

## Tech Decisions

- **GORM**: ORM for Go, handles migrations and queries
- **Pinia Colada**: Async state management for Vue, handles API calls with caching
- **DaisyUI**: Component library on top of Tailwind, fast UI development
- **Air**: Go hot reload for development
- **MediaRecorder API**: Native browser audio recording, no external libraries needed
- **UUID**: All primary keys are UUIDs for security (no sequential IDs exposed)
