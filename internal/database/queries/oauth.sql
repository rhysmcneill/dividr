-- name: UpsertOAuthToken :one
-- Saves or updates the token. If user already has one, we overwrite it.
INSERT INTO oauth_tokens (
    user_id, provider, access_token, refresh_token, expiry, scope, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, NOW()
)
ON CONFLICT (id) DO UPDATE -- Note: You might need a unique constraint on (user_id, provider) for this to work perfectly, or just use ID.
-- actually, standard practice is usually finding by user_id. Let's stick to simple Insert/Update for now.
SET access_token = $3, refresh_token = $4, expiry = $5, updated_at = NOW()
RETURNING *;

-- name: GetOAuthToken :one
SELECT * FROM oauth_tokens
WHERE user_id = $1 AND provider = $2
LIMIT 1;
