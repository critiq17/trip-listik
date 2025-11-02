# ✈️ TripListik

> Telegram bot to easily save, manage, and explore your travel wishlist 🧳

![TripListik Preview](./assets/images/ReadmeImage.jpg)

---

## 🚀 Features

- 💬 Manage your personal travel wishlist via Telegram
- 🗺️ Add and remove destinations
- 📸 User profiles with photo upload support
- 🤖 Optional AI integration for travel ideas

---

## 🧩 Tech Stack

- **Go (Fiber + GORM)** — backend & Telegram bot
- **PostgreSQL** — database
- **SvelteKit** — frontend for user profiles
- **Docker Compose** — infrastructure setup and container orchestration

---

## ⚙️ Setup & Run

### 1. Configure `.env`

Create a `.env` file in the project root:

```env
# Database connection
DSN="host=localhost user=myuser password=1234 dbname=trip_listik sslmode=disable"

# Telegram Bot Token
BOT_TOKEN=your_telegram_bot_token_here

# Optional AI API key
API_KEY_AI=your_api_key_here

```

### Run 

```docker-compose up -d```


```make run```