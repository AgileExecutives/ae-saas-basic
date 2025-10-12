-- Create token_blacklist table
CREATE TABLE IF NOT EXISTS token_blacklist (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    token_id VARCHAR(255) NOT NULL UNIQUE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    reason VARCHAR(255)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_token_blacklist_deleted_at ON token_blacklist(deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_token_blacklist_token_id ON token_blacklist(token_id);
CREATE INDEX IF NOT EXISTS idx_token_blacklist_user_id ON token_blacklist(user_id);
CREATE INDEX IF NOT EXISTS idx_token_blacklist_expires_at ON token_blacklist(expires_at);

-- Create trigger for updated_at
CREATE TRIGGER update_token_blacklist_updated_at BEFORE UPDATE
ON token_blacklist FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create cleanup job for expired tokens (optional, for maintenance)
-- This can be run periodically to clean up expired tokens
-- DELETE FROM token_blacklist WHERE expires_at < CURRENT_TIMESTAMP;