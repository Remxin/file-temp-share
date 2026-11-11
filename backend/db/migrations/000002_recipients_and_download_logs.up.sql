-- Recipients = notification targets for a transfer (email with public link).
-- They do NOT need an account; access is via short_code (+ optional transfer password).

CREATE TABLE transfer_recipients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transfer_id UUID NOT NULL REFERENCES transfers (id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT transfer_recipients_transfer_email_key UNIQUE (transfer_id, email)
);

CREATE INDEX idx_transfer_recipients_transfer_id ON transfer_recipients (transfer_id);

-- Download history (per transfer; optional file when per-file download is used).

CREATE TABLE download_logs (
    id BIGSERIAL PRIMARY KEY,
    transfer_id UUID NOT NULL REFERENCES transfers (id) ON DELETE CASCADE,
    file_id UUID REFERENCES files (id) ON DELETE SET NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    downloaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_download_logs_transfer_id ON download_logs (transfer_id);
CREATE INDEX idx_download_logs_downloaded_at ON download_logs (downloaded_at);
