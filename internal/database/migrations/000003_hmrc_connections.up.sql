CREATE TABLE hmrc_connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- MTD Identity
    mtd_id TEXT NOT NULL,

    -- Tokens (These will be encrypted at the app level before insertion)
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    token_expiry TIMESTAMPTZ NOT NULL,

    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Ensure one user usually links one HMRC account (for now)
    CONSTRAINT unique_user_hmrc UNIQUE (user_id)
);

-- Index for fast lookups during token refresh
CREATE INDEX idx_hmrc_connections_user_id ON hmrc_connections(user_id);
