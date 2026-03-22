# TripListik Monorepo

A Telegram Mini App for collaborative trip planning.

This repo is split into:

- `backend/` — Go API + migrations
- `frontend/` — SvelteKit Telegram Mini App
- `docs/` — API documentation and guides

---

## Quick Start

### Backend

```bash
cd backend
# copy .env.example to .env and fill in values
go run ./cmd/migrate -direction=up
go run ./cmd/api
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Docker Compose (all-in-one)

```bash
cp .env.example .env   # fill in values
docker-compose up --build
```

The `migrate` service runs automatically before the backend starts.

---

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `DSN` | yes | PostgreSQL connection string |
| `JWT_SECRET` | yes | Secret for signing JWT tokens |
| `TELEGRAM_WEBAPP_SECRET` | yes | Bot token used to validate Mini App init data |
| `BOT_TOKEN` | no | Telegram Bot API token (enables TG notifications) |
| `MINI_APP_URL` | no | Public URL of the Mini App (used in TG notification buttons) |
| `SUPABASE_URL` | no | Supabase project URL (for photo storage) |
| `SUPABASE_ANON_KEY` | no | Supabase anon key |
| `SUPABASE_SERVICE_ROLE_KEY` | no | Supabase service role key |
| `SUPABASE_STORAGE_BUCKET` | no | Storage bucket name (default: `trip-photos`) |
| `CORS_ORIGINS` | no | Allowed CORS origins (default: `*`) |

Frontend build args:

| Variable | Description |
|---|---|
| `PUBLIC_API_BASE_URL` | Backend URL visible from the browser |
| `PUBLIC_SUPABASE_URL` | Supabase URL for photo display |
| `PUBLIC_SUPABASE_BUCKET` | Bucket name for photo display |

---

## API Overview

Full API documentation: [`docs/trips-api.md`](docs/trips-api.md)

### Trips

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/v1/trips` | required | Create trip |
| `GET` | `/v1/trips/:id` | optional | Get trip detail |
| `PATCH` | `/v1/trips/:id` | required | Update trip (owner only) |
| `DELETE` | `/v1/trips/:id` | required | Delete trip (owner only) |
| `GET` | `/v1/trips?scope=mine` | required | List my trips |

### Invites

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/v1/trips/:id/invite` | required | Invite user by username or ID |
| `GET` | `/v1/trips/:id/invites` | required | List trip invites |
| `POST` | `/v1/invites/:id/respond` | required | Accept or decline invite |
| `DELETE` | `/v1/invites/:id` | required | Cancel invite (owner only) |

### Access rights

| Action | Who |
|---|---|
| View public trip | Anyone |
| View private trip | Trip members only |
| Edit / delete trip | Trip owner only |
| Send invite | Owner; or any member (group trips) |
| Accept / decline | The invited user |
| Cancel invite | Trip owner only |

---

## Telegram Notifications

When `BOT_TOKEN` and `MINI_APP_URL` are configured:

- **Invite sent** → invitee receives a TG message with a [View Invite] button
- **Invite accepted/declined** → TG message is deleted from invitee's chat; owner is notified
- **Invite cancelled** → TG message is deleted from invitee's chat

See [`docs/invite-guide.md`](docs/invite-guide.md) for the full flow.

---

## Running Tests

```bash
# Backend integration tests (requires a test PostgreSQL DB)
cd backend
TEST_DATABASE_URL="host=localhost user=myuser password=1234 dbname=trip_listik_test sslmode=disable" \
  go test ./internal/trips/... ./internal/invites/...
```
