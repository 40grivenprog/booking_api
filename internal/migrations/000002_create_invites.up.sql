-- Create invites table
-- This table tracks client appointments with their start times
CREATE TABLE IF NOT EXISTS invites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL REFERENCES appointments(id) ON DELETE CASCADE, -- Cascade delete when appointment is deleted
    start_time TIMESTAMP WITH TIME ZONE NOT NULL, -- Start time of the appointment
    end_time TIMESTAMP WITH TIME ZONE NOT NULL, -- End time of the appointment
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE, -- Cascade delete when client is deleted (but appointment remains)
    description TEXT NOT NULL, -- Description of the appointment
    type appointment_type NOT NULL, -- Type of the appointment (required) - personal, unavailable, split, group
    professional_name VARCHAR(255) NOT NULL, -- Name of the professional
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_invites_appointment_id ON invites(appointment_id);
CREATE INDEX IF NOT EXISTS idx_invites_client_id ON invites(client_id);
CREATE INDEX IF NOT EXISTS idx_invites_start_time ON invites(start_time);
CREATE UNIQUE INDEX IF NOT EXISTS idx_invites_appointment_id_client_id ON invites(appointment_id, client_id);
