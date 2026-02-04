-- name: CreateProfessional :one
INSERT INTO professionals (username, first_name, last_name, password_hash, chat_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: GetProfessionalInfoForNotification :one
SELECT chat_id, locale
FROM professionals
WHERE id = $1;

-- name: GetProfessionalInfoForNotificationByAppointmentID :one
SELECT chat_id, locale
FROM professionals
WHERE id = (
    SELECT professional_id FROM appointments a WHERE a.id = $1
);


-- name: GetProfessionalByUsername :one
SELECT * FROM professionals
WHERE username = $1;

-- name: GetProfessionals :many
SELECT id, first_name, last_name, chat_id, locale
FROM professionals p
WHERE p.chat_id is not null
AND p.id NOT IN (
    SELECT DISTINCT(professional_id)
    FROM subscriptions
    WHERE client_id = $1
)
ORDER BY p.first_name, p.last_name
LIMIT $2 OFFSET $3;

-- name: CountProfessionals :one
SELECT COUNT(*) FROM professionals p
WHERE p.chat_id is not null
AND p.id NOT IN (
    SELECT DISTINCT(professional_id)
    FROM subscriptions
    WHERE client_id = $1
)
;

-- name: UpdateProfessionalChatID :one
UPDATE professionals
SET chat_id = $2, locale = $3
WHERE id = $1
RETURNING *;

-- name: GetAppointmentsByProfessionalWithStatusAndDate :many
SELECT 
    a.id,
    a.start_time,
    a.end_time,
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
  AND ($2::appointment_status IS NULL OR a.status = $2)
  AND a.start_time > NOW()
  AND a.type != 'unavailable'
GROUP BY a.id, a.start_time, a.end_time
ORDER BY a.start_time ASC
LIMIT $3 OFFSET $4;

-- name: CountProfessionalAppointmentsWithStatusAndDate :one
SELECT COUNT(*)
FROM appointments a
WHERE a.professional_id = $1
    AND ($2 = '' OR a.status = $2::appointment_status)
    AND ($3 = '' OR DATE(a.start_time) = $3::date)
    AND a.start_time > NOW()
    AND a.type != 'unavailable';

-- name: GetProfessionalClients :many
SELECT id, first_name, last_name
FROM clients
WHERE id IN (
    SELECT DISTINCT(client_id)
    FROM appointments
    WHERE professional_id = $1
        AND client_id IS NOT NULL
        AND status = 'confirmed'
)
ORDER BY first_name, last_name;


-- name: UpdateProfessionalLocale :exec
UPDATE professionals
SET locale = $2
WHERE id = $1;

-- name: CreateGroupVisitAppointment :exec
INSERT INTO appointments (type, professional_id, start_time, end_time, status, description)
VALUES ('group', $1, $2, $3, 'confirmed', $4);
