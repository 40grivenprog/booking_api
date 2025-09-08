-- name: GetUserByChatID :one
SELECT 
    id,
    chat_id,
    first_name,
    last_name,
    phone_number,
    created_at,
    updated_at,
    'client' as role,
    NULL as username
FROM clients 
WHERE clients.chat_id = $1

UNION ALL

SELECT 
    id,
    chat_id,
    first_name,
    last_name,
    phone_number,
    created_at,
    updated_at,
    'professional' as role,
    username
FROM professionals 
WHERE professionals.chat_id = $1;
