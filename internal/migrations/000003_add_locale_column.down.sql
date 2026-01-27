-- Remove locale column from professionals table
ALTER TABLE professionals DROP COLUMN locale;

-- Remove locale column from clients table
ALTER TABLE clients DROP COLUMN locale;
