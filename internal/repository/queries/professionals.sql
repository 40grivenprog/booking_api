-- name: CreateProfessional :one
INSERT INTO professionals (username, first_name, last_name, phone_number, password_hash, chat_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProfessionalByUsername :one
SELECT * FROM professionals
WHERE username = $1;

-- name: GetProfessionals :many
SELECT * FROM professionals
WHERE chat_id is not null
ORDER BY created_at DESC;

-- name: UpdateProfessionalChatID :one
UPDATE professionals
SET chat_id = $2
WHERE id = $1
RETURNING *;

-- name: GetAppointmentsByProfessionalWithStatusAndDate :many
SELECT 
    a.id,
    a.type,
    a.start_time,
    a.end_time,
    a.description,
    a.status,
    a.created_at,
    a.updated_at,
    a.client_id,
    c.first_name as client_first_name,
    c.last_name as client_last_name,
    c.phone_number as client_phone_number
FROM appointments a
LEFT JOIN clients c ON a.client_id = c.id
WHERE a.professional_id = $1
    AND ($2 = '' OR a.status = $2::appointment_status)
    AND ($3 = '' OR DATE(a.start_time) = $3::date)
ORDER BY a.start_time ASC;
