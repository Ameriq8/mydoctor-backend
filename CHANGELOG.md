# Changelog

## [Unreleased]

### Refactor and Add New Modules for Facility Handling and Validation
- **Modified** `internal/handlers/facility_handler.go` to improve the facility data handling logic.
- **Added** `internal/validators/facility_validator.go` for introducing validation logic for facility data.
- **Created** `pkg/middlewares/auth.go` for implementing authentication middleware.
- **Created** `pkg/middlewares/validation.go` for providing validation middleware for incoming requests.
- **Added** `pkg/utils/errors.go` for centralized error handling utilities.
- **Updated** `pkg/utils/pg.go` to enhance PostgreSQL connection handling.
- **Improved** modularity, validation, and error handling across the application.

### Implement Modular Handlers and Services
- **Added** `city_handler.go` for city-related route handling.
- **Added** `facility_handler.go` for facility-related route handling.
- **Introduced** `handler_initializer.go` to centralize route registration.
- **Implemented** `city_service.go` for city-specific business logic.
- **Implemented** `facility_service.go` for facility-specific business logic.
- **Updated** `main.go` to integrate new handler and service structure.

### Add and Update Configurations for Monitoring and Testing
- **Added** `alert_rules.yml` for defining Prometheus alerting rules.
- **Added** `alertmanager.yml` for configuring Alertmanager.
- **Added** `artillery-config.yml` for performance testing scenarios.
- **Updated** `docker-compose.yml` to include new services and dependencies.
- **Updated** `Dockerfile` for enhanced container configuration.
- **Updated** `prometheus.yml` with new scrape configurations and metrics endpoints.
- **Improved** `README.md` with detailed instructions for running and testing.
- **Updated** `main.go` to enhance monitoring and API functionality.

### Enhanced Prometheus Metrics
- **Added** in-flight query tracking.
- **Optimized** query duration buckets for better performance monitoring.

### Integrate trackMetrics for All Repository Operations
- **Integrated** `trackMetrics` into all repository operations for better performance tracking and monitoring.

### Set Up Prometheus and Docker
- **Configured** Prometheus and Docker for monitoring and containerized environments.

### Implement Logging Functions
- **Implemented** logging functions for better traceability and debugging.

### Add All Repo Files
- **Added** all repository files and related configurations for enhanced code structure and organization.

### Add Facility Repository File to Repositories Directory
- **Added** a facility repository file to the repositories directory along with new changes to improve data access.

### Refactor Database Configuration and Connection Handling
- **Introduced** a new `LoadedConfig` struct and `SQLRepository` interface for improved maintainability and flexibility in database configuration and connections.

### Enhance Database Schema with Detailed Comments
- **Improved** organization and maintainability of SQL scripts.
- **Added** detailed comments to the database schema for clarity.

### Refactor Models to Use BaseModel for Common Fields
- **Refactored** models to use `BaseModel` for common fields to streamline code and improve maintainability.

### Add Logger Package, Middleware for Logging, and Initial API Routes
- **Added** a logger package and middleware for logging requests.
- **Updated** API routes and added initial implementation.
- **Updated** README with additional setup and usage instructions.

---

## [Version 1.0.0] - 2025-01-02

### Initial Release
- Implemented the core functionality and setup for the application.
- Added initial modules for user authentication, validation, error handling, and database configuration.
