-- name: CreateWaitlistEntry :one
INSERT INTO waitlist (
    email
) VALUES (
    $1
) RETURNING id, email, created_at;

-- name: CheckWaitlistEmail :one
SELECT EXISTS(SELECT 1 FROM waitlist WHERE email = $1);
