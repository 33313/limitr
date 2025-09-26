-- name: CreateKey :one
INSERT INTO api_keys (
    hashed_key, window_size_seconds, requests_per_window
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetKeyById :one
SELECT * FROM api_keys
WHERE id = $1;

-- name: GetKeyByHash :one
SELECT * FROM api_keys
WHERE hashed_key = $1;

-- name: ListKeys :many
SELECT * FROM api_keys
ORDER BY created_at DESC;

-- name: DeleteKey :exec
DELETE FROM api_keys
WHERE id = $1;
