DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS user_wishlist_items;
ALTER TABLE users DROP COLUMN IF EXISTS is_public;
ALTER TABLE users DROP COLUMN IF EXISTS bio;

-- Note: PostgreSQL does not support removing values from an ENUM type using ALTER TYPE natively.
-- The newly added 'tour' and 'group' values for trip_visibility will remain.
