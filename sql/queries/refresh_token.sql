-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, revoked_at, user_id)
VALUES ($1, now(), now(), $2, NULL, $3)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT user_id FROM refresh_tokens
WHERE token = $1 and revoked_at IS NULL AND expires_at > now();

-- name: UpdateByRevokeToken :exec
UPDATE refresh_tokens SET updated_at = now(), revoked_at = now()
WHERE token = $1 and revoked_at IS NULL AND expires_at > now();