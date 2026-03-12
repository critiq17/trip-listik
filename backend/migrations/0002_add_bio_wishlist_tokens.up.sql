ALTER TABLE users ADD COLUMN bio text;
ALTER TABLE users ADD COLUMN is_public boolean DEFAULT true;

-- Add new trip types
ALTER TYPE trip_visibility ADD VALUE IF NOT EXISTS 'tour';
ALTER TYPE trip_visibility ADD VALUE IF NOT EXISTS 'group';

CREATE TABLE user_wishlist_items (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    country_code text NOT NULL,
    city text,
    note text,
    created_at timestamptz DEFAULT now()
);

CREATE INDEX idx_user_wishlist_items_user_id ON user_wishlist_items(user_id);

CREATE TABLE refresh_tokens (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    token_hash text NOT NULL UNIQUE,
    expires_at timestamptz NOT NULL,
    created_at timestamptz DEFAULT now(),
    last_used_at timestamptz,
    revoked_at timestamptz
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
