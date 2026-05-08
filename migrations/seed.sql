-- ==========================================
-- APARTMENT MANAGEMENT SYSTEM - SEED DATA
-- ==========================================

-- Insert Users (3 users: 2 staff, 1 tenant)
INSERT INTO users (id, name, email, phone, password, role, created_at, updated_at) VALUES
(
  '550e8400-e29b-41d4-a716-446655440001',
  'Admin Staff',
  'admin@example.com',
  '0811111111',
  'mysecretpassword', -- Plaintext password for reference (should be hashed in production)
  'STAFF',
  NOW(),
  NOW()
),
(
  '550e8400-e29b-41d4-a716-446655440002',
  'Manager Staff',
  'manager@apartment.local',
  '0812222222',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36XQuvG',
  'STAFF',
  NOW(),
  NOW()
),
(
  '550e8400-e29b-41d4-a716-446655440003',
  'John Tenant',
  'john@tenant.local',
  '0813333333',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36XQuvG',
  'TENANT',
  NOW(),
  NOW()
);

-- Insert Rooms (5 rooms across different floors)
INSERT INTO rooms (id, room_number, level, status, created_at, updated_at) VALUES
(
  '550e8400-e29b-41d4-a716-446655440010',
  '101',
  1,
  'Occupied',
  NOW(),
  NOW()
),
(
  '550e8400-e29b-41d4-a716-446655440011',
  '102',
  1,
  'Occupied',
  NOW(),
  NOW()
),
(
  '550e8400-e29b-41d4-a716-446655440012',
  '201',
  2,
  'Available',
  NOW(),
  NOW()
),
(
  '550e8400-e29b-41d4-a716-446655440013',
  '202',
  2,
  'Occupied',
  NOW(),
  NOW()
),
(
  '550e8400-e29b-41d4-a716-446655440014',
  '301',
  3,
  'Available',
  NOW(),
  NOW()
);

-- Insert Utility Rates (1 rate record)
INSERT INTO utility_rates (id, water_rate, electric_rate, common_fee, period, created_at, updated_at) VALUES
(
  '550e8400-e29b-41d4-a716-446655440020',
  5.50,
  3.25,
  150.00,
  '2024-01',
  NOW(),
  NOW()
);

-- Insert Contracts (2 contracts linked to tenants and rooms)
INSERT INTO contracts (id, user_id, room_id, start_date, end_date, status, created_at, updated_at) VALUES
(
  '550e8400-e29b-41d4-a716-446655440030',
  '550e8400-e29b-41d4-a716-446655440003',
  '550e8400-e29b-41d4-a716-446655440010',
  '2024-01-01'::date,
  '2025-12-31'::date,
  'Active',
  NOW(),
  NOW()
),
(
  '550e8400-e29b-41d4-a716-446655440031',
  '550e8400-e29b-41d4-a716-446655440003',
  '550e8400-e29b-41d4-a716-446655440011',
  '2024-02-01'::date,
  '2025-12-31'::date,
  'Active',
  NOW(),
  NOW()
);

-- Insert Utility Usage Records (2 records linked to contracts)
-- Using meter readings (old_unit to new_unit) instead of consumption values
INSERT INTO utility_usages (id, contract_id, old_water_unit, new_water_unit, old_electric_unit, new_electric_unit, record_date, created_at, updated_at) VALUES
(
  '550e8400-e29b-41d4-a716-446655440040',
  '550e8400-e29b-41d4-a716-446655440030',
  100,
  125,
  500,
  650,
  '2024-01-31'::date,
  NOW(),
  NOW()
),
(
  '550e8400-e29b-41d4-a716-446655440041',
  '550e8400-e29b-41d4-a716-446655440031',
  200,
  230,
  600,
  780,
  '2024-01-31'::date,
  NOW(),
  NOW()
);

-- Insert Bills (Optional - for reference)
INSERT INTO bills (id, contract_id, rate_id, record_date, rent_fee, water_fee, electricity_fee, common_fee, total_amount, status, due_date, created_date, created_at, updated_at) VALUES
(
  '550e8400-e29b-41d4-a716-446655440050',
  '550e8400-e29b-41d4-a716-446655440030',
  '550e8400-e29b-41d4-a716-446655440020',
  '2024-01-31'::date,
  2500000.00,
  137.50,
  487.50,
  150.00,
  3275.00,
  'Unpaid',
  '2024-02-14'::date,
  '2024-02-01'::date,
  NOW(),
  NOW()
),
(
  '550e8400-e29b-41d4-a716-446655440051',
  '550e8400-e29b-41d4-a716-446655440031',
  '550e8400-e29b-41d4-a716-446655440020',
  '2024-01-31'::date,
  3000000.00,
  165.00,
  586.63,
  150.00,
  3901.63,
  'Unpaid',
  '2024-02-14'::date,
  '2024-02-01'::date,
  NOW(),
  NOW()
);
