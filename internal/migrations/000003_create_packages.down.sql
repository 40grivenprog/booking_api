-- Drop indexes
DROP INDEX IF EXISTS idx_packages_client_professional;
DROP INDEX IF EXISTS idx_packages_deactivated_at;
DROP INDEX IF EXISTS idx_packages_expires_at;
DROP INDEX IF EXISTS idx_packages_professional_id;
DROP INDEX IF EXISTS idx_packages_client_id;

-- Drop packages table
DROP TABLE IF EXISTS packages;
