-- name: CreateClient :one
INSERT INTO clients (first_name, last_name, phone_number, chat_id, created_by, locale)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateClientLocale :exec
UPDATE clients
SET locale = $2
WHERE id = $1;