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
WHERE user_type = 'professional'
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

-- name: CreateAppointment :one
INSERT INTO appointments (type, client_id, professional_id, start_time, end_time, status)
VALUES ('appointment', $1, $2, $3, $4, 'pending')
RETURNING *;

-- name: CreateAppointmentWithDetails :one
WITH new_appointment AS (
    INSERT INTO appointments (type, client_id, professional_id, start_time, end_time, status)
    VALUES ('appointment', $1, $2, $3, $4, 'pending')
    RETURNING *
)
SELECT 
    na.*,
    c.id as client_id_full,
    c.username as client_username,
    c.first_name as client_first_name,
    c.last_name as client_last_name,
    c.phone_number as client_phone_number,
    c.chat_id as client_chat_id,
    c.created_at as client_created_at,
    c.updated_at as client_updated_at,
    p.id as professional_id_full,
    p.username as professional_username,
    p.first_name as professional_first_name,
    p.last_name as professional_last_name,
    p.phone_number as professional_phone_number,
    p.chat_id as professional_chat_id,
    p.created_at as professional_created_at,
    p.updated_at as professional_updated_at
FROM new_appointment na
LEFT JOIN users c ON c.id = na.client_id
LEFT JOIN users p ON p.id = na.professional_id;
