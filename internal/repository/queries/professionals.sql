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

-- name: CreateGroupVisitAppointment :one
INSERT INTO appointments (type, professional_id, start_time, end_time, status, description)
VALUES ($1, $2, $3, $4, 'confirmed', $5)
RETURNING id;


-- name: GetPreviousAppointmentsByProfessionalID :many
SELECT 
    a.id, 
    a.start_time, 
    a.end_time, 
    a.description, 
    a.type, 
    COUNT(DISTINCT ca.client_id) AS users_count
FROM appointments a
LEFT JOIN client_appointments ca ON ca.appointment_id = a.id
WHERE a.professional_id = $1
  AND a.status = 'confirmed'
  AND a.type != 'unavailable'
  AND a.end_time < NOW()
  AND ($2::timestamptz IS NULL OR $2 < '1900-01-01'::timestamptz OR a.start_time >= $2)
  AND ($3::timestamptz IS NULL OR $3 < '1900-01-01'::timestamptz OR a.start_time < $3)
GROUP BY a.id, a.start_time, a.end_time, a.description, a.type
ORDER BY a.start_time DESC
LIMIT $4 OFFSET $5;

-- name: CountPreviousAppointmentsByProfessionalID :one
SELECT COUNT(DISTINCT a.id)
FROM appointments a
WHERE a.professional_id = $1
  AND a.status = 'confirmed'
  AND a.type != 'unavailable'
  AND a.end_time < NOW()
  AND ($2::timestamptz IS NULL OR $2 < '1900-01-01'::timestamptz OR a.start_time >= $2)
  AND ($3::timestamptz IS NULL OR $3 < '1900-01-01'::timestamptz OR a.start_time < $3);

-- name: GetPreviousAppointmentsByProfessionalIDAndClientID :many
SELECT a.id, a.start_time, a.end_time, a.description, a.type, COUNT(DISTINCT ca.client_id) AS users_count
FROM appointments a
LEFT JOIN client_appointments ca ON ca.appointment_id = a.id
WHERE a.professional_id = $1
  AND a.status = 'confirmed'
  AND a.type != 'unavailable'
  AND a.end_time < NOW()
  AND a.id IN (
    SELECT DISTINCT(appointment_id)
    FROM client_appointments ca
    WHERE ca.client_id = $2
  )
  AND ($3::timestamptz IS NULL OR $3 < '1900-01-01'::timestamptz OR a.start_time >= $3)
  AND ($4::timestamptz IS NULL OR $4 < '1900-01-01'::timestamptz OR a.start_time < $4)
GROUP BY a.id, a.start_time, a.end_time, a.description, a.type
ORDER BY a.start_time DESC
LIMIT $5 OFFSET $6;

-- name: CountPreviousAppointmentsByProfessionalIDAndClientID :one
SELECT COUNT(DISTINCT a.id)
FROM appointments a
WHERE a.professional_id = $1
  AND a.status = 'confirmed'
  AND a.type != 'unavailable'
  AND a.end_time < NOW()
  AND a.id IN (
    SELECT DISTINCT(appointment_id)
    FROM client_appointments ca
    WHERE ca.client_id = $2
  )
  AND ($3::timestamptz IS NULL OR $3 < '1900-01-01'::timestamptz OR a.start_time >= $3)
  AND ($4::timestamptz IS NULL OR $4 < '1900-01-01'::timestamptz OR a.start_time < $4);


-- name: GetAppointmentDetailsByProfessionalIDAndAppointmentID :one
SELECT a.id, a.start_time, a.end_time, a.description, a.type
FROM appointments a
WHERE a.professional_id = $1
AND a.id = $2;


-- name: GetAppintmentClientsCountByAppointmentID :one
SELECT COUNT(*)
FROM client_appointments ca
WHERE ca.appointment_id = $1;


-- name: GetMissingInviteUsersForAppointment :many
SELECT c.id, c.first_name, c.last_name, c.chat_id, c.locale
FROM clients c
WHERE c.id NOT IN (
    SELECT DISTINCT(client_id)
    FROM invites i
    WHERE i.appointment_id = $1
)
AND c.id NOT IN(
  SELECT DISTINCT(client_id)
  FROM client_appointments ca
  WHERE ca.appointment_id = $1
);

-- name: GetPendingInviteUsersForAppointment :many
SELECT c.id, c.first_name, c.last_name
FROM clients c
WHERE c.id IN (
    SELECT DISTINCT(client_id)
    FROM invites i
    WHERE i.appointment_id = $1
);

-- name: GetMissingClientsForPreviousAppointment :many
SELECT c.id, c.first_name, c.last_name, c.chat_id, c.locale
FROM clients c
WHERE c.id NOT IN (
    SELECT DISTINCT(client_id)
    FROM client_appointments ca
    WHERE ca.appointment_id = $1
)
AND c.id IN (
  SELECT s.client_id
  FROM subscriptions s
  WHERE s.professional_id = $2
);

-- name: GetClientInfoByID :one
SELECT id, first_name, last_name, chat_id, locale FROM clients WHERE id = $1;


-- name: CheckProfessionalAppointmentConflict :one
SELECT EXISTS(
    SELECT 1 FROM appointments a
    WHERE a.professional_id = $1
    AND a.start_time = $2
    AND a.end_time = $3
    AND a.status = 'confirmed'
);