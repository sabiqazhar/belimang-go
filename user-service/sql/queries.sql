-- name: CreateUser :one
INSERT INTO users (id, email, username, password, is_admin)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT id, email, username, password, is_admin, created_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, email, username, password, is_admin, created_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT id, email, username, password, is_admin, created_at
FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT id, email, username, is_admin, created_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET email = $2, username = $3, is_admin = $4
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: IsAdminEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND is_admin = TRUE);

-- name: CountUsers :one
SELECT COUNT(*) FROM users;