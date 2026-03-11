DROP INDEX IF EXISTS idx_users_username;

ALTER TABLE trip_invites
  DROP COLUMN IF EXISTS comment;
