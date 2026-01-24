-- name: UpsertHMRCConnection :one
INSERT INTO hmrc_connections (user_id, mtd_id, access_token, refresh_token, token_expiry, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW())
ON CONFLICT (user_id) DO UPDATE
SET mtd_id = EXCLUDED.mtd_id,
    access_token = EXCLUDED.access_token,
    refresh_token = EXCLUDED.refresh_token,
    token_expiry = EXCLUDED.token_expiry,
    updated_at = NOW()
RETURNING id, user_id, mtd_id, token_expiry;

-- name: GetHMRCConnectionByUserID :one
SELECT * FROM hmrc_connections
WHERE user_id = $1
LIMIT 1;

-- name: DeleteHMRCConnection :exec
DELETE FROM hmrc_connections
WHERE user_id = $1;
