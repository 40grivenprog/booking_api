-- name: CreateClient :one
INSERT INTO users (username, first_name, last_name, user_type, phone_number, password_hash, chat_id)
VALUES ($1, $2, $3, 'client', $4, NULL, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserByChatID :one
SELECT * FROM users
WHERE chat_id = $1;

-- name: GetProfessionals :many
SELECT * FROM users
WHERE user_type = 'professional' AND chat_id IS NOT NULL
ORDER BY created_at DESC;

-- name: UpdateUserChatID :one
UPDATE users
SET chat_id = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserByUsername :one
UPDATE users
SET chat_id = $2
WHERE username = $1
RETURNING *;

-- name: CreateProfessional :one
INSERT INTO users (username, first_name, last_name, user_type, phone_number, password_hash, chat_id)
VALUES ($1, $2, $3, 'professional', $4, $5, $6)
RETURNING *;

-- name: VerifyProfessionalCredentials :one
SELECT * FROM users
WHERE username = $1 AND user_type = 'professional' AND password_hash = $2;
