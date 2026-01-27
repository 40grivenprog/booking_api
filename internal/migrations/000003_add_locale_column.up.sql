-- Add locale column to professionals table
ALTER TABLE professionals ADD COLUMN locale VARCHAR(10);

-- Add locale column to clients table
ALTER TABLE clients ADD COLUMN locale VARCHAR(10);
