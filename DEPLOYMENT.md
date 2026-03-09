# Deployment

This repo is set up for:

- `backend` -> Railway
- `frontend` -> Vercel
- database + storage -> Supabase

## 1. Supabase

Create a Supabase project and use:

- Postgres connection string for `DSN`
- project URL for `SUPABASE_URL`
- service role key for `SUPABASE_SERVICE_ROLE_KEY`
- storage bucket name for `SUPABASE_STORAGE_BUCKET`

The backend uses:

- Supabase Postgres through `DSN`
- Supabase Storage for trip photo uploads

It does **not** currently use Supabase Auth. Telegram auth is handled by the backend and the backend issues its own JWT.

## 2. Backend on Railway

Create a Railway service from this repo and set the service root to:

```text
backend
```

Railway will pick up [`backend/railway.toml`](/home/critiq/development/trip-listik/backend/railway.toml) and build with the Dockerfile in that directory.

Required Railway env vars:

```env
DSN=postgresql://...
JWT_SECRET=change-me
TELEGRAM_WEBAPP_SECRET=your_telegram_bot_token
CORS_ORIGINS=https://your-frontend-domain.vercel.app

SUPABASE_URL=https://your-project.supabase.co
SUPABASE_SERVICE_ROLE_KEY=your_service_role_key
SUPABASE_STORAGE_BUCKET=trip-photos
SUPABASE_ANON_KEY=optional
SUPABASE_JWT_SECRET=optional
```

Optional:

```env
HTTP_ADDR=:8080
APP_ENV=prod
BOT_TOKEN=if you also deploy the bot separately
```

Notes:

- Railway injects `PORT`. The backend now respects `PORT` automatically if `HTTP_ADDR` is not set.
- Run DB migrations against the same Supabase Postgres before first production use.

## 3. Frontend on Vercel

Create a Vercel project from this repo and set the root directory to:

```text
frontend
```

Required Vercel env vars:

```env
PUBLIC_API_BASE_URL=https://your-backend.up.railway.app
PUBLIC_SUPABASE_URL=https://your-project.supabase.co
PUBLIC_SUPABASE_BUCKET=trip-photos
```

The frontend uses the Vercel SvelteKit adapter and reads the public env vars at runtime.

## 4. Telegram Mini App

After deploy, use the frontend production URL as the Telegram Mini App URL.

Recommended order:

1. Deploy Supabase and apply migrations
2. Deploy backend to Railway
3. Set backend URL in Vercel env vars
4. Deploy frontend to Vercel
5. Update Telegram Mini App URL
