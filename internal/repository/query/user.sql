-- name: CreateUser :one
INSERT INTO users (name, username, email, password, role, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetAllByRole :many
SELECT * FROM users
WHERE role = $1;

-- name: UpdateUser :exec
UPDATE users
SET name = $2, username = $3, email = $4, password = $5, role = $6, status = $7, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
