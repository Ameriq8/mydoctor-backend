-- ======================================
-- 1) Enable necessary extensions
-- ======================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS btree_gist;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- ======================================
-- 2) Create ENUM type for facility_type
-- ======================================
CREATE TYPE facility_type AS ENUM (
    'Public Hospital',
    'Teaching Hospital',
    'Private Hospital',
    'Rehabilitation Center',
    'Medical Complex',
    'Clinic',
    'Pharmacy',
    'Laboratory',
    'Imaging Center'
);

-- ======================================
-- 3) Create facility_categories
-- ======================================
CREATE TABLE facility_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id BIGINT REFERENCES facility_categories(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ======================================
-- 4) Create cities
-- ======================================
CREATE TABLE cities (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    population INTEGER,
    image_url VARCHAR(255),
    timezone VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ======================================
-- 5) Create facilities (referenced by doctors, etc.)
-- ======================================
CREATE TABLE facilities (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    type facility_type NOT NULL,
    category_id BIGINT REFERENCES facility_categories(id),
    city_id BIGINT REFERENCES cities(id),
    location VARCHAR(255) NOT NULL,
    coordinates POINT,
    phone VARCHAR(20),
    emergency_phone VARCHAR(20),
    email VARCHAR(255),
    website VARCHAR(255),
    rating DECIMAL(3,2),
    bed_capacity INTEGER,
    is_24_hours BOOLEAN DEFAULT false,
    has_emergency BOOLEAN DEFAULT false,
    has_parking BOOLEAN DEFAULT false,
    has_ambulance BOOLEAN DEFAULT false,
    accepts_insurance BOOLEAN DEFAULT false,
    description TEXT,
    image_url VARCHAR(255),
    amenities JSONB,
    accreditations JSONB,
    meta_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ======================================
-- 6) Create doctors (references facilities)
-- ======================================
CREATE TABLE doctors (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    specialty VARCHAR(100),
    primary_facility_id BIGINT REFERENCES facilities(id),
    contact_number VARCHAR(20),
    email VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ======================================
-- 7) Create reviews (generic reference via entity_type/entity_id)
-- ======================================
CREATE TABLE reviews (
    id BIGSERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id BIGINT NOT NULL,
    user_id BIGINT,
    rating DECIMAL(3,2) CHECK (rating BETWEEN 0 AND 5),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ======================================
-- 8) Create facility_departments
-- ======================================
CREATE TABLE facility_departments (
    id BIGSERIAL PRIMARY KEY,
    facility_id BIGINT REFERENCES facilities(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    floor_number VARCHAR(10),
    head_doctor_id BIGINT REFERENCES doctors(id),
    contact_number VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ======================================
-- 9) Create facility_equipment
-- ======================================
CREATE TABLE facility_equipment (
    id BIGSERIAL PRIMARY KEY,
    facility_id BIGINT REFERENCES facilities(id),
    department_id BIGINT REFERENCES facility_departments(id),
    name VARCHAR(200) NOT NULL,
    model VARCHAR(100),
    manufacturer VARCHAR(100),
    purchase_date DATE,
    last_maintenance_date DATE,
    next_maintenance_date DATE,
    status VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ======================================
-- 10) Create facility_operating_hours
-- ======================================
CREATE TABLE facility_operating_hours (
    id BIGSERIAL PRIMARY KEY,
    facility_id BIGINT REFERENCES facilities(id),
    department_id BIGINT REFERENCES facility_departments(id),
    day_of_week SMALLINT CHECK (day_of_week BETWEEN 0 AND 6),
    start_time TIME,
    end_time TIME,
    is_closed BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (facility_id, department_id, day_of_week)
);

-- ======================================
-- 11) Create facility_insurance_providers
-- ======================================
CREATE TABLE facility_insurance_providers (
    facility_id BIGINT REFERENCES facilities(id),
    insurance_provider_id BIGINT,
    coverage_details JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (facility_id, insurance_provider_id)
);

-- ======================================
-- 12) Create facility_certifications
-- ======================================
CREATE TABLE facility_certifications (
    id BIGSERIAL PRIMARY KEY,
    facility_id BIGINT REFERENCES facilities(id),
    name VARCHAR(200) NOT NULL,
    issuing_authority VARCHAR(200),
    issue_date DATE,
    expiry_date DATE,
    status VARCHAR(50),
    document_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ======================================
-- 13) Create indexes for facility-related tables
-- ======================================
CREATE INDEX idx_facilities_type_city ON facilities(type, city_id);
CREATE INDEX idx_facilities_coordinates ON facilities USING GIST(coordinates);
CREATE INDEX idx_facilities_rating_type ON facilities(rating DESC, type);
CREATE INDEX idx_facilities_amenities ON facilities USING GIN(amenities);
CREATE INDEX idx_facilities_meta_data ON facilities USING GIN(meta_data);

CREATE UNIQUE INDEX idx_facility_stats_facility_id ON facility_stats (facility_id);
CREATE INDEX idx_facility_departments_facility ON facility_departments(facility_id);
CREATE INDEX idx_facility_equipment_facility ON facility_equipment(facility_id, status);
CREATE INDEX idx_facility_operating_hours_facility ON facility_operating_hours(facility_id, day_of_week);
CREATE INDEX idx_facility_insurance_facility ON facility_insurance_providers(facility_id);
CREATE INDEX idx_facility_certifications_facility ON facility_certifications(facility_id, status);

