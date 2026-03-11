CREATE TABLE IF NOT EXISTS user_wishlist (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  country_code text NOT NULL,
  city text,
  note text,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE(user_id, country_code)
);

ALTER TABLE users
  ADD COLUMN IF NOT EXISTS bio text DEFAULT '',
  ADD COLUMN IF NOT EXISTS is_public boolean DEFAULT true;
