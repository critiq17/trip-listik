ALTER TABLE trip_invites
  ADD COLUMN IF NOT EXISTS comment text;

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
