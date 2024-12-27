-- ======================================
-- 1) Drop extensions
-- ======================================
DROP EXTENSION IF EXISTS "uuid-ossp" CASCADE;
DROP EXTENSION IF EXISTS btree_gist CASCADE;
DROP EXTENSION IF EXISTS pg_trgm CASCADE;

-- ======================================
-- 2) Drop ENUM type for facility_type
-- ======================================
DROP TYPE IF EXISTS facility_type CASCADE;

-- ======================================
-- 3) Drop materialized view
-- ======================================
DROP MATERIALIZED VIEW IF EXISTS facility_stats CASCADE;

-- ======================================
-- 4) Drop functions
-- ======================================
DROP FUNCTION IF EXISTS refresh_facility_stats CASCADE;
DROP FUNCTION IF EXISTS archive_old_facility_metrics CASCADE;
DROP FUNCTION IF EXISTS log_audit CASCADE;

-- ======================================
-- 5) Drop tables
-- ======================================
DROP TABLE IF EXISTS audit_log CASCADE;
DROP TABLE IF EXISTS facility_plans CASCADE;
DROP TABLE IF EXISTS plans CASCADE;
DROP TABLE IF EXISTS facility_metrics_q2_2024 CASCADE;
DROP TABLE IF EXISTS facility_metrics_q1_2024 CASCADE;
DROP TABLE IF EXISTS facility_metrics CASCADE;
DROP TABLE IF EXISTS facility_certifications CASCADE;
DROP TABLE IF EXISTS facility_insurance_providers CASCADE;
DROP TABLE IF EXISTS facility_operating_hours CASCADE;
DROP TABLE IF EXISTS facility_equipment CASCADE;
DROP TABLE IF EXISTS facility_departments CASCADE;
DROP TABLE IF EXISTS reviews CASCADE;
DROP TABLE IF EXISTS doctors CASCADE;
DROP TABLE IF EXISTS facilities CASCADE;
DROP TABLE IF EXISTS cities CASCADE;
DROP TABLE IF EXISTS facility_categories CASCADE;

-- ======================================
-- 6) Drop indexes (optional, will be removed with tables)
-- ======================================
-- No separate DROP INDEX statements are needed as they are removed with tables.

-- ======================================
-- 7) Drop triggers
-- ======================================
DROP TRIGGER IF EXISTS refresh_facility_stats_on_change ON facilities CASCADE;
DROP TRIGGER IF EXISTS audit_facilities ON facilities CASCADE;
DROP TRIGGER IF EXISTS audit_facility_departments ON facility_departments CASCADE;
DROP TRIGGER IF EXISTS audit_facility_equipment ON facility_equipment CASCADE;
DROP TRIGGER IF EXISTS audit_facility_operating_hours ON facility_operating_hours CASCADE;
DROP TRIGGER IF EXISTS audit_facility_insurance_providers ON facility_insurance_providers CASCADE;
DROP TRIGGER IF EXISTS audit_facility_certifications ON facility_certifications CASCADE;
DROP TRIGGER IF EXISTS audit_plans ON plans CASCADE;
DROP TRIGGER IF EXISTS audit_facility_plans ON facility_plans CASCADE;
