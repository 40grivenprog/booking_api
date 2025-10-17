-- Create appointment_type enum
DO $$ BEGIN
    CREATE TYPE appointment_type AS ENUM ('appointment', 'unavailable');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create appointment_status enum
DO $$ BEGIN
    CREATE TYPE appointment_status AS ENUM ('pending', 'confirmed', 'cancelled', 'completed');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create professionals table
CREATE TABLE IF NOT EXISTS professionals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id BIGINT UNIQUE, -- Telegram chat ID (NULL for professionals created via admin)
    first_name VARCHAR(255) NOT NULL, -- First name (required)
    last_name VARCHAR(255) NOT NULL, -- Last name (required)
    phone_number VARCHAR(20), -- Optional
    username VARCHAR(255) NOT NULL UNIQUE, -- Username (required)
    password_hash VARCHAR(255), -- Required
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create clients table
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id BIGINT UNIQUE, -- Telegram chat ID (NULL for clients without Telegram)
    first_name VARCHAR(255) NOT NULL, -- First name (required)
    last_name VARCHAR(255) NOT NULL, -- Last name (required)
    phone_number VARCHAR(20), -- Optional
    created_by UUID REFERENCES professionals(id), -- Optional ID of the professional who created the client
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create appointments table
CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type appointment_type NOT NULL, -- Type of the appointments (required)
    client_id UUID REFERENCES clients(id), -- can be null for unavailable appointments and for clients without Telegram
    professional_id UUID NOT NULL REFERENCES professionals(id), -- Required
    start_time TIMESTAMP WITH TIME ZONE NOT NULL, -- Required
    end_time TIMESTAMP WITH TIME ZONE NOT NULL, -- Required
    status appointment_status DEFAULT 'pending',
    cancellation_reason TEXT, -- Reason for cancellation (optional)
    cancelled_by_professional_id UUID REFERENCES professionals(id), -- Who cancelled (optional)
    cancelled_by_client_id UUID REFERENCES clients(id), -- Who cancelled (optional)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_professionals_chat_id ON professionals(chat_id);
CREATE INDEX IF NOT EXISTS idx_professionals_username ON professionals(username);
CREATE INDEX IF NOT EXISTS idx_clients_chat_id ON clients(chat_id);
CREATE INDEX IF NOT EXISTS idx_clients_created_by ON clients(created_by);
CREATE INDEX IF NOT EXISTS idx_appointments_client_id ON appointments(client_id);
CREATE INDEX IF NOT EXISTS idx_appointments_professional_id ON appointments(professional_id);
CREATE INDEX IF NOT EXISTS idx_appointments_start_time ON appointments(start_time);
CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments(status);
CREATE INDEX IF NOT EXISTS idx_appointments_type ON appointments(type);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_professionals_updated_at BEFORE UPDATE ON professionals FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_clients_updated_at BEFORE UPDATE ON clients FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_appointments_updated_at BEFORE UPDATE ON appointments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
