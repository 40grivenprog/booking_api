-- name: CreatePersonalAppointment :one
INSERT INTO appointments (type, professional_id, start_time, end_time, status)
VALUES ('personal', $1, $2, $3, 'pending')
RETURNING id, type, professional_id, start_time, end_time, status, created_at, updated_at;

-- name: GetAppointmentWithDetails :one
SELECT 
    a.id,
    a.type,
    a.professional_id,
    a.start_time,
    a.end_time,
    a.status,
    a.created_at,
    a.updated_at,
    c.id as client_id,
    c.first_name as client_first_name,
    c.last_name as client_last_name,
    c.chat_id as client_chat_id,
    c.locale as client_locale,
    p.chat_id as professional_chat_id,
    p.locale as professional_locale
FROM appointments a
LEFT JOIN client_appointments ca ON ca.appointment_id = a.id
LEFT JOIN clients c ON c.id = ca.client_id
LEFT JOIN professionals p ON p.id = a.professional_id
WHERE a.id = $1;

-- name: GetAppointmentByID :one
SELECT * FROM appointments
WHERE id = $1;

-- name: DeleteAppointmentById :exec
DELETE FROM appointments
WHERE id = $1;

-- name: GetAppointmentInfoByAppointmentID :one
SELECT start_time, end_time FROM appointments a
WHERE a.id = $1;

-- -- name: CreateAppointmentWithDetails :one
-- WITH new_appointment AS (
--     INSERT INTO appointments (type, professional_id, start_time, end_time, status, description)
--     VALUES ('appointment', $1, $2, $3, $4, 'pending', $5)
--     RETURNING *
-- )
-- SELECT 
--     na.start_time,
--     na.end_time,
--     na.description,
--     c.first_name as client_first_name,
--     c.last_name as client_last_name,
--     p.chat_id as professional_chat_id,
--     p.locale as professional_locale
-- FROM new_appointment na
-- LEFT JOIN clients c ON c.id = na.client_id
-- LEFT JOIN professionals p ON p.id = na.professional_id;

-- name: ConfirmAppointmentById :exec
UPDATE appointments
SET status = 'confirmed', updated_at = NOW()
WHERE id = $1;

-- name: GetAppointmentClientsByAppointmentID :many
SELECT c.chat_id as client_chat_id, c.locale as client_locale
FROM client_appointments ca
JOIN clients c ON c.id = ca.client_id
WHERE ca.appointment_id = $1;


-- -- name: GetAppointmentsByProfessionalWithStatus :many
-- SELECT 
--     a.*,
--     c.id AS client_id,
--     c.first_name AS client_first_name,
--     c.last_name AS client_last_name,
--     c.chat_id AS client_chat_id
-- FROM appointments a
-- LEFT JOIN clients c ON c.id = a.client_id
-- WHERE a.professional_id = $1
--   AND a.status = $2
--   AND a.start_time > NOW()
--   AND a.type = 'appointment'
-- ORDER BY a.start_time DESC;


-- -- name: CancelAppointmentByProfessionalWithDetails :one
-- WITH updated_appointment AS (
--     UPDATE appointments
--     SET 
--         status = 'cancelled',
--         cancellation_reason = $3,
--         cancelled_by_professional_id = $2,
--         updated_at = NOW()
--     WHERE appointments.id = $1 
--     AND professional_id = $2
--     AND status IN ('pending', 'confirmed')
--     RETURNING *
-- )
-- SELECT 
--     ua.id,
--     ua.type,
--     ua.professional_id,
--     ua.start_time,
--     ua.end_time,
--     ua.status,
--     ua.cancellation_reason,
--     ua.cancelled_by_professional_id,
--     ua.cancelled_by_client_id,
--     ua.created_at,
--     ua.updated_at,
--     c.id as client_id_full,
--     c.first_name as client_first_name,
--     c.last_name as client_last_name,
--     c.phone_number as client_phone_number,
--     c.chat_id as client_chat_id,
--     c.locale as client_locale,
--     p.id as professional_id_full,
--     p.username as professional_username,
--     p.first_name as professional_first_name,
--     p.last_name as professional_last_name,
--     p.phone_number as professional_phone_number,
--     p.chat_id as professional_chat_id
-- FROM updated_appointment ua
-- LEFT JOIN clients c ON c.id = ua.client_id
-- LEFT JOIN professionals p ON p.id = ua.professional_id;

-- name: GetAppointmentsByClientWithStatus :many
SELECT 
    a.id,
    a.start_time,
    a.end_time,
    a.type,
    a.status,
    p.first_name AS professional_first_name,
    p.last_name AS professional_last_name,
    p.chat_id AS professional_chat_id
