-- Remove the strict user link from sessions
ALTER TABLE sessions DROP COLUMN user_id;

-- Add the data storage column required by SCS
ALTER TABLE sessions ADD COLUMN data BYTEA NOT NULL;

-- Add index for cleaning up old sessions efficiently
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
