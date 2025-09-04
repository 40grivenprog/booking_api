-- Create user_type enum
CREATE TYPE user_type AS ENUM ('client', 'professional');

-- Create appointment_type enum
CREATE TYPE appointment_type AS ENUM ('appointment', 'unavailable');

-- Create appointment_status enum
CREATE TYPE appointment_status AS ENUM ('pending', 'confirmed', 'cancelled', 'completed');

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id BIGINT UNIQUE, -- Telegram chat ID (NULL for clients without Telegram)
    username VARCHAR(255) NOT NULL, -- Username (required)
    first_name VARCHAR(255) NOT NULL, -- First name (required)
    last_name VARCHAR(255) NOT NULL, -- Last name (required)
    user_type user_type NOT NULL, -- User type (required)
    phone_number VARCHAR(20), -- Optional for both
    password_hash VARCHAR(255), -- For professionals only required if user_type is 'professional'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create appointments table
CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type appointment_type NOT NULL, -- Type of the appointments (required)
    client_id UUID REFERENCES users(id), -- can be null for unavailable appointments and for clients without Telegram
    professional_id UUID NOT NULL REFERENCES users(id), -- Required
    start_time TIMESTAMP WITH TIME ZONE NOT NULL, -- Required
    end_time TIMESTAMP WITH TIME ZONE NOT NULL, -- Required
    status appointment_status DEFAULT 'pending',
    cancellation_reason TEXT, -- Reason for cancellation (optional)
    cancelled_by UUID REFERENCES users(id), -- Who cancelled (optional)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_chat_id ON users(chat_id);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_user_type ON users(user_type);
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
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_appointments_updated_at BEFORE UPDATE ON appointments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
