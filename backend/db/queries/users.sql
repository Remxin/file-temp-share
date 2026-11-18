-- name: CreateUser :one
INSERT INTO users (email, password_hash, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2
WHERE id = $1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created_at DESC;

-- name: SetUserActive :one
UPDATE users
SET is_active = $2
WHERE id = $1
RETURNING *;

-- name: CountActiveAdmins :one
SELECT COUNT(*)::bigint AS count
FROM users
WHERE role = 'ADMIN'
  AND is_active = TRUE;
