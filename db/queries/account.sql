-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
    ) VALUES(
        $1, $2, $3
    ) RETURNING *;

-- name: GetAccount :one
SELECT * from accounts
WHERE id = $1
LIMIT 1;

-- name: ListAccount :many
SELECT * from accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $1
WHERE id = $2
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;
