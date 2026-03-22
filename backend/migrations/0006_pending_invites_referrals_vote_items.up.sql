-- pending_invites: for users not yet registered in the bot
CREATE TABLE IF NOT EXISTS pending_invites (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id     UUID NOT NULL,
    telegram_id BIGINT NOT NULL,
    invited_by  UUID NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pending_invites_telegram_id ON pending_invites(telegram_id);

-- referrals: track who referred whom via the profile share link
CREATE TABLE IF NOT EXISTS referrals (
    referrer_id UUID NOT NULL,
    referred_id UUID NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (referrer_id, referred_id)
);

-- vote_items: hotels, airlines, activities that trip members can vote on
CREATE TABLE IF NOT EXISTS vote_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id     UUID NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    category    TEXT NOT NULL DEFAULT 'OTHER',
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    image_url   TEXT NOT NULL DEFAULT '',
    added_by    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_vote_items_trip_id ON vote_items(trip_id);

-- vote_item_votes: each user can vote once per item
CREATE TABLE IF NOT EXISTS vote_item_votes (
    item_id    UUID NOT NULL REFERENCES vote_items(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (item_id, user_id)
);

-- add comment column to trip_invites if not exists
ALTER TABLE trip_invites ADD COLUMN IF NOT EXISTS comment TEXT;
