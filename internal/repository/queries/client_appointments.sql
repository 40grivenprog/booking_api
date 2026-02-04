-- name: CreateClientAppointment :exec
INSERT INTO client_appointments (client_id, appointment_id)
VALUES ($1, $2);
