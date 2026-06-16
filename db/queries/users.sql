-- name: CreateUser :one
INSERT INTO users (name, dob)
VALUES ($1, $2)
RETURNING id, name, dob, created_at;

-- name: GetUser :one
SELECT id, name, dob, created_at
FROM users
WHERE id = $1;

-- name: GetAllUsers :many
SELECT id, name, dob, created_at
FROM users
ORDER BY id DESC;

-- name: UpdateUser :one
UPDATE users
SET name = $1, dob = $2
WHERE id = $3
RETURNING id, name, dob, created_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;