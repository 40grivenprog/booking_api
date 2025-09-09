-- name: CreateProfessional :one
INSERT INTO professionals (username, first_name, last_name, phone_number, password_hash, chat_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProfessionalByID :one
SELECT * FROM professionals
WHERE id = $1;

-- name: GetProfessionalByUsername :one
SELECT * FROM professionals
WHERE username = $1;

-- name: GetProfessionalByChatID :one
SELECT * FROM professionals
WHERE chat_id = $1;

-- name: GetProfessionals :many
SELECT * FROM professionals
WHERE chat_id is not null
ORDER BY created_at DESC;

-- name: UpdateProfessionalChatID :one
UPDATE professionals
SET chat_id = $2
WHERE id = $1
RETURNING *;

-- name: VerifyProfessionalCredentials :one
SELECT * FROM professionals
WHERE username = $1 AND password_hash = $2;
