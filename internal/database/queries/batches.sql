-- name: CreateImportBatch :one
INSERT INTO transaction_import_batches (
    user_id, source_filename
) VALUES (
    $1, $2
) RETURNING *;
