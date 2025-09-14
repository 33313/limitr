-- name: CreateKey :one
INSERT INTO api_keys (
    hashed_key, limit_per_minute
) VALUES (
    $1, $2
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
