# TripListik

A Telegram Mini App for collaborative trip planning. Plan trips, invite friends, vote on destinations, and coordinate travel — all inside Telegram.

Landing: [triplistik.netlify.app](https://triplistik.netlify.app)

---

## Structure

```
backend/    Go API (Fiber + PostgreSQL)
frontend/   SvelteKit Telegram Mini App
landing/    Static landing page (SvelteKit)
docs/       API documentation
```

## Quick Start

### Docker Compose

```bash
cp .env.example .env   # fill in values
docker-compose up --build
```

Migrations run automatically before the backend starts.

### Manual

```bash
# Backend
cd backend
go run ./cmd/migrate -direction=up
go run ./cmd/api

# Frontend
cd frontend
pnpm install
pnpm dev
```

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `DSN` | yes | PostgreSQL connection string |
| `JWT_SECRET` | yes | JWT signing secret |
| `TELEGRAM_WEBAPP_SECRET` | yes | Bot token for Mini App init data validation |
| `BOT_TOKEN` | no | Telegram Bot API token (enables notifications) |
| `MINI_APP_URL` | no | Public Mini App URL (used in notification buttons) |
| `SUPABASE_URL` | no | Supabase project URL |
| `SUPABASE_ANON_KEY` | no | Supabase anon key |
| `SUPABASE_SERVICE_ROLE_KEY` | no | Supabase service role key |
| `SUPABASE_STORAGE_BUCKET` | no | Storage bucket name (default: `trip-photos`) |
| `CORS_ORIGINS` | no | Allowed CORS origins (default: `*`) |

Frontend build vars:

| Variable | Description |
|---|---|
| `PUBLIC_API_BASE_URL` | Backend URL |
| `PUBLIC_SUPABASE_URL` | Supabase URL for photo display |
| `PUBLIC_SUPABASE_BUCKET` | Bucket name for photo display |

## Tests

```bash
cd backend
TEST_DATABASE_URL="host=localhost user=myuser password=1234 dbname=trip_listik_test sslmode=disable" \
  go test ./internal/trips/... ./internal/invites/...
```

## Docs

- [API Reference](docs/trips-api.md)
- [Invite Flow](docs/invite-guide.md)