-- ======================================
-- 14) Create materialized view for facility statistics
-- ======================================
CREATE MATERIALIZED VIEW facility_stats AS
SELECT 
    f.id AS facility_id,
    f.name AS facility_name,
    f.type,
    f.city_id,
    COUNT(DISTINCT fd.id) AS department_count,
    COUNT(DISTINCT d.id) AS doctor_count,
    COUNT(DISTINCT r.id) AS review_count,
    COALESCE(AVG(r.rating), 0) AS avg_rating,
    COUNT(DISTINCT fe.id) AS equipment_count,
    COUNT(DISTINCT fc.id) AS certification_count
FROM facilities f
LEFT JOIN facility_departments fd ON fd.facility_id = f.id
LEFT JOIN doctors d ON d.primary_facility_id = f.id
LEFT JOIN reviews r ON r.entity_type = 'facility' AND r.entity_id = f.id
LEFT JOIN facility_equipment fe ON fe.facility_id = f.id
LEFT JOIN facility_certifications fc ON fc.facility_id = f.id
GROUP BY f.id, f.name, f.type, f.city_id;

-- ======================================
-- 15) Create function to refresh materialized view
-- ======================================
CREATE OR REPLACE FUNCTION refresh_facility_stats()
RETURNS TRIGGER AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY facility_stats;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to refresh facility stats on changes to facilities
CREATE TRIGGER refresh_facility_stats_on_change
AFTER INSERT OR UPDATE OR DELETE ON facilities
FOR EACH STATEMENT EXECUTE FUNCTION refresh_facility_stats();

-- ======================================
-- 16) Create an audit logging system
-- ======================================
CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    table_name VARCHAR(50),
    operation VARCHAR(10),
    old_data JSONB,
    new_data JSONB,
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION log_audit()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO audit_log (table_name, operation, old_data, changed_at)
        VALUES (TG_TABLE_NAME, TG_OP, row_to_json(OLD), CURRENT_TIMESTAMP);
        RETURN OLD;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO audit_log (table_name, operation, new_data, changed_at)
        VALUES (TG_TABLE_NAME, TG_OP, row_to_json(NEW), CURRENT_TIMESTAMP);
        RETURN NEW;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO audit_log (table_name, operation, old_data, new_data, changed_at)
        VALUES (TG_TABLE_NAME, TG_OP, row_to_json(OLD), row_to_json(NEW), CURRENT_TIMESTAMP);
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for audit logging on key tables
CREATE TRIGGER audit_facilities
AFTER INSERT OR UPDATE OR DELETE ON facilities
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER audit_facility_departments
AFTER INSERT OR UPDATE OR DELETE ON facility_departments
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER audit_facility_equipment
AFTER INSERT OR UPDATE OR DELETE ON facility_equipment
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER audit_facility_operating_hours
AFTER INSERT OR UPDATE OR DELETE ON facility_operating_hours
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER audit_facility_insurance_providers
AFTER INSERT OR UPDATE OR DELETE ON facility_insurance_providers
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER audit_facility_certifications
AFTER INSERT OR UPDATE OR DELETE ON facility_certifications
FOR EACH ROW EXECUTE FUNCTION log_audit();

-- ======================================
-- 17) Create plans table
-- ======================================
CREATE TABLE plans (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    monthly_price DECIMAL(10, 2) NOT NULL,
    yearly_price DECIMAL(10, 2) NOT NULL,
    description TEXT NOT NULL,
    features JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ======================================
-- 18) Seed the plans table
--   (You may also do this in separate insert statements or via your app)
-- ======================================
INSERT INTO plans (name, monthly_price, yearly_price, description, features)
VALUES
(
    'Basic',
    39,            -- Monthly price
    390,           -- Yearly price
    'Essential features for small practices',
    '["Standard profile", "Basic search listing", "Basic analytics (views, clicks)"]'::jsonb
),
(
    'Pro',
    99,            -- Monthly price
    990,           -- Yearly price
    'Advanced features for growing clinics',
    '["Enhanced profile (photos and videos)", "Featured placement in search", "Advanced analytics dashboard"]'::jsonb
),
(
    'Enterprise',
    190,           -- Monthly price
    1900,          -- Yearly price
    'Comprehensive solution for large hospitals',
    '["AI-powered recommendations", "Real-time analytics dashboard", "Unlimited profile enhancements", "Priority support", "Free booking fees (unlimited appointments)"]'::jsonb
);

-- ======================================
-- 19) Create facility_plans table to link facilities to plans
-- ======================================
CREATE TABLE facility_plans (
    id BIGSERIAL PRIMARY KEY,
    facility_id BIGINT REFERENCES facilities(id) ON DELETE CASCADE,
    plan_id BIGINT REFERENCES plans(id) ON DELETE CASCADE,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_date TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- If you only want one active plan per facility, you could do:
    -- UNIQUE (facility_id) WHERE is_active
    --
    -- If multiple active plans are allowed, remove or modify that constraint.
    UNIQUE (facility_id, plan_id)
);

-- (Optional) Add triggers for auditing plan changes
CREATE TRIGGER audit_plans
AFTER INSERT OR UPDATE OR DELETE ON plans
FOR EACH ROW EXECUTE FUNCTION log_audit();

CREATE TRIGGER audit_facility_plans
AFTER INSERT OR UPDATE OR DELETE ON facility_plans
FOR EACH ROW EXECUTE FUNCTION log_audit();
