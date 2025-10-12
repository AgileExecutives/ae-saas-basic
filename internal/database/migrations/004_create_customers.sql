-- Create customers table
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    street VARCHAR(255),
    zip VARCHAR(20),
    city VARCHAR(255),
    country VARCHAR(255),
    tax_id VARCHAR(255),
    vat VARCHAR(255),
    plan_id INTEGER NOT NULL REFERENCES plans(id),
    organization_id INTEGER NOT NULL REFERENCES organizations(id),
    status VARCHAR(50) DEFAULT 'active',
    payment_method VARCHAR(255),
    active BOOLEAN DEFAULT true
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_customers_deleted_at ON customers(deleted_at);
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
CREATE INDEX IF NOT EXISTS idx_customers_plan_id ON customers(plan_id);
CREATE INDEX IF NOT EXISTS idx_customers_organization_id ON customers(organization_id);
CREATE INDEX IF NOT EXISTS idx_customers_active ON customers(active);
CREATE INDEX IF NOT EXISTS idx_customers_status ON customers(status);

-- Create trigger for updated_at
CREATE TRIGGER update_customers_updated_at BEFORE UPDATE
ON customers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();