CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hashed_key TEXT UNIQUE NOT NULL,
    limit_per_minute INT NOT NULL DEFAULT 60,
    created_at TIMESTAMP DEFAULT NOW()
);
