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

-- name: CreateTransferRecipient :one
INSERT INTO transfer_recipients (transfer_id, email)
VALUES ($1, $2)
RETURNING *;
