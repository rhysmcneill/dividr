-- name: CreateAuditEvent :one
INSERT INTO audit_events (
    user_id,
    event_type,
    ip_address,
    details
) VALUES (
    $1, $2, $3, $4
) RETURNING id, created_at;
