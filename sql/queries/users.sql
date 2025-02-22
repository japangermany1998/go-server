-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(), now(), now(), $1, $2
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET updated_at = now(), email = $1, hashed_password = $2
WHERE id = $3
RETURNING *;

-- name: UpdateUserToChirpyRed :one
UPDATE users SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;
