# Trips & Invites API

All endpoints are under `/v1`. Authenticated endpoints require `Authorization: Bearer <token>`.

---

## Trips

### Create trip

```
POST /v1/trips
Authorization: Bearer <token>
```

**Body:**
```json
{
  "title": "Paris 2026",
  "description": "Our summer trip.",
  "start_date": "2026-07-01",
  "end_date": "2026-07-10",
  "visibility": "public",
  "status": "draft",
  "country_code": "FR",
  "city": "Paris",
  "cover_photo_url": ""
}
```

| Field | Type | Required | Constraints |
|---|---|---|---|
| `title` | string | yes | 1–100 chars |
| `description` | string | no | max 2000 chars |
| `start_date` | string | no | `YYYY-MM-DD` |
| `end_date` | string | no | `YYYY-MM-DD`, ≥ start_date |
| `visibility` | string | no | `public` (default) \| `private` \| `group` |
| `status` | string | no | `draft` (default) \| `planned` \| `completed` \| `canceled` |

**Response `201`:** Trip object.

---

### Get trip

```
GET /v1/trips/:id
```

Public endpoint. Returns 403 for private trips if caller is not a member.

**Response `200`:**
```json
{
  "trip": { ... },
  "member_count": 4,
  "vote_count": 12,
  "vote_average": 4.2,
  "comment_count": 7,
  "photo_count": 3,
  "viewer_is_member": true
}
```

---

### Update trip

```
PATCH /v1/trips/:id
Authorization: Bearer <token>
```

Owner only. All fields are optional (partial update).

**Body:** same fields as create, all optional (use `null` to skip a field).

**Response `200`:** Updated trip object.

---

### Delete trip

```
DELETE /v1/trips/:id
Authorization: Bearer <token>
```

Owner only. Soft-deletes the trip (sets `deleted_at`).

**Response `204` No Content.**

---

### List my trips

```
GET /v1/trips?scope=mine&status=upcoming
Authorization: Bearer <token>
```

| Query param | Values | Description |
|---|---|---|
| `scope` | `mine` | Required |
| `status` | `upcoming` \| `past` \| `draft` \| `` (all) | Filter |

**Response `200`:**
```json
{ "items": [ { ...trip }, ... ] }
```

---

## Invites

### Send invite

```
POST /v1/trips/:id/invite
Authorization: Bearer <token>
```

Owner (or member of group trip) can invite by `username` or `user_id`.
A Telegram notification is sent to the invitee and the message_id is stored
so it can be deleted when the invite is acted upon.

**Body:**
```json
{ "username": "johndoe" }
```
or
```json
{ "user_id": "uuid" }
```

**Response `200`:**
```json
{ "invite_id": "uuid", "status": "pending" }
```
If a pending invite already exists, returns the existing one (idempotent).

---

### List trip invites

```
GET /v1/trips/:id/invites
Authorization: Bearer <token>
```

Owner or group member only.

**Response `200`:**
```json
{ "items": [ { ...invite with trip + inviter info }, ... ] }
```

---

### Get invite

```
GET /v1/invites/:id
Authorization: Bearer <token>
```

**Response `200`:**
```json
{ "invite": { ...invite view } }
```

---

### Respond to invite

```
POST /v1/invites/:id/respond
Authorization: Bearer <token>
```

Invited user only.

**Body:**
```json
{ "action": "accept", "comment": "Can't wait!" }
```

| `action` | Effect |
|---|---|
| `accept` | Status → `accepted`, user added as trip member, TG invite message deleted, owner notified |
| `decline` | Status → `declined`, TG invite message deleted, owner notified |

**Response `200`:**
```json
{ "status": "accepted" }
```

---

### Cancel invite

```
DELETE /v1/invites/:id
Authorization: Bearer <token>
```

Trip owner only. Sets status to `cancelled` and deletes the TG invite message from the invitee's chat.

**Response `204` No Content.**

---

## Access rights summary

| Action | Who can do it |
|---|---|
| Create trip | Any authenticated user |
| View public trip | Anyone |
| View private trip | Trip members only |
| Edit trip | Trip owner only |
| Delete trip | Trip owner only |
| Send invite | Trip owner; or any member (group trips only) |
| Accept/decline invite | The invited user only |
| Cancel invite | Trip owner only |
| View invite list | Trip owner; or any member (group trips only) |
