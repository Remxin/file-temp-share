-- Initial schema: users, transfers, files

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'USER',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_email_key UNIQUE (email),
    CONSTRAINT users_role_check CHECK (role IN ('USER', 'ADMIN'))
);

CREATE TABLE transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    short_code VARCHAR(12) NOT NULL,
    title VARCHAR(100),
    message TEXT,
    password_hash VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    expires_at TIMESTAMPTZ NOT NULL,
    max_downloads INTEGER,
    download_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT transfers_short_code_key UNIQUE (short_code),
    CONSTRAINT transfers_status_check CHECK (status IN ('active', 'expired', 'deleted')),
    CONSTRAINT transfers_max_downloads_check CHECK (max_downloads IS NULL OR max_downloads > 0),
    CONSTRAINT transfers_download_count_check CHECK (download_count >= 0)
);

CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transfer_id UUID NOT NULL REFERENCES transfers (id) ON DELETE CASCADE,
    original_name VARCHAR(255) NOT NULL,
    storage_key VARCHAR(512) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT files_storage_key_key UNIQUE (storage_key),
    CONSTRAINT files_size_bytes_check CHECK (size_bytes >= 0)
);

CREATE INDEX idx_transfers_short_code ON transfers (short_code);
CREATE INDEX idx_transfers_expires_at ON transfers (expires_at);
CREATE INDEX idx_transfers_user_id ON transfers (user_id);
CREATE INDEX idx_transfers_status ON transfers (status);
CREATE INDEX idx_files_transfer_id ON files (transfer_id);
