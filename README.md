# Tourism Reviewer Recruitment Platform

Platform for recruiting tourism location reviewers. Candidates register, take a quiz, record an audio presentation about a tourist location they visited, and get evaluated by AI.

## Architecture

```
Browser --> Vue 3 (Vite, DaisyUI) :5173
              |
              v
           Go API (Gin, GORM) :8080
              |
       +------+------+--------+
       v      v      v        v
   Postgres  Whisper  Ollama  Filesystem
    :5432    :9000   :11434   (uploads/)
```

- **Frontend**: Vue 3 + Composition API, Pinia Colada, DaisyUI/Tailwind
- **Backend**: Go + Gin + GORM, JWT auth, REST API
- **Database**: PostgreSQL 16
- **Whisper**: Audio transcription (faster_whisper, auto language detection)
- **Ollama**: AI analysis with qwen3:4b (tourism reviewer scoring)

## AI Pipeline

When a candidate uploads audio, a background goroutine runs:

1. **Transcription** (Whisper) - audio to text
2. **Analysis** (Ollama) - evaluates on 5 criteria: descriptive skills (0-25), critical thinking (0-25), engagement (0-20), structure (0-15), practical info (0-15)
3. Score, strengths, weaknesses and recommendation saved to DB

Status tracked as: `pending` > `transcribing` > `analyzing` > `completed` / `failed`

## Quick Start

```bash
docker compose up
```

First run downloads the Ollama model (~3GB), takes a few minutes.

- Frontend: http://localhost:5173
- API: http://localhost:8080/api/health
- Admin: http://localhost:5173/admin (login: `admin@test.com` / `password`)

## Candidate Flow

1. Register (name + email)
2. Complete the quiz (one question at a time)
3. Record audio presentation (max 2 min) describing a tourist location
4. Done - AI processes the audio in background

## Admin Panel

- Candidates list with quiz score, AI score, audio player
- Candidate detail: full answers, transcript, AI analysis with sub-scores
- Retry button if AI analysis fails
- Questionnaire management (create, view, activate)

## Dev Notes

- Backend hot-reload via `air`
- Frontend hot-reload via Vite
- Ollama model persisted in `ollama_data` Docker volume
- `docker compose down -v` to reset everything (DB + models)