FROM appointments a
LEFT JOIN professionals p ON p.id = a.professional_id
WHERE
  a.id in (
    select appointment_id from client_appointments where client_id = $1
  )
  AND a.status = $2
  AND a.start_time > NOW()
  AND a.type != 'unavailable'
ORDER BY a.start_time ASC
LIMIT $3 OFFSET $4;

-- name: CountClientAppointmentsWithStatus :one
SELECT COUNT(*)
FROM appointments a
WHERE
    a.id in (
        select appointment_id from client_appointments where client_id = $1
    )
  AND a.status = $2
  AND a.start_time > NOW()
  AND a.type != 'unavailable';

-- -- name: CancelAppointmentByClientWithDetails :one
-- WITH updated_appointment AS (
--     UPDATE appointments
--     SET 
--         status = 'cancelled',
--         cancellation_reason = $3,
--         cancelled_by_client_id = $2,
--         updated_at = NOW()
--     WHERE appointments.id = $1 
--     AND client_id = $2
--     AND status IN ('pending', 'confirmed')
--     RETURNING *
-- )
-- SELECT 
--     ua.start_time,
--     ua.end_time,
--     ua.cancellation_reason,
--     ua.description,
--     c.first_name as client_first_name,
--     c.last_name as client_last_name,
--     p.chat_id as professional_chat_id,
--     p.locale as professional_locale
-- FROM updated_appointment ua
-- LEFT JOIN clients c ON c.id = ua.client_id
-- LEFT JOIN professionals p ON p.id = ua.professional_id;

-- name: CreateUnavailableAppointment :exec
INSERT INTO appointments (type, professional_id, start_time, end_time, status, description)
VALUES ('unavailable', $1, $2, $3, 'confirmed', $4);

-- -- name: GetAppointmentsByProfessionalAndDate :many
-- SELECT * FROM appointments
-- WHERE professional_id = $1
--   AND DATE(start_time) = $2
--   AND type = 'appointment' or type = 'unavailable'
--   AND status not in ('cancelled', 'pending')
-- ORDER BY start_time ASC;

-- name: CheckClientAppointmentConflict :one
SELECT EXISTS(
    SELECT 1 FROM appointments a
    INNER JOIN client_appointments ca ON ca.appointment_id = a.id
    WHERE ca.client_id = $1 
    AND a.professional_id = $2 
    AND a.start_time = $3 
    AND a.status IN ('pending', 'confirmed')
) as has_conflict;

-- name: GetAppointmentsByProfessionalByDate :many
SELECT 
    a.start_time,
    a.end_time
FROM appointments a
WHERE a.professional_id = $1
  AND DATE(a.start_time) = $2
  AND a.status NOT IN ('pending')
ORDER BY a.start_time ASC;

-- -- name: GetProfessionalAppointmentDates :many
-- SELECT DISTINCT DATE(start_time) AS appointment_date
-- FROM appointments
-- WHERE professional_id = $1
--   AND type = 'appointment'
--   AND status = 'confirmed'
--   AND start_time >= $2
--   AND start_time < $3
-- ORDER BY appointment_date ASC;

-- name: GetProfessionalTimetable :many
SELECT 
    a.id,
    a.start_time,
    a.end_time,
    a.description,
    a.type,
    COALESCE(
        ARRAY_AGG(
            c.first_name || ' ' || c.last_name
            ORDER BY c.last_name
        ) FILTER (WHERE c.id IS NOT NULL),
        ARRAY[]::text[]
    )::text[] AS clients
FROM appointments a
LEFT JOIN client_appointments ca ON ca.appointment_id = a.id
LEFT JOIN clients c ON c.id = ca.client_id
WHERE a.professional_id = $1
  AND a.status = 'confirmed'
  AND DATE(a.start_time) = $2
  AND a.end_time > NOW()
GROUP BY a.id, a.start_time, a.end_time, a.description
ORDER BY a.start_time ASC;

-- -- name: GetPreviousProfessionalAppointmentsByClient :many
-- SELECT id, start_time, end_time, description
-- FROM appointments
-- WHERE client_id = $1
--   AND professional_id = $2
--   AND status = 'confirmed'
--   AND end_time < NOW()
-- ORDER BY start_time DESC;

-- -- name: GetPreviousAppointmentsByClientForMonth :many
-- SELECT id, start_time, end_time, description
-- FROM appointments
-- WHERE client_id = sqlc.arg(client_id)
--   AND professional_id = sqlc.arg(professional_id)
--   AND status = 'confirmed'
--   AND end_time < NOW()
--   AND DATE_TRUNC('month', start_time) = DATE_TRUNC('month', sqlc.arg(month_date)::date)
-- ORDER BY start_time DESC;
