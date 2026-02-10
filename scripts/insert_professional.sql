-- Insert professional user
-- Password: password1
-- To generate a new hash, run: go run scripts/hash_password.go <password>

INSERT INTO professionals (username, first_name, last_name, password_hash, chat_id)
VALUES (
    'nefedman',
    'Viktor',
    'Nefedov',
    '$2a$10$YOUR_HASHED_PASSWORD_HERE',  -- Replace with actual bcrypt hash
    NULL  -- Set chat_id if you have a Telegram chat ID, otherwise NULL
)
ON CONFLICT (username) DO NOTHING;

-- To generate the password hash, run this command in the booking_api directory:
-- go run scripts/hash_password.go password1
--
-- Then replace $2a$10$YOUR_HASHED_PASSWORD_HERE with the output
