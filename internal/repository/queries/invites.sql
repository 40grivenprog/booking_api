-- name: CreateInvites :exec
WITH a AS (
  SELECT x AS appointment_id, ord
  FROM unnest($1::uuid[]) WITH ORDINALITY AS t(x, ord)
),
b AS (
  SELECT x AS start_time, ord
  FROM unnest($2::timestamptz[]) WITH ORDINALITY AS t(x, ord)
),
c AS (
  SELECT x AS end_time, ord
  FROM unnest($3::timestamptz[]) WITH ORDINALITY AS t(x, ord)
),
d AS (
  SELECT x AS client_id, ord
  FROM unnest($4::uuid[]) WITH ORDINALITY AS t(x, ord)
),
e AS (
  SELECT x AS description, ord
  FROM unnest($5::text[]) WITH ORDINALITY AS t(x, ord)
),
f AS (
  SELECT x AS type, ord
  FROM unnest($6::appointment_type[]) WITH ORDINALITY AS t(x, ord)
),
g AS (
  SELECT x AS professional_name, ord
  FROM unnest($7::text[]) WITH ORDINALITY AS t(x, ord)
)
INSERT INTO invites (appointment_id, start_time, end_time, client_id, description, type, professional_name)
SELECT a.appointment_id, b.start_time, c.end_time, d.client_id, e.description, f.type, g.professional_name
FROM a
JOIN b USING (ord)
JOIN c USING (ord)
JOIN d USING (ord)
JOIN e USING (ord)
JOIN f USING (ord)
JOIN g USING (ord);


-- name: GetInvitesByAppointmentIDAndClientIs :many
SELECT id, client_id FROM invites
WHERE appointment_id = $1
AND client_id = ANY($2::uuid[]);

-- name: GetClientInvites :many
SELECT id, appointment_id, start_time, end_time, client_id, description, type, professional_name FROM invites
WHERE client_id = $1
AND start_time > NOW()
ORDER BY start_time DESC;

-- name: GetInviteByID :one
SELECT id, appointment_id, start_time, end_time, client_id, description, type, professional_name FROM invites
WHERE id = $1
AND client_id = $2;

-- name: DeleteInviteByID :exec
DELETE FROM invites
WHERE id = $1
AND client_id = $2;


-- name: GetInfoForAcceptInviteNotification :one
select a.id, a.description, a.type, a.start_time, a.end_time, p.chat_id, p.locale from appointments a
left join professionals p on p.id = a.professional_id
where a.id = $1;