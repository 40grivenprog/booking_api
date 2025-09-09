-- name: CreateClient :one
INSERT INTO clients (first_name, last_name, phone_number, chat_id, created_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

