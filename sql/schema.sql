CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hashed_key TEXT UNIQUE NOT NULL,
    window_size_seconds INT NOT NULL,
    requests_per_window INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
