-- name: CreateFile :one
INSERT INTO files (
    transfer_id,
    original_name,
    storage_key,
    mime_type,
    size_bytes
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: ListFilesByTransferID :many
SELECT * FROM files
WHERE transfer_id = $1
ORDER BY created_at ASC;
