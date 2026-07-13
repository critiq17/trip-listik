-- invite_links: shareable multi-use trip invite tokens.
-- Sent into any Telegram chat as t.me/<bot>?startapp=join_<token>;
-- the recipient does not need to be a bot user yet.
CREATE TABLE IF NOT EXISTS invite_links (
    token      TEXT PRIMARY KEY,
    trip_id    UUID NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_invite_links_trip_creator ON invite_links(trip_id, created_by);
