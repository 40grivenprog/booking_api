-- Drop indexes
DROP INDEX IF EXISTS idx_invites_start_time;
DROP INDEX IF EXISTS idx_invites_client_id;
DROP INDEX IF EXISTS idx_invites_appointment_id;
DROP INDEX IF EXISTS idx_invites_appointment_id_client_id;

-- Drop table
DROP TABLE IF EXISTS invites;
