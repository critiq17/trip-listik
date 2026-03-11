# 🚀 Feature Implementation Plans

Детальные планы для всех запрошенных фич.

---

## Feature 1 — UI Consistency + Bottom Nav Минимализм

### Проблема

- Trip cards на разных страницах выглядят по-разному (feed vs trips page)
- Bottom nav громоздкий: `add_circle` FAB выпирает, labels занимают место
- Нет единого компонента карточки поездки

### План

#### 1.1 Design System расширение

##### [MODIFY] [tokens.css](file:///home/critiq/development/trip-listik/frontend/src/lib/styles/tokens.css)

- Добавить недостающие токены: `--radius-2xl: 16px`, `--radius-3xl: 24px`
- Добавить `--transition-spring`, `--transition-smooth` для единых анимаций
- Добавить `--card-bg`, `--card-border` для унифицированных карточек
- Добавить `--nav-height` для единого отступа

#### 1.2 Единый TripCard компонент

##### [MODIFY] [TripCard.svelte](file:///home/critiq/development/trip-listik/frontend/src/lib/components/TripCard.svelte)

- Добавить prop `variant: 'feed' | 'compact' | 'horizontal'`
- `feed` — текущий вид с 4/5 aspect ratio
- `compact` — как на trips page, 16/10 ratio + status badge + date pill
- `horizontal` — для history list, маленький thumbnail + info
- Glassmorphism vote badge
- Hover: subtle scale + glow ring animation

#### 1.3 Bottom Nav минимализм

##### [MODIFY] [BottomNav.svelte](file:///home/critiq/development/trip-listik/frontend/src/lib/components/BottomNav.svelte)

**Новый дизайн:**

- Убрать labels по умолчанию — только иконки (4 + 1)
- Label появляется ТОЛЬКО у активного таба с `fly` анимацией
- FAB «Create» — маленькая точка-индикатор вместо громоздкой выпуклости
- При нажатии на любой таб: иконка scale spring bounce `[1 → 0.85 → 1.1 → 1]`
- Active tab: заливка иконки (FILL 1) + dot indicator снизу вместо text
- Фон: `backdrop-filter: blur(24px)` + чуть прозрачнее
- Общая высота ≈ 52px вместо текущих ~72px

```
┌─────────────────────────────────────────────┐
│   🏠     🗺️     ➕     🔍     👤           │
│          •               Explore            │
└─────────────────────────────────────────────┘
          ↑ active: dot + label fade-in
```

**Анимации:**

1. Icon tap → spring bounce `{scale: [1, 0.85, 1.12, 1], duration: 0.35}`
2. Active label → `fly|fade {y: 4, duration: 200}`
3. Indicator dot → `spring` slide between tabs (shared element)
4. Page transition → icons crossfade fill ↔ outline

#### 1.4 Trips Page карточки

##### [MODIFY] [trips/+page.svelte](file:///home/critiq/development/trip-listik/frontend/src/routes/trips/+page.svelte)

- Заменить inline card markup на `<TripCard trip={trip} variant="compact" />`
- Удалить 200+ строк дублированных стилей

#### 1.5 Profile page history

##### [MODIFY] [profile/+page.svelte](file:///home/critiq/development/trip-listik/frontend/src/routes/profile/+page.svelte)

- Заменить inline history items на `<TripCard trip={trip} variant="horizontal" />`

---

## Feature 2 — Интерактивная карта мира

### Концепция

SVG-карта мира на профиле с цветовой heatmap: чем больше поездок в страну — тем насыщеннее зелёный. Неизвестные страны — тёмные/прозрачные.

### План

#### 2.1 Backend — данные для карты

##### [NEW] `GET /v1/me/map` endpoint

##### [MODIFY] [profile.go](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/profile.go)

```json
// Response
{
  "countries": [
    { "code": "JP", "visit_count": 3, "last_visit": "2025-08-01" },
    { "code": "IT", "visit_count": 1, "last_visit": "2024-12-15" }
  ],
  "total_countries": 12,
  "world_explored_percent": 6.15
}
```

##### [MODIFY] [stats.go](file:///home/critiq/development/trip-listik/backend/internal/store/stats.go)

- Добавить `GetUserCountryVisits(ctx, userID) → []CountryVisit`
- SQL: `SELECT country_code, COUNT(*) as visit_count, MAX(end_date) as last_visit FROM trips JOIN trip_members... GROUP BY country_code`

##### [NEW] Migration `0002_add_city_coordinates.up.sql`

- Добавить `latitude`, `longitude` к `trips` table для будущей карты городов

#### 2.2 Frontend — WorldMap компонент

