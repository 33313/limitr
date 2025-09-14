-- name: CreateKey :one
INSERT INTO api_keys (
    hashed_key, limit_per_minute
) VALUES (
    $1, $2
)
RETURNING *;

-- name: ListKeys :many
SELECT * FROM api_keys
ORDER BY created_at DESC;

-- name: DeleteKey :exec
DELETE FROM api_keys
WHERE id = $1;
