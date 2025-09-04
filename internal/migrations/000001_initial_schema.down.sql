-- Drop triggers
DROP TRIGGER IF EXISTS update_appointments_updated_at ON appointments;
DROP TRIGGER IF EXISTS update_clients_updated_at ON clients;
DROP TRIGGER IF EXISTS update_professionals_updated_at ON professionals;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_appointments_type;
DROP INDEX IF EXISTS idx_appointments_status;
DROP INDEX IF EXISTS idx_appointments_start_time;
DROP INDEX IF EXISTS idx_appointments_professional_id;
DROP INDEX IF EXISTS idx_appointments_client_id;
DROP INDEX IF EXISTS idx_clients_created_by;
DROP INDEX IF EXISTS idx_clients_chat_id;
DROP INDEX IF EXISTS idx_professionals_username;
DROP INDEX IF EXISTS idx_professionals_chat_id;

-- Drop tables
DROP TABLE IF EXISTS appointments;
DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS professionals;

-- Drop enums
DROP TYPE IF EXISTS appointment_status;
DROP TYPE IF EXISTS appointment_type;