##### [NEW] `WorldMap.svelte`

**Технология:** Inline SVG мира (Natural Earth projection)

- Скачать public domain SVG карту, включить inline
- Каждый `<path>` имеет `data-country="XX"` (ISO alpha-2)

**Визуализация:**
| Кол-во визитов | Цвет | Opacity |
|---|---|---|
| 0 (не был) | `#1a1f1c` (тёмный, почти невидимый) | 0.3 |
| 1 визит | `#2d6b44` (начальный зелёный) | 0.6 |
| 2-3 визита | `#4d9d6d` (основной primary) | 0.8 |
| 4+ визитов | `#7fbf99` (яркий зелёный) + glow | 1.0 + `drop-shadow` |

**Интерактивность:**

1. **Hover** → tooltip с `{country_name} — {visit_count} trips`
2. **Tap** на мобильном → bottom sheet с деталями (поездки в эту страну)
3. **Mount animation** → stagger fill animation по странам, самые посещённые — в конце для эффекта reveal
4. **Pulse dots** → маленькие `<circle>` с `pulseGlow` анимацией на последних посещённых городах

**Zoom/Pan:**

- `viewBox` manipulation через touch gestures (pinch-to-zoom)
- Или использовать D3-geo для более продвинутой проекции

**Компоненты:**

```
WorldMap.svelte          — основной SVG контейнер
├── CountryPath.svelte   — отдельный path с hover/tap
├── CityDot.svelte       — pulse dot для городов
└── MapTooltip.svelte    — floating tooltip
```

#### 2.3 Интеграция в профиль

##### [MODIFY] [profile/+page.svelte](file:///home/critiq/development/trip-listik/frontend/src/routes/profile/+page.svelte)

- Заменить `"View Map"` кнопку на inline `<WorldMap>` компонент
- Добавить "World Explored" секцию с анимированным процентом + progress ring SVG
- Под картой — список «Top countries» с флагами и счётчиками

---

## Feature 3 — Инвайты в поездку + Telegram сообщения

### Концепция

Пользователь может найти друзей по username/имени, отправить приглашение. Друг получает сообщение через Telegram бота и может принять, отклонить, или оставить комментарий.

### План

#### 3.1 Backend API

##### [NEW] `POST /v1/trips/:id/invite` — отправить приглашение

```json
// Request
{ "username": "john_doe" }
// или
{ "user_id": "uuid-xxx" }

// Response
{ "invite_id": "uuid", "status": "pending" }
```

##### [NEW] `GET /v1/trips/:id/invites` — список приглашений (для owner)

##### [NEW] `POST /v1/invites/:id/respond` — ответ на приглашение

```json
// Request
{ "action": "accept" | "decline", "comment": "Извини, я в те даты занят :(" }
```

##### [NEW] `GET /v1/users/search?q=` — поиск пользователей

```json
// Response
{
  "items": [
    {
      "id": "uuid",
      "username": "john",
      "first_name": "John",
      "photo_url": "..."
    }
  ]
}
```

#### 3.2 Database

##### [NEW] Migration `0003_invite_comments.up.sql`

- Добавить `comment` TEXT к `trip_invites` table
- Добавить индекс `users(username)` для быстрого поиска

#### 3.3 Telegram Bot интеграция

##### [NEW] `internal/bot/notifications.go`

Когда создаётся invite:

1. Backend ищет `telegram_id` приглашённого пользователя
2. Отправляет Telegram сообщение через Bot API:

```
🗺 Тебя приглашают в поездку!

✈️ "Summer in Tuscany"
📅 15 Jun – 28 Jun 2026
👤 Пригласил: @alex_traveler

[Принять ✅] [Отклонить ❌] [Открыть 🔗]
```

3. Inline keyboard с callback queries `/accept:{invite_id}`, `/decline:{invite_id}`
4. При нажатии — бот вызывает `/v1/invites/:id/respond` через internal API

#### 3.4 Frontend — Invite UI

##### [NEW] `InviteModal.svelte` (bottom sheet)

**Шаги:**

1. На trip detail page, Members tab → кнопка «Invite friends»
2. Открывается bottom sheet с search input
3. Поиск по username/имени → real-time debounced запросы к `/v1/users/search`
4. Результаты — список с avatar + username + «Invite» кнопка
5. При нажатии «Invite» — POST `/v1/trips/:id/invite` + toast «Invitation sent ✉️»
6. Список pending invites внизу с статусами

##### [MODIFY] [trips/[id]/+page.svelte](file:///home/critiq/development/trip-listik/frontend/src/routes/trips/%5Bid%5D/+page.svelte)

