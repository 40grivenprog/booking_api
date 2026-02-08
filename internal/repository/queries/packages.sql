-- name: CreatePackage :exec
INSERT INTO packages (client_id, professional_id, issued_at, expires_at, apppointments_number)
VALUES ($1, $2, $3, $4, $5);

-- name: DeactivatePackage :exec
UPDATE packages SET deactivated_at = NOW() WHERE id = $1;

-- name: GetPackageByClientIDAndProfessionalID :one
SELECT * FROM packages WHERE client_id = $1 AND professional_id = $2;

-- name: GetAppintmentNumberByClientIDAndProfessionalID :one
SELECT COUNT(*) FROM appointments
WHERE professional_id = $1
AND id IN (SELECT appointment_id FROM client_appointments WHERE client_id = $2)
AND start_time > $3;

-- name: GetListOfPackagesByProfessionalID :many
SELECT p.id, c.first_name, c.last_name, p.issued_at, p.expires_at, p.apppointments_number FROM packages p
LEFT JOIN clients c ON c.id = p.client_id
WHERE professional_id = $1
and deactivated_at is null;

-- name: GetPackageById :one
SELECT id, client_id, professional_id, issued_at, expires_at, apppointments_number, deactivated_at, created_at, updated_at FROM packages WHERE id = $1;

-- name: GetAppointmentsForThePackage :many
SELECT id, start_time, end_time, type FROM appointments
WHERE id IN (SELECT appointment_id FROM client_appointments WHERE client_id = $1)
and professional_id = $2
and start_time > $3
and end_time < NOW();

-- name: GetPackagesByClientId :many
SELECT p.id, p.issued_at, p.expires_at, p.apppointments_number, p.deactivated_at, pr.first_name, pr.last_name
FROM packages p
LEFT JOIN professionals pr
ON pr.id = p.professional_id
WHERE client_id = $1;