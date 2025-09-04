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


-- name: ConfirmAppointmentWithDetails :one
WITH updated_appointment AS (
    UPDATE appointments
    SET status = 'confirmed', updated_at = NOW()
    WHERE appointments.id = $1 AND appointments.professional_id = $2
    RETURNING *
)
SELECT 
    ua.id,
    ua.type,
    ua.client_id,
    ua.professional_id,
    ua.start_time,
    ua.end_time,
    ua.status,
    ua.created_at,
    ua.updated_at,
    c.id as client_id,
    c.first_name as client_first_name,
    c.last_name as client_last_name,
    c.chat_id as client_chat_id
FROM updated_appointment ua
LEFT JOIN clients c ON c.id = ua.client_id;


-- name: GetAppointmentsByProfessionalWithStatus :many
SELECT 
    a.*,
    c.id AS client_id,
    c.first_name AS client_first_name,
    c.last_name AS client_last_name,
    c.phone_number AS client_phone_number,
    c.chat_id AS client_chat_id
FROM appointments a
LEFT JOIN clients c ON c.id = a.client_id
WHERE a.professional_id = $1
  AND a.status = $2
  AND a.start_time > NOW()
ORDER BY a.start_time DESC;