- Добавить «Invite» кнопку в Members tab
- Показывать pending invites рядом с join requests для owner

#### 3.5 Notification интеграция

##### [MODIFY] [inbox/+page.svelte](file:///home/critiq/development/trip-listik/frontend/src/routes/inbox/+page.svelte)

- Показывать входящие инвайты с кнопками Accept/Decline/Comment
- Карточка инвайта: trip cover + details + action buttons + comment input

---

## Feature 4 — Публичные профили

### Концепция

Каждый пользователь имеет публичную страницу с историей поездок, статистикой, картой, и "wishlist" — куда хочет поехать.

### План

#### 4.1 Backend

##### [NEW] `GET /v1/users/:id/profile` — публичный профиль (no auth required)

```json
{
  "user": { "id": "...", "username": "alex", "first_name": "Alex", "photo_url": "..." },
  "stats": { "total_trips": 14, "countries_visited": 8, "cities_visited": 22 },
  "world_explored_percent": 4.1,
  "public_trips": [
    { "id": "...", "title": "Iceland Adventure", "cover_photo_url": "...", "country_code": "IS", ... }
  ],
  "country_map": [
    { "code": "IS", "visit_count": 2 }
  ],
  "wishlist": [
    { "country_code": "JP", "note": "Хочу увидеть Фудзи" }
  ]
}
```

##### [NEW] `GET/POST/DELETE /v1/me/wishlist` — управление списком желаний

##### [NEW] Migration `0004_wishlist_and_bio.up.sql`

```sql
CREATE TABLE user_wishlist (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id),
  country_code text NOT NULL,
  city text,
  note text,
  created_at timestamptz DEFAULT now(),
  UNIQUE(user_id, country_code)
);

ALTER TABLE users ADD COLUMN bio text DEFAULT '';
ALTER TABLE users ADD COLUMN is_public boolean DEFAULT true;
```

#### 4.2 Frontend

##### [NEW] `/profile/[username]/+page.svelte` — публичный профиль

**Layout:**

```
┌─────────────────────────────────────────┐
│ [← Back]                    [Share 🔗]  │
│                                         │
│        [Avatar with gradient ring]      │
│         @alex_traveler                  │
│    "I chase sunsets and cheap flights"  │
│                                         │
│   ┌────────────────────────────────┐    │
│   │  14 trips · 8 countries · 22🏙️ │    │
│   └────────────────────────────────┘    │
│                                         │
│   [========= World Map ===========]    │
│                                         │
│   ── Public Trips ──                    │
│   [TripCard: Iceland]                   │
│   [TripCard: Italy]                     │
│   [TripCard: Japan]                     │
│                                         │
│   ── Wishlist 🌟 ──                     │
│   🇯🇵 Japan — "Хочу увидеть Фудзи"      │
│   🇳🇿 New Zealand — "Бандж-джамп"       │
└─────────────────────────────────────────┘
```

**Анимации:**

- Stats counter count-up при mount
- Trip cards stagger fade-in
- World map страны animate на mount
- Wishlist items slide-in с delay

##### [MODIFY] [profile/+page.svelte](file:///home/critiq/development/trip-listik/frontend/src/routes/profile/+page.svelte)

- Добавить «Share profile» кнопку → копирует ссылку `t.me/bot?startapp=profile_alex`
- Добавить секцию «Wishlist» с возможностью добавлять страны
- Добавить «Bio» input в Edit profile

---

## Feature 5 — Фотоальбом и описание поездки (post-trip / during-trip)

### Концепция

После (или во время) поездки — можно добавить фотографии, написать описание, сделать "дневник поездки" с записями по дням.

### План

#### 5.1 Backend

##### [NEW] Migration `0005_trip_diary.up.sql`

```sql
CREATE TABLE trip_diary_entries (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  trip_id uuid NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  day_number int,
  title text,
  body text NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- photos уже есть, добавим caption
ALTER TABLE trip_photos ADD COLUMN caption text DEFAULT '';
ALTER TABLE trip_photos ADD COLUMN diary_entry_id uuid REFERENCES trip_diary_entries(id);
```

##### [NEW] CRUD endpoints:

- `POST /v1/trips/:id/diary` — создать запись
- `GET /v1/trips/:id/diary` — список записей с фото
- `PATCH /v1/trips/:id/diary/:entryId` — редактировать
- `DELETE /v1/trips/:id/diary/:entryId` — удалить

##### [MODIFY] [photos.go](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/photos.go)

- Добавить `caption` в create/update photo
- Привязка фото к diary entry

#### 5.2 Frontend

##### [MODIFY] Trip Detail — Photos tab

