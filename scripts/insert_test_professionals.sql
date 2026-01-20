-- Insert 15 test professionals (coaches) with different names and data
-- All passwords are hashed with bcrypt for "password123" (you should change these in production)

INSERT INTO professionals (username, first_name, last_name, phone_number, password_hash, chat_id) VALUES
('coach_alex', 'Alex', 'Johnson', '+1234567890', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000001),
('coach_sarah', 'Sarah', 'Williams', '+1234567891', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000002),
('coach_michael', 'Michael', 'Brown', '+1234567892', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000003),
('coach_emily', 'Emily', 'Davis', '+1234567893', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000004),
('coach_david', 'David', 'Miller', '+1234567894', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000005),
('coach_jessica', 'Jessica', 'Wilson', '+1234567895', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000006),
('coach_james', 'James', 'Moore', '+1234567896', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000007),
('coach_olivia', 'Olivia', 'Taylor', '+1234567897', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000008),
('coach_robert', 'Robert', 'Anderson', '+1234567898', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000009),
('coach_sophia', 'Sophia', 'Thomas', '+1234567899', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000010),
('coach_william', 'William', 'Jackson', '+1234567900', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000011),
('coach_isabella', 'Isabella', 'White', '+1234567901', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000012),
('coach_christopher', 'Christopher', 'Harris', '+1234567902', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000013),
('coach_ava', 'Ava', 'Martin', '+1234567903', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000014),
('coach_daniel', 'Daniel', 'Thompson', '+1234567904', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 1000015)
ON CONFLICT (username) DO NOTHING;

-- Verify the insert
SELECT id, username, first_name, last_name, phone_number, chat_id, created_at 
FROM professionals 
ORDER BY created_at DESC 
LIMIT 15;
