-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id, import_batch_id, bank_date, description, amount, source_row_hash
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUnprocessedTransactions :many
-- Used for the "Swiping" UI. Returns a page of items.
SELECT * FROM transactions
WHERE user_id = $1
  AND status = 'unprocessed'
ORDER BY bank_date DESC, created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateTransactionCategory :one
-- The core "Swipe Right" action.
UPDATE transactions
SET status = $2,
    stream = $3,
    category = $4,
    updated_at = NOW()
WHERE id = $1 AND user_id = $5
RETURNING *;

-- name: MarkTransactionPersonal :one
-- The "Swipe Left" action. No stream/category needed.
UPDATE transactions
SET status = 'personal',
    stream = NULL,
    category = NULL,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: GetSubmissionData :many
-- Fetches business data for a SPECIFIC stream (Trade OR Property)
SELECT category, SUM(amount) as total_amount
FROM transactions
WHERE user_id = $1
  AND status = 'business'
  AND stream = $2 -- Passed in as 'TRADE' or 'PROPERTY'
  AND bank_date >= sqlc.arg(start_date)::date
  AND bank_date <= sqlc.arg(end_date)::date
GROUP BY category;

-- name: DeleteTransactionsForPeriod :exec
-- The Purge. Runs after submission.
DELETE FROM transactions
WHERE user_id = $1
  AND status = 'business'
  AND stream = $2
  AND bank_date >= sqlc.arg(start_date)::date
  AND bank_date <= sqlc.arg(end_date)::date;

-- name: GetDashboardStats :one
-- UPDATED: Now counts Trade and Property separately
SELECT
    COUNT(*) FILTER (WHERE status = 'unprocessed') as unprocessed_count,
    COUNT(*) FILTER (WHERE status = 'business' AND stream = 'TRADE') as trade_count,
    COUNT(*) FILTER (WHERE status = 'business' AND stream = 'PROPERTY') as property_count
FROM transactions
WHERE user_id = $1;
