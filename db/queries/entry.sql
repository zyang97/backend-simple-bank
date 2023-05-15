-- name: CreateEntry :one
INSERT INTO entries (
    account_id,
    amount
    ) VALUES(
        $1, $2
    ) RETURNING *;

-- name: GetEntry :one
SELECT * from entries
WHERE id = $1
LIMIT 1;

-- name: ListEntry :many
SELECT * from entries
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateEntry :one
UPDATE entries
SET amount = $1
WHERE id = $2
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;
