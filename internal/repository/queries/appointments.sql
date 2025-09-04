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
LEFT JOIN clients c ON c.id = na.client_id
LEFT JOIN professionals p ON p.id = na.professional_id;

-- name: GetAppointmentByID :one
SELECT * FROM appointments
WHERE id = $1;

-- name: GetAppointmentsByClient :many
SELECT * FROM appointments
WHERE client_id = $1
ORDER BY start_time DESC;

-- name: GetAppointmentsByProfessional :many
SELECT * FROM appointments
WHERE professional_id = $1
ORDER BY start_time DESC;

-- name: UpdateAppointmentStatus :one
UPDATE appointments
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;
