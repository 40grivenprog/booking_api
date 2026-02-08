-- Create test appointments in the past
-- Professional ID: 2811eb7e-b1a8-485d-80ad-135e94d4bb57
-- Client IDs: 13e4a18d-4374-4e25-9818-52e102b8bb16, 6d6832bb-ebe7-4ef6-9c11-0f431fd9bbd4

-- Step 1: Create appointments
-- Dates: 2026-02-01, 2026-02-02, 2026-02-03, 2026-02-04, 2026-02-05
-- Times: 10:00-11:00 and 12:00-13:00

INSERT INTO appointments (id, type, professional_id, start_time, end_time, status, description) VALUES
-- February 1st
('a1b2c3d4-e5f6-4789-a012-345678901234', 'personal', '2811eb7e-b1a8-485d-80ad-135e94d4bb57', '2026-02-01 10:00:00+00:00', '2026-02-01 11:00:00+00:00', 'confirmed', 'Test appointment 1'),
('a1b2c3d4-e5f6-4789-a012-345678901235', 'personal', '2811eb7e-b1a8-485d-80ad-135e94d4bb57', '2026-02-01 12:00:00+00:00', '2026-02-01 13:00:00+00:00', 'confirmed', 'Test appointment 2'),

-- February 2nd
('a1b2c3d4-e5f6-4789-a012-345678901236', 'personal', '2811eb7e-b1a8-485d-80ad-135e94d4bb57', '2026-02-02 10:00:00+00:00', '2026-02-02 11:00:00+00:00', 'confirmed', 'Test appointment 3'),
('a1b2c3d4-e5f6-4789-a012-345678901237', 'personal', '2811eb7e-b1a8-485d-80ad-135e94d4bb57', '2026-02-02 12:00:00+00:00', '2026-02-02 13:00:00+00:00', 'confirmed', 'Test appointment 4'),

-- February 3rd
('a1b2c3d4-e5f6-4789-a012-345678901238', 'personal', '2811eb7e-b1a8-485d-80ad-135e94d4bb57', '2026-02-03 10:00:00+00:00', '2026-02-03 11:00:00+00:00', 'confirmed', 'Test appointment 5'),
('a1b2c3d4-e5f6-4789-a012-345678901239', 'personal', '2811eb7e-b1a8-485d-80ad-135e94d4bb57', '2026-02-03 12:00:00+00:00', '2026-02-03 13:00:00+00:00', 'confirmed', 'Test appointment 6'),

-- February 4th
('a1b2c3d4-e5f6-4789-a012-345678901240', 'personal', '2811eb7e-b1a8-485d-80ad-135e94d4bb57', '2026-02-04 10:00:00+00:00', '2026-02-04 11:00:00+00:00', 'confirmed', 'Test appointment 7'),
('a1b2c3d4-e5f6-4789-a012-345678901241', 'personal', '2811eb7e-b1a8-485d-80ad-135e94d4bb57', '2026-02-04 12:00:00+00:00', '2026-02-04 13:00:00+00:00', 'confirmed', 'Test appointment 8'),

-- February 5th
('a1b2c3d4-e5f6-4789-a012-345678901242', 'personal', '2811eb7e-b1a8-485d-80ad-135e94d4bb57', '2026-02-05 10:00:00+00:00', '2026-02-05 11:00:00+00:00', 'confirmed', 'Test appointment 9'),
('a1b2c3d4-e5f6-4789-a012-345678901243', 'personal', '2811eb7e-b1a8-485d-80ad-135e94d4bb57', '2026-02-05 12:00:00+00:00', '2026-02-05 13:00:00+00:00', 'confirmed', 'Test appointment 10');

-- Step 2: Create client_appointments for each appointment
-- Each appointment will have both clients linked

INSERT INTO client_appointments (client_id, appointment_id) VALUES
-- February 1st appointments
('13e4a18d-4374-4e25-9818-52e102b8bb16', 'a1b2c3d4-e5f6-4789-a012-345678901234'),
('6d6832bb-ebe7-4ef6-9c11-0f431fd9bbd4', 'a1b2c3d4-e5f6-4789-a012-345678901234'),
('13e4a18d-4374-4e25-9818-52e102b8bb16', 'a1b2c3d4-e5f6-4789-a012-345678901235'),
('6d6832bb-ebe7-4ef6-9c11-0f431fd9bbd4', 'a1b2c3d4-e5f6-4789-a012-345678901235'),

-- February 2nd appointments
('13e4a18d-4374-4e25-9818-52e102b8bb16', 'a1b2c3d4-e5f6-4789-a012-345678901236'),
('6d6832bb-ebe7-4ef6-9c11-0f431fd9bbd4', 'a1b2c3d4-e5f6-4789-a012-345678901236'),
('13e4a18d-4374-4e25-9818-52e102b8bb16', 'a1b2c3d4-e5f6-4789-a012-345678901237'),
('6d6832bb-ebe7-4ef6-9c11-0f431fd9bbd4', 'a1b2c3d4-e5f6-4789-a012-345678901237'),

-- February 3rd appointments
('13e4a18d-4374-4e25-9818-52e102b8bb16', 'a1b2c3d4-e5f6-4789-a012-345678901238'),
('6d6832bb-ebe7-4ef6-9c11-0f431fd9bbd4', 'a1b2c3d4-e5f6-4789-a012-345678901238'),
('13e4a18d-4374-4e25-9818-52e102b8bb16', 'a1b2c3d4-e5f6-4789-a012-345678901239'),
('6d6832bb-ebe7-4ef6-9c11-0f431fd9bbd4', 'a1b2c3d4-e5f6-4789-a012-345678901239'),

-- February 4th appointments
('13e4a18d-4374-4e25-9818-52e102b8bb16', 'a1b2c3d4-e5f6-4789-a012-345678901240'),
('6d6832bb-ebe7-4ef6-9c11-0f431fd9bbd4', 'a1b2c3d4-e5f6-4789-a012-345678901240'),
('13e4a18d-4374-4e25-9818-52e102b8bb16', 'a1b2c3d4-e5f6-4789-a012-345678901241'),
('6d6832bb-ebe7-4ef6-9c11-0f431fd9bbd4', 'a1b2c3d4-e5f6-4789-a012-345678901241'),

-- February 5th appointments
('13e4a18d-4374-4e25-9818-52e102b8bb16', 'a1b2c3d4-e5f6-4789-a012-345678901242'),
('6d6832bb-ebe7-4ef6-9c11-0f431fd9bbd4', 'a1b2c3d4-e5f6-4789-a012-345678901242'),
('13e4a18d-4374-4e25-9818-52e102b8bb16', 'a1b2c3d4-e5f6-4789-a012-345678901243'),
('6d6832bb-ebe7-4ef6-9c11-0f431fd9bbd4', 'a1b2c3d4-e5f6-4789-a012-345678901243');
