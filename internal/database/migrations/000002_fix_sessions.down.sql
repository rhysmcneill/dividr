-- 1. Remove the Index we added
DROP INDEX IF EXISTS sessions_expiry_idx;

-- 2. Remove the 'data' column we added
ALTER TABLE sessions DROP COLUMN data;

-- 3. TRUNCATE the table to clear data.
-- We MUST do this because we are about to add a NOT NULL column (user_id).
-- If we don't clear the table, this migration will fail on any existing rows.
TRUNCATE TABLE sessions;

-- 4. Add the 'user_id' column back (Restoring original state)
ALTER TABLE sessions
    ADD COLUMN user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE;
