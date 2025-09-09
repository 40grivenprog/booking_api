-- Initialize the booking database
-- This script runs when the PostgreSQL container starts for the first time

-- Set timezone to Europe/Berlin
SET timezone = 'Europe/Berlin';

-- Create extensions if they don't exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- The database and user are already created by the environment variables
-- in docker-compose.yml, so we just need to ensure extensions are available
