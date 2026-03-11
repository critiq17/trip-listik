# ✈️ TripListik

> Backend API for the TripListik Telegram Mini App 🧳

<img src="../assets/images/ReadmeImage.jpg" alt="TripListik Preview" width="150"/>


---

## 🚀 Features

- ✈️ Trips, members, comments, votes, photos
- 🔐 Telegram WebApp auth + JWT
- 🤖 Optional AI integration for travel ideas

---

## 🧩 Tech Stack

- **Go (Fiber + GORM)** — backend API
- **PostgreSQL** — database
- **Docker Compose** — infrastructure setup and container orchestration

---

## ⚙️ Setup & Run

### 1. Configure `.env`

Create a `.env` file in the project root:

```env
# Database connection
DSN="host=localhost user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB sslmode=disable"

# Postgres settings for docker-compose
POSTGRES_USER=your_db_user
POSTGRES_PASSWORD=your_db_password
POSTGRES_DB=trip_listik

# Optional AI API key
API_KEY_AI=your_api_key_here

# API server
HTTP_ADDR=:8080
JWT_SECRET=your_jwt_secret_here
TELEGRAM_WEBAPP_SECRET=your_telegram_bot_token_here
CORS_ORIGINS=http://localhost:5173

# Supabase (optional, for photo uploads)
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your_anon_key
SUPABASE_SERVICE_ROLE_KEY=your_service_role_key
SUPABASE_JWT_SECRET=your_supabase_jwt_secret
SUPABASE_STORAGE_BUCKET=trip-photos

```

### Run 

```docker-compose up -d```


```make run```

### Migrations

```bash
go run ./cmd/migrate -direction=up
```

### API server

```bash
go run ./cmd/api
```
