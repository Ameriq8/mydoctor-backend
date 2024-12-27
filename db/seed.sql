-- ======================================
-- Seed Data for Cities (18 Iraqi Cities)
-- ======================================
INSERT INTO cities (id, name, population, image_url, timezone)
VALUES
(1, 'Baghdad', 8000000, 'www.facility.com', 'Asia/Baghdad'),
(2, 'Basra', 2715000, 'www.facility.com', 'Asia/Baghdad'),
(3, 'Mosul', 1800000, 'www.facility.com', 'Asia/Baghdad'),
(4, 'Erbil', 1400000, 'www.facility.com', 'Asia/Baghdad'),
(5, 'Sulaymaniyah', 879000, 'www.facility.com', 'Asia/Baghdad'),
(6, 'Karbala', 700000, 'www.facility.com', 'Asia/Baghdad'),
(7, 'Najaf', 900000, 'www.facility.com', 'Asia/Baghdad'),
(8, 'Kirkuk', 950000, 'www.facility.com', 'Asia/Baghdad'),
(9, 'Duhok', 900000, 'www.facility.com', 'Asia/Baghdad'),
(10, 'Samarra', 150000, 'www.facility.com', 'Asia/Baghdad'),
(11, 'Ramadi', 600000, 'www.facility.com', 'Asia/Baghdad'),
(12, 'Fallujah', 275000, 'www.facility.com', 'Asia/Baghdad'),
(13, 'Amara', 420000, 'www.facility.com', 'Asia/Baghdad'),
(14, 'Nasiriyah', 560000, 'www.facility.com', 'Asia/Baghdad'),
(15, 'Kut', 400000, 'www.facility.com', 'Asia/Baghdad'),
(16, 'Diwaniyah', 500000, 'www.facility.com', 'Asia/Baghdad'),
(17, 'Tikrit', 160000, 'www.facility.com', 'Asia/Baghdad'),
(18, 'Hillah', 530000, 'www.facility.com', 'Asia/Baghdad');

-- ======================================
-- Seed Data for Facility Categories (20 Rows)
-- ======================================
INSERT INTO facility_categories (id, name, description)
SELECT i, 'Category ' || i, 'Description for Category ' || i
FROM generate_series(1, 20) AS s(i);

-- ======================================
-- Seed Data for Facilities (10,000 Rows)
-- ======================================
INSERT INTO facilities (name, type, category_id, city_id, location, phone, email, website, bed_capacity, is_24_hours, has_emergency, description)
SELECT 
    'Facility ' || i, 
    (ARRAY['Public Hospital', 'Teaching Hospital', 'Private Hospital', 'Rehabilitation Center', 'Medical Complex', 'Clinic', 'Pharmacy', 'Laboratory', 'Imaging Center'])[ceil(random() * 9)]::facility_type,
    (SELECT id FROM facility_categories ORDER BY random() LIMIT 1), -- Random valid category_id
    (SELECT id FROM cities ORDER BY random() LIMIT 1), -- Ensure valid city_id
    'Location ' || i,
    '964-771-' || (1000000 + i)::TEXT,
    'facility' || i || '@example.com',
    'www.facility' || i || '.com',
    ceil(random() * 500), 
    random() > 0.5, 
    random() > 0.5, 
    'Description for Facility ' || i
FROM generate_series(1, 10000) AS s(i);

-- ======================================
-- Seed Data for Doctors (10,000 Rows)
-- ======================================
INSERT INTO doctors (name, specialty, primary_facility_id, contact_number, email)
SELECT 
    'Doctor ' || i, 
    (ARRAY['Cardiologist', 'Neurologist', 'Radiologist', 'General Practitioner', 'Pediatrician', 'Dentist', 'Dermatologist', 'Orthopedic', 'Surgeon', 'Oncologist'])[ceil(random() * 10)],
    (SELECT id FROM facilities ORDER BY random() LIMIT 1), -- Ensure valid primary_facility_id
    '964-771-' || (2000000 + i)::TEXT,
    'doctor' || i || '@example.com'
FROM generate_series(1, 10000) AS s(i);

-- ======================================
-- Seed Data for Facility Departments (10,000 Rows)
-- ======================================
INSERT INTO facility_departments (facility_id, name, description, floor_number)
SELECT 
    (SELECT id FROM facilities ORDER BY random() LIMIT 1), -- Ensure valid facility_id
    'Department ' || i, 
    'Description for Department ' || i, 
    (ARRAY['1', '2', '3', '4', 'Ground', 'Basement'])[ceil(random() * 6)]
FROM generate_series(1, 10000) AS s(i);

-- ======================================
-- Seed Data for Facility Equipment (10,000 Rows)
-- ======================================
INSERT INTO facility_equipment (facility_id, department_id, name, model, manufacturer, status)
SELECT 
    (SELECT id FROM facilities ORDER BY random() LIMIT 1), -- Ensure valid facility_id
    (SELECT id FROM facility_departments ORDER BY random() LIMIT 1), -- Ensure valid department_id
    'Equipment ' || i, 
    'Model ' || i, 
    'Manufacturer ' || ceil(random() * 100), 
    (ARRAY['Operational', 'Under Maintenance', 'Decommissioned'])[ceil(random() * 3)]
FROM generate_series(1, 10000) AS s(i);

-- ======================================
-- Seed Data for Facility Operating Hours
-- ======================================
INSERT INTO facility_operating_hours (facility_id, department_id, day_of_week, start_time, end_time, is_closed)
SELECT 
    facility_id, 
    department_id, 
    gs.day_of_week, 
    '08:00', 
    '16:00', 
    random() > 0.8
FROM (
    SELECT 
        (SELECT id FROM facilities ORDER BY random() LIMIT 1) AS facility_id,
        (SELECT id FROM facility_departments ORDER BY random() LIMIT 1) AS department_id,
        generate_series(0, 6) AS day_of_week -- Unique days of the week
) gs
LIMIT 10000; -- Generate 10,000 rows

-- ======================================
-- Seed Data for Facility Plans (10,000 Rows)
-- ======================================
INSERT INTO facility_plans (facility_id, plan_id, start_date, end_date, is_active)
SELECT DISTINCT ON (facility_id, plan_id) -- Ensure unique combinations
    facility_id,
    plan_id,
    CURRENT_DATE - (random() * 365)::INT, 
    CURRENT_DATE + (random() * 365)::INT, 
    random() > 0.5
FROM (
    SELECT 
        (SELECT id FROM facilities ORDER BY random() LIMIT 1) AS facility_id,
        generate_series(1, 3) AS plan_id -- Assume 3 unique plans
) subquery
LIMIT 10000; -- Generate up to 10,000 rows

-- ======================================
-- Seed Data for Reviews (10,000 Rows)
-- ======================================
INSERT INTO reviews (entity_type, entity_id, user_id, rating, comment)
SELECT 
    (ARRAY['facility', 'doctor'])[ceil(random() * 2)], 
    (CASE WHEN random() > 0.5 THEN (SELECT id FROM facilities ORDER BY random() LIMIT 1) ELSE (SELECT id FROM doctors ORDER BY random() LIMIT 1) END),
    ceil(random() * 10000), 
    ROUND((random() * 5)::NUMERIC, 2), -- Corrected rounding for rating
    'Review comment ' || i
FROM generate_series(1, 10000) AS s(i);
