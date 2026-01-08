CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. USERS
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. SESSIONS
CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expiry TIMESTAMPTZ NOT NULL
);

-- 3. OAUTH_TOKENS (HMRC Access)
CREATE TABLE oauth_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL DEFAULT 'hmrc',
    access_token TEXT NOT NULL,  -- Application layer should encrypt this
    refresh_token TEXT NOT NULL, -- Application layer should encrypt this
    expiry TIMESTAMPTZ NOT NULL,
    scope TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 4. MAPPING_PROFILES (CSV Parsing Rules)
CREATE TABLE mapping_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    fingerprint TEXT NOT NULL,
    mapping_rules JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5. TRANSACTION_IMPORT_BATCHES (Upload History)
CREATE TABLE transaction_import_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_filename TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 6. TRANSACTIONS (Staging Workspace)
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    import_batch_id UUID REFERENCES transaction_import_batches(id) ON DELETE SET NULL,

    -- Source Data
    bank_date DATE NOT NULL,
    description TEXT NOT NULL,
    amount NUMERIC(14,2) NOT NULL,

    -- State
    status TEXT NOT NULL DEFAULT 'unprocessed',

    -- Stream & Category (The Logical Divide)
    -- Stream is nullable until the user swipes "Business".
    -- Then it must be 'TRADE' or 'PROPERTY'.
    stream TEXT,
    category TEXT, -- e.g. 'expenses.rent'

    -- Dedupe
    source_row_hash TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 7. RECEIPTS (Permanent Evidence)
CREATE TABLE receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),

    correlation_id TEXT NOT NULL,
    receipt_timestamp TIMESTAMPTZ NOT NULL,
    submission_type TEXT NOT NULL, -- 'quarterly_trade', 'final', etc.
    period_key TEXT NOT NULL,

    -- GDPR-safe evidence
    payload_hash TEXT NOT NULL,
    totals_json JSONB NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 8. AUDIT_EVENTS (Security Log)
CREATE TABLE audit_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id), -- Nullable for failed logins
    event_type TEXT NOT NULL,
    details JSONB,
    ip_address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- INDEXES
CREATE INDEX idx_transactions_user_status_created ON transactions (user_id, status, created_at DESC);
CREATE INDEX idx_transactions_submission ON transactions (user_id, status, stream, bank_date);
CREATE UNIQUE INDEX uniq_transactions_user_rowhash ON transactions (user_id, source_row_hash) WHERE source_row_hash IS NOT NULL;
