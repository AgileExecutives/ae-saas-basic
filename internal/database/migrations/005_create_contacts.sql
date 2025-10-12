-- Create contacts table
CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    mobile VARCHAR(50),
    street VARCHAR(255),
    zip VARCHAR(20),
    city VARCHAR(255),
    country VARCHAR(255),
    type VARCHAR(50) DEFAULT 'contact',
    notes TEXT,
    active BOOLEAN DEFAULT true
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_contacts_deleted_at ON contacts(deleted_at);
CREATE INDEX IF NOT EXISTS idx_contacts_email ON contacts(email);
CREATE INDEX IF NOT EXISTS idx_contacts_type ON contacts(type);
CREATE INDEX IF NOT EXISTS idx_contacts_active ON contacts(active);
CREATE INDEX IF NOT EXISTS idx_contacts_name ON contacts(first_name, last_name);

-- Create trigger for updated_at
CREATE TRIGGER update_contacts_updated_at BEFORE UPDATE
ON contacts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();