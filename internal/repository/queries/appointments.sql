-- name: GetAppointmentByID :one
SELECT * FROM appointments
WHERE appointments.id = $1;

-- name: CreateAppointmentWithDetails :one
WITH new_appointment AS (
    INSERT INTO appointments (type, client_id, professional_id, start_time, end_time, status, description)
    VALUES ('appointment', $1, $2, $3, $4, 'pending', $5)
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
    c.chat_id as client_chat_id,
    p.id as professional_id_full,
    p.username as professional_username,
    p.first_name as professional_first_name,
    p.last_name as professional_last_name
FROM updated_appointment ua
LEFT JOIN clients c ON c.id = ua.client_id
LEFT JOIN professionals p ON p.id = ua.professional_id;


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
  AND a.type = 'appointment'
ORDER BY a.start_time DESC;


-- name: CancelAppointmentByProfessionalWithDetails :one
WITH updated_appointment AS (
    UPDATE appointments
    SET 
        status = 'cancelled',
        cancellation_reason = $3,
        cancelled_by_professional_id = $2,
        updated_at = NOW()
    WHERE appointments.id = $1 
    AND professional_id = $2
    AND status IN ('pending', 'confirmed')
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
    ua.cancellation_reason,
    ua.cancelled_by_professional_id,
    ua.cancelled_by_client_id,
    ua.created_at,
    ua.updated_at,
    c.id as client_id_full,
    c.first_name as client_first_name,
    c.last_name as client_last_name,
    c.phone_number as client_phone_number,
    c.chat_id as client_chat_id,
    p.id as professional_id_full,
    p.username as professional_username,
    p.first_name as professional_first_name,
    p.last_name as professional_last_name,
    p.phone_number as professional_phone_number,
    p.chat_id as professional_chat_id
FROM updated_appointment ua
LEFT JOIN clients c ON c.id = ua.client_id
LEFT JOIN professionals p ON p.id = ua.professional_id;

-- name: GetAppointmentsByClientWithStatus :many
SELECT 
    a.*,
    c.id AS client_id_full,
    c.first_name AS client_first_name,
    c.last_name AS client_last_name,
    c.phone_number AS client_phone_number,
    c.chat_id AS client_chat_id,
    p.id AS professional_id_full,
    p.username AS professional_username,
    p.first_name AS professional_first_name,
    p.last_name AS professional_last_name,
    p.phone_number AS professional_phone_number,
    p.chat_id AS professional_chat_id
FROM appointments a
LEFT JOIN clients c ON c.id = a.client_id
LEFT JOIN professionals p ON p.id = a.professional_id
WHERE a.client_id = $1
  AND a.status = $2
  AND a.start_time > NOW()
  AND a.type = 'appointment'
ORDER BY a.start_time DESC;

-- name: CancelAppointmentByClientWithDetails :one
WITH updated_appointment AS (
    UPDATE appointments
    SET 
        status = 'cancelled',
        cancellation_reason = $3,
        cancelled_by_client_id = $2,
        updated_at = NOW()
    WHERE appointments.id = $1 
    AND client_id = $2
    AND status IN ('pending', 'confirmed')
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
    ua.cancellation_reason,
    ua.cancelled_by_professional_id,
    ua.cancelled_by_client_id,
    ua.created_at,
    ua.updated_at,
    c.id as client_id_full,
    c.first_name as client_first_name,
    c.last_name as client_last_name,
    c.phone_number as client_phone_number,
    c.chat_id as client_chat_id,
    p.id as professional_id_full,
    p.username as professional_username,
    p.first_name as professional_first_name,
    p.last_name as professional_last_name,
    p.phone_number as professional_phone_number,
    p.chat_id as professional_chat_id
FROM updated_appointment ua
LEFT JOIN clients c ON c.id = ua.client_id
LEFT JOIN professionals p ON p.id = ua.professional_id;

-- name: CreateUnavailableAppointment :one
INSERT INTO appointments (type, professional_id, start_time, end_time, status, description)
VALUES ('unavailable', $1, $2, $3, 'confirmed', $4)
RETURNING *;

-- name: GetAppointmentsByProfessionalAndDate :many
SELECT * FROM appointments
WHERE professional_id = $1
  AND DATE(start_time) = $2
  AND type = 'appointment' or type = 'unavailable'
  AND status not in ('cancelled', 'pending')
ORDER BY start_time ASC;

-- name: GetAppointmentsByProfessionalAndDateWithClient :many
SELECT 
    a.id,
    a.professional_id,
    a.client_id,
    a.start_time,
    a.end_time,
    a.description,
    a.type,
    a.status,
    a.created_at,
    a.updated_at,
    c.first_name as client_first_name,
    c.last_name as client_last_name
FROM appointments a
LEFT JOIN clients c ON a.client_id = c.id
WHERE a.professional_id = $1
  AND DATE(a.start_time) = $2
  AND (a.type = 'appointment' OR a.type = 'unavailable')
  AND a.status NOT IN ('cancelled', 'pending')
ORDER BY a.start_time ASC;
