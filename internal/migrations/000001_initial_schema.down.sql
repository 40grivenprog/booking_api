-- Drop triggers
DROP TRIGGER IF EXISTS update_appointments_updated_at ON appointments;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_appointments_type;
DROP INDEX IF EXISTS idx_appointments_status;
DROP INDEX IF EXISTS idx_appointments_start_time;
DROP INDEX IF EXISTS idx_appointments_professional_id;
DROP INDEX IF EXISTS idx_appointments_client_id;
DROP INDEX IF EXISTS idx_users_user_type;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_chat_id;

-- Drop tables
DROP TABLE IF EXISTS appointments;
DROP TABLE IF EXISTS users;

-- Drop enums
DROP TYPE IF EXISTS appointment_status;
DROP TYPE IF EXISTS appointment_type;
DROP TYPE IF EXISTS user_type;
