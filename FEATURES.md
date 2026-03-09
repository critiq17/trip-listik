# FEATURES

This document defines the backend feature plan and implementation details for the Travel Planning system and the Telegram Mini App. It is written to production-grade (Big Tech) standards, with a clear order of implementation and the full set of backend capabilities required by the frontend.

## 1) Goals

- Simple but powerful travel planning tool.
- Organize trips, invite friends, discover trips, share experiences, track travel history.
- World explored progress (starting at 3.5%).
- Telegram Mini App + Bot workflow.
- Supabase-ready backend (Postgres + Auth + Storage + Realtime) with only `.env` setup required.

## 2) Backend Plan (Order of Implementation)

1. Architecture and foundations
- Migrate to Fiber-based HTTP API (keep Telegram bot as a separate entrypoint).
- Introduce structured logging, config validation, and request tracing.
- Add migrations and schema management.

2. Data model and Supabase compatibility
- Define DB schema for all core tables.
- Add Supabase-compatible auth integration and RLS policies.
- Implement storage for trip photos via Supabase Storage.

3. Authentication and identity
- Support Telegram login for Mini App (Telegram WebApp auth flow).
- Issue internal JWT session (or validate Supabase JWT directly).
- Create user profiles and link Telegram IDs.

4. Core features APIs
- Trips CRUD + membership + join requests.
- Feed (public trips) + trip discovery.
- Voting, comments, and photos.

5. Statistics + world map data
- Compute and persist travel statistics.
- Maintain country visit data and world explored percentage.

6. Performance + safety
- Caching, pagination, rate limits.
- Indexes and query tuning.
- Observability: metrics, logs, tracing.

7. Telegram bot integration
- Bot triggers and notifications tied to the same backend.
- Mini App deep links and share links.

## 2.1) Backend Plan Aligned to PAGES_UI Screens

1. Auth and session
- Telegram WebApp initData validation.
- `POST /v1/auth/telegram` returns JWT (or Supabase session).

2. Feed
- `GET /v1/feed` with filters: All, Friends, Popular, Nearby.
- Cursor pagination and sort by recency + engagement.

3. Trip detail
- `GET /v1/trips/:id` returns details, members, votes, photos count, discussion count.
- `GET /v1/trips/:id/members`
- `GET /v1/trips/:id/photos`
- `GET /v1/trips/:id/discussion`

4. Create trip (3 steps)
- `POST /v1/trips` supports draft creation and step-by-step updates.
- `PATCH /v1/trips/:id` for step 2 and 3 updates.

5. My trips
- `GET /v1/trips?scope=mine&status=upcoming|past|draft`

6. Explore
- `GET /v1/explore` with search, country filter, trending.

7. Inbox / Notifications
- `GET /v1/inbox` with join requests, approvals, comments, invites.

8. Profile + stats + world map
- `GET /v1/me`
- `GET /v1/me/stats`
- `GET /v1/me/world`

9. Photos upload
- `POST /v1/trips/:id/photos/presign`
- `POST /v1/trips/:id/photos`

10. Voting + comments
- `POST /v1/trips/:id/votes`
- `GET/POST /v1/trips/:id/comments`

## 3) Architecture

- API server: Go + Fiber.
- DB: PostgreSQL (Supabase-compatible).
- Auth: Supabase Auth or internal JWT with Supabase user linkage.
- Storage: Supabase Storage for photos.
- Realtime: Supabase Realtime for comments and votes (optional), fallback to polling.

## 4) Configuration (.env)

- `DATABASE_URL` (Supabase or local Postgres)
- `SUPABASE_URL`
- `SUPABASE_ANON_KEY`
- `SUPABASE_SERVICE_ROLE_KEY`
- `SUPABASE_JWT_SECRET`
- `TELEGRAM_BOT_TOKEN`
- `TELEGRAM_WEBAPP_SECRET`
- `APP_BASE_URL`
- `CORS_ORIGINS`

## 5) Data Model

### users
- `id` (uuid, PK)
- `telegram_id` (bigint, unique)
- `username` (text)
- `first_name` (text)
- `last_name` (text)
- `photo_url` (text)
- `created_at` (timestamptz)
- `updated_at` (timestamptz)

### trips
- `id` (uuid, PK)
- `owner_id` (uuid, FK -> users)
- `title` (text)
- `description` (text)
- `start_date` (date)
- `end_date` (date)
- `visibility` (enum: public, private)
- `status` (enum: draft, planned, completed, canceled)
- `country_code` (text)  -- ISO 3166-1 alpha-2
- `city` (text)
- `cover_photo_url` (text)
- `created_at` (timestamptz)
- `updated_at` (timestamptz)

### trip_members
- `id` (uuid, PK)
- `trip_id` (uuid, FK -> trips)
- `user_id` (uuid, FK -> users)
- `role` (enum: owner, member)
- `joined_at` (timestamptz)

### trip_join_requests
- `id` (uuid, PK)
- `trip_id` (uuid, FK -> trips)
- `user_id` (uuid, FK -> users)
- `status` (enum: pending, approved, rejected)
- `created_at` (timestamptz)
- `updated_at` (timestamptz)

### trip_votes
- `id` (uuid, PK)
- `trip_id` (uuid, FK -> trips)
- `user_id` (uuid, FK -> users)
- `vote` (smallint)  -- 1..5 or -1/1 (define final behavior)
- `created_at` (timestamptz)
- `updated_at` (timestamptz)

### trip_comments
- `id` (uuid, PK)
- `trip_id` (uuid, FK -> trips)
- `user_id` (uuid, FK -> users)
- `body` (text)
- `created_at` (timestamptz)
- `updated_at` (timestamptz)

