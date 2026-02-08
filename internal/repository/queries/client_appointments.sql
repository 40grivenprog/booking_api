-- name: CreateClientAppointment :exec
INSERT INTO client_appointments (client_id, appointment_id)
VALUES ($1, $2);

-- name: DeleteClientAppointmentsByAppointmentIDAndClientIDs :exec
DELETE FROM client_appointments
WHERE appointment_id = $1 AND client_id = ANY($2::uuid[]);
