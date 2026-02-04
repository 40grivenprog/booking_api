-- Drop triggers and functions
DROP TRIGGER IF EXISTS trigger_delete_appointments_on_client_delete ON clients;
DROP FUNCTION IF EXISTS delete_appointments_on_client_delete();

-- Drop indexes
DROP INDEX IF EXISTS idx_subscriptions_professional_id;
DROP INDEX IF EXISTS idx_subscriptions_client_id;
DROP INDEX IF EXISTS idx_client_appointments_appointment_id;
DROP INDEX IF EXISTS idx_client_appointments_client_id;
DROP INDEX IF EXISTS idx_appointments_status;
DROP INDEX IF EXISTS idx_appointments_start_time;
DROP INDEX IF EXISTS idx_appointments_professional_id;

-- Drop tables
DROP TABLE IF EXISTS client_appointments;
DROP TABLE IF EXISTS appointments;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS professionals;
DROP TABLE IF EXISTS clients;

-- Drop enum types
DROP TYPE IF EXISTS appointment_status;
DROP TYPE IF EXISTS appointment_type;
