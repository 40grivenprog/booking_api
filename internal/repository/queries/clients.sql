-- name: CreateClient :one
INSERT INTO clients (first_name, last_name, phone_number, chat_id, created_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetClientByID :one
SELECT * FROM clients
WHERE id = $1;

-- name: GetClientByChatID :one
SELECT * FROM clients
WHERE chat_id = $1;

-- name: GetClients :many
SELECT * FROM clients
ORDER BY created_at DESC;

-- name: UpdateClientChatID :one
UPDATE clients
SET chat_id = $2
WHERE id = $1
RETURNING *;
