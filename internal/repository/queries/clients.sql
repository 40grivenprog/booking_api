-- name: CreateClient :one
INSERT INTO clients (first_name, last_name, chat_id, locale)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateClientLocale :exec
UPDATE clients
SET locale = $2
WHERE id = $1;