- Добавить «caption» input при загрузке фото
- Masonry grid с captions
- Lightbox viewer: swipe between photos + caption overlay

##### [NEW] Trip Diary tab (или секция в Details)

Добавить в trip detail tabs → «Story» tab:

```
┌─────────────────────────────────────────┐
│  ── Day 1 · Arrival ──                  │
│  [Full-width hero photo]                │
│  "Landed in Reykjavik at midnight.      │
│   The sun was still up. 🌅"             │
│                                         │
│  [Photo grid: 3 photos with captions]   │
│                                         │
│  ── Day 2 · Golden Circle ──            │
│  [Photo]                                │
│  "Geysir was insane..."                 │
│                                         │
│  [+ Add entry]                          │
└─────────────────────────────────────────┘
```

**Компоненты:**

```
DiaryEntry.svelte     — один день/запись с текстом и фотками
DiaryEditor.svelte    — форма создания/редактирования записи
PhotoUploader.svelte  — drag & drop + caption input
```

**UX Flow:**

1. Trip status = `planned` → show «Start your travel diary» CTA
2. Trip status = `completed` → show «Write about your trip» CTA
3. Во время поездки (между start_date и end_date) → push notification «How's day ${N}?»
4. Diary entries автоматически нумеруются по дням от start_date

---

## Feature 6 — Bottom Nav: минималистичная переработка (детальная спецификация)

### Текущие проблемы

- Labels всегда видны — занимают место
- FAB [Create](file:///home/critiq/development/trip-listik/backend/internal/store/trips.go#12-15) с `translateY(-6px)` выглядит "торчащим"
- Анимация bounce — минимальная

### Новый дизайн

##### [MODIFY] [BottomNav.svelte](file:///home/critiq/development/trip-listik/frontend/src/lib/components/BottomNav.svelte)

**Спецификация:**

```css
/* Общая высота: 48px + safe-area */
.nav {
  height: 48px;
  padding: 0 1.5rem;
  gap: 0;
  background: rgba(22, 28, 24, 0.75);
  backdrop-filter: blur(28px) saturate(1.8);
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

/* Каждый таб: только иконка */
.tab {
  width: 48px;
  height: 48px;
  display: grid;
  place-items: center;
  position: relative;
}

/* Active indicator — точка под иконкой */
.dot {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--primary);
  position: absolute;
  bottom: 6px;
  /* animated slide between tabs */
}
```

**Анимации (Motion One):**

1. **Tab switch — icon morph:**

   ```js
   // Outgoing tab
   animate(
     prevIcon,
     { scale: [1, 0.8], opacity: [1, 0.5] },
     { duration: 0.15 },
   );
   // Incoming tab
   animate(
     nextIcon,
     { scale: [0.8, 1.15, 1], opacity: [0.5, 1] },
     { duration: 0.3, easing: spring() },
   );
   ```

2. **Dot indicator slide:**

   ```js
   // Shared dot slides horizontally
   animate(
     dot,
     { x: [prevX, nextX] },
     { duration: 0.25, easing: [0.34, 1.56, 0.64, 1] },
   );
   ```

3. **Create button — special:**
   - Вместо FAB — обычная иконка `add`, НО при нажатии:

   ```js
   // Rotate 45° (becomes ×), then navigate
   animate(plusIcon, { rotate: ["0deg", "135deg"] }, { duration: 0.2 });
   ```

   - Если на create page — иконка становится `×` (close/cancel)

4. **Micro-interaction ripple:**
   ```js
   // CSS circular ripple from tap point
   @keyframes ripple {
     from { scale: 0; opacity: 0.4; }
     to { scale: 2.5; opacity: 0; }
   }
   ```

---

## Порядок реализации

| Приоритет | Feature                            | Effort          | Зависимости                |
| --------- | ---------------------------------- | --------------- | -------------------------- |
| 🥇 1      | Bottom Nav redesign                | S (2-3 часа)    | Нет                        |
| 🥇 2      | UI consistency (TripCard variants) | M (4-5 часов)   | Нет                        |
| 🥈 3      | Интерактивная карта мира           | L (8-10 часов)  | Backend: `/v1/me/map`      |
| 🥈 4      | Инвайты + Telegram messages        | L (10-12 часов) | Bot integration, migration |
| 🥉 5      | Публичные профили                  | M (6-8 часов)   | Migration, new route       |
| 🥉 6      | Фотоальбом + дневник               | L (8-10 часов)  | Migration, new endpoints   |

> [!TIP]
> Рекомендую начать с **Bottom Nav + UI consistency** — это быстрый visual win, потом **карта мира** для wow-эффекта, и дальше по плану.
