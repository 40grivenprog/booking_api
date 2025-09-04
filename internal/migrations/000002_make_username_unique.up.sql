-- Add unique constraint to username field
ALTER TABLE users ADD CONSTRAINT users_username_unique UNIQUE (username);
