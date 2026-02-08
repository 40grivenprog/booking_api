-- Create packages table
-- This table tracks appointment packages issued to clients by professionals
CREATE TABLE IF NOT EXISTS packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE, -- Cascade delete when client is deleted
    professional_id UUID NOT NULL REFERENCES professionals(id) ON DELETE CASCADE, -- Cascade delete when professional is deleted
    issued_at TIMESTAMP WITH TIME ZONE NOT NULL, -- When the package was issued (required)
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL, -- When the package expires (required)
    apppointments_number INTEGER NOT NULL, -- Number of appointments in the package (required) - Note: keeping original spelling
    deactivated_at TIMESTAMP WITH TIME ZONE, -- When the package was deactivated (optional)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_packages_client_id ON packages(client_id);
CREATE INDEX IF NOT EXISTS idx_packages_professional_id ON packages(professional_id);
CREATE INDEX IF NOT EXISTS idx_packages_expires_at ON packages(expires_at);
CREATE INDEX IF NOT EXISTS idx_packages_deactivated_at ON packages(deactivated_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_packages_client_professional_issued_at ON packages(client_id, professional_id, issued_at);
