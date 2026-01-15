-- Create the waitlist table to store email signups
CREATE TABLE waitlist (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Prevent duplicate signups for the same email
CREATE UNIQUE INDEX idx_waitlist_email ON waitlist(email);
