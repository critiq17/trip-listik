# 🔍 Code Review — Trip-Listik

Полный анализ кодовой базы: Go backend (Fiber + GORM) и SvelteKit frontend (Svelte 5 + Motion One).

---

## Общая оценка

| Аспект | Оценка | Комментарий |
|---|---|---|
| Архитектура | ⭐⭐⭐ | Хорошее разделение слоёв, но есть дублирование и мёртвый код |
| Безопасность | ⭐⭐⚠️ | Auth работает, но есть дыры в валидации и авторизации |
| Производительность | ⭐⭐⭐ | Cursor pagination, индексы есть, но N+1 запросы и нет кэширования |
| Код frontend | ⭐⭐⭐ | Animations красивые, но Svelte 4 + Svelte 5 API смешаны |
| Тесты | ⚠️ | Только 1 тест (jwt_test.go), остальное — ноль покрытия |
| DevOps | ⭐⭐⭐ | Docker Compose, migrations, health checks — всё на месте |

---

## 🔴 Критические проблемы

### 1. Два мёртвых entrypoint — [app.go](file:///home/critiq/development/trip-listik/backend/internal/app/app.go) и [cmd/main.go](file:///home/critiq/development/trip-listik/backend/cmd/main.go)

[app.go](file:///home/critiq/development/trip-listik/backend/internal/app/app.go) — это старый entrypoint для Telegram бота с `gorm.Open()` и package `repository/services/state`. Новый API запускается через [cmd/api/main.go](file:///home/critiq/development/trip-listik/backend/cmd/api/main.go) с другим `store` слоем.

**Проблема:** Две параллельные архитектуры (`repository/services/state` vs `store/httpapi/server`) не стыкуются. Бот и API не могут работать вместе.

### 2. Мёртвый `models` пакет (`internal/models`)

[model.go](file:///home/critiq/development/trip-listik/backend/internal/models/model.go) содержит старые [User](file:///home/critiq/development/trip-listik/backend/internal/store/models/models.go#9-19)/[Place](file:///home/critiq/development/trip-listik/backend/internal/models/model.go#13-18) модели с `gorm.Model`. Это конфликтует с новыми моделями в [store/models/models.go](file:///home/critiq/development/trip-listik/backend/internal/store/models/models.go).

### 3. SQL Injection через ILIKE

В [trips.go:122-123](file:///home/critiq/development/trip-listik/backend/internal/store/trips.go#L122-L123):
```go
like := "%" + query + "%"
q = q.Where("title ILIKE ? OR city ILIKE ?", like, like)
```
Хоть GORM параметризует запросы, символы `%` и `_` в input'е меняют семантику ILIKE. Нужна санитизация спецсимволов.

### 4. Отсутствует `auth_date` проверка в Telegram initData

[telegram.go](file:///home/critiq/development/trip-listik/backend/internal/auth/telegram.go) валидирует HMAC-подпись, но не проверяет `auth_date`. Это позволяет реюзать старый `initData` бесконечно → replay attack.

### 5. Транзакция в ApproveJoin использует неправильный контекст

В [members.go:170-175](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/members.go#L170-L175):
```go
err = h.Store.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    if err := h.Store.UpdateJoinRequestStatus(ctx, tripID, userID, "approved"); err != nil {
        return err
    }
    return h.Store.AddTripMember(ctx, tripID, userID, "member")
})
```
Операции внутри транзакции используют `h.Store.DB` (оригинальное соединение), а не `tx`. **Транзакция ничего не защищает.**

---

## 🟡 Серьёзные проблемы

### 6. [ListMyTrips](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/trips.go#237-279) — фильтрация в Go вместо SQL

[trips.go:252-275](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/trips.go#L252-L275) загружает **все** поездки пользователя, а потом фильтрует по `status` в Go-коде. При 1000+ поездках это убьёт производительность.

### 7. [ComputeUserStats](file:///home/critiq/development/trip-listik/backend/internal/store/stats.go#17-60) — 3 отдельных SQL запроса

[stats.go](file:///home/critiq/development/trip-listik/backend/internal/store/stats.go) делает 3 отдельных `Raw()` запроса вместо одного оптимизированного. Плюс, второй `Scan()` перезаписывает `TotalTrips` нулём, потому что он не входит в SELECT.

### 8. SSE stream без graceful disconnect

[stream.go](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/stream.go) не слушает context cancellation. Если клиент отключится, goroutine зависнет навсегда до следующего ping/event. Утечка горутин.

### 9. JWT TTL = 30 дней без refresh token

[auth.go:67](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/auth.go#L67): `30*24*time.Hour` — слишком долго для single token без механизма отзыва.

### 10. Frontend: Svelte 4 ([$:](file:///home/critiq/development/trip-listik/.env)) и Svelte 5 (`$props()`) смешаны

[+layout.svelte](file:///home/critiq/development/trip-listik/frontend/src/routes/+layout.svelte) использует Svelte 5 runes (`$props()`), но [+page.svelte](file:///home/critiq/development/trip-listik/frontend/src/routes/+page.svelte) и все другие — реактивные statements [$:](file:///home/critiq/development/trip-listik/.env) и `export let`. Это работает, но создаёт путаницу и несовместимости.

### 11. [GetTrip](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/trips.go#104-163) — public endpoint без auth, но подсчёты — 4 отдельных запроса

[trips.go:135-150](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/trips.go#L135-L150) — каждый [GetTrip](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/trips.go#104-163) запрос делает **5 SQL-запросов** (trip + member_count + vote_stats + comment_count + photo_count).

### 12. Нет пагинации для photos и members

Endpoints `GET /trips/:id/members` и `GET /trips/:id/photos` возвращают все записи без limit/cursor.

### 13. [ensureTripAccess](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/access.go#12-38) создаёт свой контекст

[access.go](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/access.go) создаёт новый `context.WithTimeout` вместо получения из Fiber. Это значит context cancellation от HTTP-клиента не работает.

---

## 🟢 Менее критичные проблемы

### 14. Docker Compose `version: '3.8'` deprecated
### 15. DSN в [docker-compose.yml](file:///home/critiq/development/trip-listik/docker-compose.yml) — вложенные env vars не работают (строка 26)
### 16. Нет [.env.example](file:///home/critiq/development/trip-listik/frontend/.env.example) в root (есть только в frontend)
### 17. [Feed](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/feed.go#31-112) endpoint не требует auth, но `friends` filter в нём — да → неловкий UX-флоу
### 18. `token.css` radius-2xl не определён, хотя используется в trip detail page
### 19. Нет обработки `enter` key в comment input
### 20. Cover photo preview в create trip не загружается на бэкенд — теряется при step-2 transition

---

## 📋 План фикса критических и серьёзных проблем

### Фаза 1 — Безопасность (приоритет: немедленно)

| # | Что | Файл | Как |
|---|---|---|---|
| 1 | Добавить `auth_date` проверку | [auth/telegram.go](file:///home/critiq/development/trip-listik/backend/internal/auth/telegram.go) | Парсить `auth_date` из initData, отклонять если старше 5 мин |
| 2 | Фикс транзакции ApproveJoin | [handlers/members.go](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/members.go) | Передавать `tx *gorm.DB` в store-методы или создать `Store.WithTx()` |
| 3 | Санитизация ILIKE | [store/trips.go](file:///home/critiq/development/trip-listik/backend/internal/store/trips.go) | Экранировать `%`, `_`, `\` в пользовательском input |
| 4 | Уменьшить JWT TTL + refresh | [auth/jwt.go](file:///home/critiq/development/trip-listik/backend/internal/auth/jwt.go) + [handlers/auth.go](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/auth.go) | TTL=2h, добавить `/v1/auth/refresh` endpoint |

### Фаза 2 — Производительность (приоритет: высокий)

| # | Что | Файл | Как |
|---|---|---|---|
| 5 | Фильтрация в SQL для MyTrips | [store/trips.go](file:///home/critiq/development/trip-listik/backend/internal/store/trips.go) | Добавить `ListUserTripsByStatus(ctx, userID, status)` с WHERE в SQL |
| 6 | Объединить stats queries | [store/stats.go](file:///home/critiq/development/trip-listik/backend/internal/store/stats.go) | Один CTE-запрос вместо трёх |
| 7 | Оптимизировать GetTrip | [store/counts.go](file:///home/critiq/development/trip-listik/backend/internal/store/counts.go) | Batch запрос для одного trip ID или один combined query |
| 8 | Добавить пагинацию в members/photos | [handlers/members.go](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/members.go), [handlers/photos.go](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/photos.go) | Cursor + limit |
| 9 | Добавить context propagation | [handlers/access.go](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/access.go) + все handlers | Использовать `c.Context()` вместо `context.Background()` |

### Фаза 3 — Код качество (приоритет: средний)

| # | Что | Файл | Как |
|---|---|---|---|
| 10 | Удалить мёртвый код | `internal/models/`, `internal/app/`, `internal/bot/`, `internal/repository/`, `internal/services/`, `internal/state/` | Удалить полностью или мигрировать бот в новую архитектуру |
| 11 | Мигрировать на Svelte 5 runes | Все [.svelte](file:///home/critiq/development/trip-listik/frontend/src/routes/+page.svelte) файлы | `export let` → `let { x } = $props()`, [$:](file:///home/critiq/development/trip-listik/.env) → `$derived()` / `$effect()` |
| 12 | Добавить graceful SSE disconnect | [handlers/stream.go](file:///home/critiq/development/trip-listik/backend/internal/httpapi/handlers/stream.go) | Слушать `c.Context().Done()` |
| 13 | Input validation библиотека | Все handlers | Добавить `go-playground/validator` для struct tags |
| 14 | Structured logging | `server/` + handlers | Заменить `log.Printf` на `slog` |
| 15 | Добавить тесты | `store/`, `handlers/` | Минимум: store integration tests, handler unit tests |
