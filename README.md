# TripListik Monorepo

This repo is split into:

- `backend/` Go API + migrations
- `frontend/` SvelteKit Telegram Mini App

## Backend

```bash
cd backend
# set .env (see backend/README.md)

go run ./cmd/migrate -direction=up

go run ./cmd/api
```

## Frontend

```bash
cd frontend
npm run dev
```

## Docker Compose (All-in-One)

```bash
docker-compose up --build
```
