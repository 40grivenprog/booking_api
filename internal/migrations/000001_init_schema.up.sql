-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types (with existence check)
DO $$ BEGIN
    CREATE TYPE appointment_type AS ENUM ('personal', 'unavailable', 'split', 'group');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE appointment_status AS ENUM ('pending', 'confirmed', 'cancelled');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create clients table
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id BIGINT UNIQUE, -- Telegram chat ID (NULL for clients without Telegram)
    first_name VARCHAR(255) NOT NULL, -- First name (required)
    last_name VARCHAR(255) NOT NULL, -- Last name (required)
    locale VARCHAR(255) NOT NULL DEFAULT 'en', -- Locale (required)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create professionals table
CREATE TABLE IF NOT EXISTS professionals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id BIGINT UNIQUE, -- Telegram chat ID (NULL for professionals created via admin)
    first_name VARCHAR(255) NOT NULL, -- First name (required)
    last_name VARCHAR(255) NOT NULL, -- Last name (required)
    username VARCHAR(255) NOT NULL UNIQUE, -- Username (required) @username
    password_hash VARCHAR(255) NOT NULL, -- Required
    locale VARCHAR(255) NOT NULL DEFAULT 'en', -- Locale (required)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create subscriptions table
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    professional_id UUID NOT NULL REFERENCES professionals(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(client_id, professional_id)
);

-- Create appointments table
CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type appointment_type NOT NULL, -- Type of the appointments (required) - personal, unavailable, split, group
    professional_id UUID NOT NULL REFERENCES professionals(id) ON DELETE CASCADE, -- Required, cascade delete when professional is deleted
    start_time TIMESTAMP WITH TIME ZONE NOT NULL, -- Required
    end_time TIMESTAMP WITH TIME ZONE NOT NULL, -- Required
    status appointment_status DEFAULT 'pending', -- 'pending', 'confirmed'
    description TEXT, -- Description of the appointment
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create client_appointments table
CREATE TABLE IF NOT EXISTS client_appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE, -- Cascade delete when client is deleted
    appointment_id UUID NOT NULL REFERENCES appointments(id) ON DELETE CASCADE, -- Cascade delete when appointment is deleted
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(client_id, appointment_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_appointments_professional_id ON appointments(professional_id);
CREATE INDEX IF NOT EXISTS idx_appointments_start_time ON appointments(start_time);
CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments(status);
CREATE INDEX IF NOT EXISTS idx_client_appointments_client_id ON client_appointments(client_id);
CREATE INDEX IF NOT EXISTS idx_client_appointments_appointment_id ON client_appointments(appointment_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_client_id ON subscriptions(client_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_professional_id ON subscriptions(professional_id);

-- Create function to delete appointments when client is deleted
-- This is needed because appointments don't have a direct FK to clients
-- Appointments are linked to clients through client_appointments table
CREATE OR REPLACE FUNCTION delete_appointments_on_client_delete()
RETURNS TRIGGER AS $$
DECLARE
    appointment_ids_to_delete UUID[];
BEGIN
    -- Collect appointment IDs that are linked only to the deleted client
    -- We need to check this BEFORE client_appointments are deleted by CASCADE
    SELECT ARRAY_AGG(appointment_id)
    INTO appointment_ids_to_delete
    FROM client_appointments
    WHERE client_id = OLD.id
    AND appointment_id NOT IN (
        SELECT DISTINCT appointment_id
        FROM client_appointments
        WHERE client_id != OLD.id
    );
    
    -- Delete appointments that are linked only to this client
    -- (for group appointments with other clients, only the link is removed via CASCADE)
    IF appointment_ids_to_delete IS NOT NULL THEN
        DELETE FROM appointments
        WHERE id = ANY(appointment_ids_to_delete);
    END IF;
    
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to delete appointments when client is deleted
-- Use BEFORE DELETE to check client_appointments before CASCADE deletes them
DROP TRIGGER IF EXISTS trigger_delete_appointments_on_client_delete ON clients;
CREATE TRIGGER trigger_delete_appointments_on_client_delete
BEFORE DELETE ON clients
FOR EACH ROW
EXECUTE FUNCTION delete_appointments_on_client_delete();
