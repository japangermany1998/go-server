-- name: CreateChirp :one
INSERT INTO chirp (id, created_at, updated_at, body, user_id)
VALUES (
           gen_random_uuid(), now(), now(), $1, $2
    )
RETURNING *;

-- name: DeleteChirp :one
DELETE FROM chirp
WHERE user_id = $1 AND id = $2
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirp
WHERE user_id = $1 OR $2
ORDER BY created_at;

-- name: GetChirp :one
SELECT * FROM chirp
WHERE id = $1;