### trip_photos
- `id` (uuid, PK)
- `trip_id` (uuid, FK -> trips)
- `user_id` (uuid, FK -> users)
- `storage_path` (text) -- Supabase Storage path
- `url` (text)
- `created_at` (timestamptz)

### user_countries
- `id` (uuid, PK)
- `user_id` (uuid, FK -> users)
- `country_code` (text)
- `first_visit_at` (timestamptz)

### user_stats (optional materialized table)
- `user_id` (uuid, PK)
- `total_trips` (int)
- `countries_visited` (int)
- `cities_visited` (int)
- `trips_with_friends` (int)
- `solo_trips` (int)
- `updated_at` (timestamptz)

## 6) API Surface (Required by Frontend)

### Auth
- `POST /v1/auth/telegram`:
  - Validate Telegram WebApp initData.
  - Create or update user.
  - Return session token (JWT) or Supabase session.

### Users
- `GET /v1/me`:
  - Returns profile + stats + world explored percent.

### Feed
- `GET /v1/feed?cursor=&limit=`:
  - Public trips sorted by recent activity.

### Trips
- `POST /v1/trips`:
  - Create trip.
- `GET /v1/trips/:id`:
  - Trip details, members, vote summary, comments count, photos.
- `PATCH /v1/trips/:id`:
  - Update trip.
- `DELETE /v1/trips/:id`:
  - Soft delete (optional).

### Members
- `POST /v1/trips/:id/join`:
  - Request to join trip.
- `POST /v1/trips/:id/join/approve`:
  - Owner approves join request.
- `POST /v1/trips/:id/join/reject`:
  - Owner rejects join request.
- `DELETE /v1/trips/:id/members/:userId`:
  - Remove member (owner only).

### Votes
- `POST /v1/trips/:id/votes`:
  - Cast or update vote.
- `GET /v1/trips/:id/votes`:
  - Aggregate stats (avg, count).

### Comments
- `GET /v1/trips/:id/comments?cursor=&limit=`
- `POST /v1/trips/:id/comments`

### Photos
- `POST /v1/trips/:id/photos/presign`:
  - Returns signed upload URL (Supabase Storage).
- `POST /v1/trips/:id/photos`:
  - Save photo metadata after upload.
- `DELETE /v1/trips/:id/photos/:photoId`

### Stats
- `GET /v1/me/stats`:
  - Full travel statistics.
- `GET /v1/me/world`:
  - Visited countries and world explored percentage.

## 7) Supabase Integration Details

- Auth:
  - Validate Telegram initData on backend.
  - Either create a Supabase user via service role or maintain internal JWT and map to user.
  - Prefer Supabase JWT so frontend can directly access Supabase Realtime/Storage with RLS.

- RLS Policies:
  - `users`: user can read own profile.
  - `trips`: public trips readable by all; private trips only by members.
  - `trip_members`, `trip_join_requests`, `trip_votes`, `trip_comments`, `trip_photos`: members only, owner can manage.

- Storage:
  - Bucket: `trip-photos`.
  - Path format: `trip_photos/{trip_id}/{photo_id}.jpg`.
  - Signed upload via backend using service role.

## 8) Travel Statistics

Definitions:
- `total_trips`: total trips created by user or joined.
- `countries_visited`: distinct country codes from completed trips or user_countries.
- `cities_visited`: distinct city names from completed trips.
- `trips_with_friends`: trips with member count > 1.
- `solo_trips`: trips with member count == 1.

World explored:
- Use static list of countries (ISO 3166-1 alpha-2) count = 195.
- `world_explored_percent = countries_visited / 195 * 100`.
- Start from 3.5% as display until user has data.

Implementation:
- Update `user_countries` on trip completion.
- Recompute `user_stats` via background job or on-demand with caching.

## 9) Telegram Mini App

### Minimal UI Structure
- Feed: public trips list.
- Trip Page: details, members, voting, discussion, photos.
- Create Trip: form.
- Profile: trip history, world map, travel stats.

### Mini App Backend Needs
- Auth endpoint with Telegram initData validation.
- Feed and trip detail endpoints optimized for mobile.
- Presigned photo uploads.
- Realtime comments/votes (optional). If not, short polling.

## 10) Svelte Mini App (Design Placeholder)

The user will provide HTML design. Implementation plan in Svelte:

1. Project setup
- SvelteKit for routing and SSR where needed.
- Use Vite + Svelte adapter for static or node hosting.

2. Performance
- Use route-based code splitting.
- Fetch data with streaming and skeleton states.
- Use image lazy loading and responsive sizes.

3. UI/UX and animation
- Minimal layout, whitespace, strong typography.
- Use `svelte/animate` and `svelte/transition` for smooth transitions.
- Staggered list transitions for feed, comments, photos.
- Avoid heavy libraries; use CSS variables, `prefers-reduced-motion` support.

4. Integration
- Telegram WebApp init data parsed on load.
- Call `POST /v1/auth/telegram` and store session token.
- Use token for subsequent API calls.

5. Feature mapping
- Feed: `/v1/feed`.
- Trip details: `/v1/trips/:id`.
- Create trip: `/v1/trips`.
- Profile: `/v1/me`, `/v1/me/stats`, `/v1/me/world`.

## 11) Non-Functional Requirements

- Latency: p95 < 200ms for reads.
- Errors: structured error response with code + message.
- Pagination: cursor-based.
- Rate limits: per user and per IP.
- Auditing: track updates to trips and members.
- Security: input validation, SQL injection safe, RLS enforced.

## 12) Implementation Notes

- Add migration tooling (golang-migrate or goose).
- Introduce DTOs and input validation (go-playground/validator).
- Use context timeouts for DB and HTTP.
- Use dependency injection for services.
- Add tests: repository, service, and API handlers.
