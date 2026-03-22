# Invite Guide

This guide explains how the invite system works in TripListik, including
the Telegram notification flow.

---

## Flow overview

```
Trip owner                       Invitee
    │                                │
    │  POST /trips/:id/invite        │
    │  { username: "johndoe" }       │
    ├─────────────────────────────►  TG message: "Alex invited you to Paris 2026"
    │                                │              [View Invite] button
    │                                │
    │                          (invitee taps button, opens Mini App)
    │                                │
    │                          POST /invites/:id/respond
    │                          { action: "accept" }
    │                                │
    ◄─ TG message: "John accepted"   │  TG invite message ← deleted
    │                                │
    │                         invitee is now a trip member
```

---

## Sending an invite

Invite by **username** (resolved via user search):

```http
POST /v1/trips/{trip_id}/invite
Authorization: Bearer <token>

{ "username": "johndoe" }
```

Invite by **user_id** (if you already have it from a search result):

```http
POST /v1/trips/{trip_id}/invite
Authorization: Bearer <token>

{ "user_id": "550e8400-e29b-41d4-a716-446655440000" }
```

**Duplicate guard:** if the user already has a pending invite, the endpoint
returns `200` with the existing invite — it does NOT create a duplicate.

---

## What happens on Telegram

**When invite is sent:**
- The invitee receives a Telegram message from the bot:
  ```
  Alex invited you to join Paris 2026
  Paris · Jul 1, 2026 – Jul 10, 2026
  [View Invite]
  ```
- The `[View Invite]` button opens the Mini App on the Inbox/Notifications page.
- The `message_id` of this notification is stored in the database.

**When invitee accepts or declines:**
- The TG invite notification is **deleted** from the invitee's chat (clean inbox).
- The trip owner receives a notification:
  - Accept: `✅ John принял(а) приглашение в поездку «Paris 2026»`
  - Decline: `❌ John отклонил(а) приглашение в поездку «Paris 2026»`

**When owner cancels invite:**
- The TG invite notification is **deleted** from the invitee's chat.
- No further notification is sent.

---

## Responding to an invite (in-app)

```http
POST /v1/invites/{invite_id}/respond
Authorization: Bearer <token>

{ "action": "accept" }
```

or

```http
{ "action": "decline", "comment": "Busy that week, sorry!" }
```

---

## Cancelling an invite (owner)

```http
DELETE /v1/invites/{invite_id}
Authorization: Bearer <token>
```

Only the trip owner can cancel. Returns `204`.

---

## Finding a user to invite (search)

```http
GET /v1/users/search?q=john&limit=5
Authorization: Bearer <token>
```

**Response:**
```json
{
  "items": [
    {
      "id": "uuid",
      "username": "johndoe",
      "first_name": "John",
      "last_name": "Doe",
      "photo_url": "..."
    }
  ]
}
```

The frontend invite modal uses this endpoint with a 300ms debounce
on the username input field.
