-- name: CreateSubscription :exec
INSERT INTO subscriptions (client_id, professional_id)
VALUES ($1, $2);

-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE client_id = $1 AND professional_id = $2;

-- name: GetSubscriptionsByClientID :many
SELECT p.id, p.first_name, p.last_name, p.chat_id FROM subscriptions s
JOIN professionals p ON s.professional_id = p.id
WHERE s.client_id = $1
ORDER BY p.first_name, p.last_name
LIMIT $2 OFFSET $3;

-- name: CountSubscriptionsByClientID :one
SELECT COUNT(*) FROM subscriptions
WHERE client_id = $1;

