-- name: CreateSession :one
INSERT INTO sessions (token, user_id, expiry)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE token = $1
AND expiry > NOW()
LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token = $1;
