-- name: CreateTransfer :one
INSERT INTO transfers (
    user_id,
    short_code,
    title,
    message,
    password_hash,
    status,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, 'active', $6
) RETURNING *;

-- name: GetTransferByID :one
SELECT * FROM transfers
WHERE id = $1
LIMIT 1;

-- name: GetTransferByShortCode :one
SELECT * FROM transfers
WHERE short_code = $1
LIMIT 1;

-- name: ListTransfersWithStatsByUserID :many
SELECT
    t.id,
    t.user_id,
    t.short_code,
    t.title,
    t.message,
    t.password_hash,
    t.status,
    t.expires_at,
    t.max_downloads,
    t.download_count,
    t.created_at,
    (SELECT COUNT(*) FROM files f WHERE f.transfer_id = t.id)::bigint AS file_count,
    (SELECT COALESCE(SUM(f.size_bytes), 0) FROM files f WHERE f.transfer_id = t.id)::bigint AS total_size_bytes,
    (SELECT COUNT(*) FROM transfer_recipients tr WHERE tr.transfer_id = t.id)::bigint AS recipient_count,
    COALESCE(
        (SELECT f.original_name FROM files f WHERE f.transfer_id = t.id ORDER BY f.created_at ASC LIMIT 1),
        ''
    )::text AS first_file_name
FROM transfers t
WHERE t.user_id = $1
  AND t.status <> 'deleted'
ORDER BY t.created_at DESC;

-- name: CountRecipientsByTransferID :one
SELECT COUNT(*)::bigint AS count
FROM transfer_recipients
WHERE transfer_id = $1;

-- name: IncrementTransferDownloadCount :one
UPDATE transfers
SET download_count = download_count + 1
WHERE id = $1
RETURNING *;

-- name: CreateDownloadLog :exec
INSERT INTO download_logs (transfer_id, file_id, ip_address, user_agent)
VALUES ($1, $2, $3, $4);

-- name: CreateTransferRecipient :one
INSERT INTO transfer_recipients (transfer_id, email)
VALUES ($1, $2)
RETURNING *;
