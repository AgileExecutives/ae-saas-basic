-- Create plans table
CREATE TABLE IF NOT EXISTS plans (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'EUR',
    invoice_period VARCHAR(50) NOT NULL DEFAULT 'monthly',
    max_users INTEGER DEFAULT 10,
    max_clients INTEGER DEFAULT 100,
    features JSONB,
    active BOOLEAN DEFAULT true
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_plans_deleted_at ON plans(deleted_at);
CREATE INDEX IF NOT EXISTS idx_plans_active ON plans(active);
CREATE INDEX IF NOT EXISTS idx_plans_slug ON plans(slug);

-- Create trigger for updated_at
CREATE TRIGGER update_plans_updated_at BEFORE UPDATE
ON plans FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();