-- name: CreateReceipt :one
INSERT INTO receipts (
    user_id, correlation_id, receipt_timestamp, submission_type, period_key, payload_hash, totals_json
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: ListReceipts :many
SELECT * FROM receipts
WHERE user_id = $1
ORDER BY created_at DESC;
