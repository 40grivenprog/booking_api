-- Remove unique constraint from username field
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_unique;
