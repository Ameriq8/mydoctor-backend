package repositories

import (
	"fmt"
	"strings"
	"time"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// FacilityAppointmentRepository defines CRUD operations for the FacilityAppointment model.
type FacilityAppointmentRepository interface {
	Repository[models.FacilityAppointment]
}

// facilityAppointmentRepository is an implementation of FacilityAppointmentRepository.
type facilityAppointmentRepository struct {
	db *sqlx.DB
}

// NewFacilityAppointmentRepository initializes a new FacilityAppointmentRepository.
func NewFacilityAppointmentRepository(db *sqlx.DB) FacilityAppointmentRepository {
	return &facilityAppointmentRepository{db: db}
}

// Find fetches a facility appointment by its ID.
func (r *facilityAppointmentRepository) Find(id int64) (*models.FacilityAppointment, error) {
	start := time.Now()

	var appointment models.FacilityAppointment
	query := `SELECT * FROM facility_appointments WHERE id = $1`
	err := r.db.Get(&appointment, query, id)

	trackMetrics("Find", "facility_appointments", start, err)

	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

// FindMany fetches facility appointments based on the filter.
func (r *facilityAppointmentRepository) FindMany(filter map[string]interface{}) ([]models.FacilityAppointment, error) {
	start := time.Now()

	var appointments []models.FacilityAppointment
	query := "SELECT * FROM facility_appointments"
	err := r.db.Select(&appointments, query)

	trackMetrics("FindMany", "facility_appointments", start, err)

	if err != nil {
		return nil, err
	}
	return appointments, nil
}

// Create adds a new facility appointment.
func (r *facilityAppointmentRepository) Create(entity *models.FacilityAppointment) (*models.FacilityAppointment, error) {
	start := time.Now()

	query := `
		INSERT INTO facility_appointments (patient_name, patient_contact, facility_id, doctor_id, appointment_time, status, reason_for_appointment)
		VALUES (:patient_name, :patient_contact, :facility_id, :doctor_id, :appointment_time, :status, :reason_for_appointment)
		RETURNING *
	`
	rows, err := r.db.NamedQuery(query, entity)
	if err != nil {
		trackMetrics("Create", "facility_appointments", start, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(entity); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	trackMetrics("Create", "facility_appointments", start, nil)
	return entity, nil
}

// CreateMany adds multiple facility appointments.
func (r *facilityAppointmentRepository) CreateMany(entities []models.FacilityAppointment) ([]models.FacilityAppointment, error) {
	start := time.Now()

	query := `
    INSERT INTO facility_appointments (patient_name, patient_contact, facility_id, doctor_id, appointment_time, status, reason_for_appointment)
    VALUES (:patient_name, :patient_contact, :facility_id, :doctor_id, :appointment_time, :status, :reason_for_appointment)
    RETURNING *;`
	tx, err := r.db.Beginx()
	if err != nil {
		trackMetrics("CreateMany", "facility_appointments", start, err)
		return nil, err
	}
	defer tx.Rollback()

	var results []models.FacilityAppointment
	for _, entity := range entities {
		rows, err := r.db.NamedQuery(query, entity)
		if err != nil {
			trackMetrics("CreateMany", "facility_appointments", start, err)
			return nil, err
		}
		defer rows.Close()

		var result models.FacilityAppointment
		if rows.Next() {
			if err := rows.StructScan(&result); err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	if err := tx.Commit(); err != nil {
		trackMetrics("CreateMany", "facility_appointments", start, err)
		return nil, err
	}

	trackMetrics("CreateMany", "facility_appointments", start, nil)
	return results, nil
}

// Update modifies a facility appointment by ID.
func (r *facilityAppointmentRepository) Update(id int64, updates map[string]interface{}) (*models.FacilityAppointment, error) {
	start := time.Now()

	setClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}
	args = append(args, id)

	query := fmt.Sprintf(`UPDATE facility_appointments SET %s WHERE id = $%d RETURNING *`, strings.Join(setClauses, ", "), argIndex)

	var appointment models.FacilityAppointment
	err := r.db.QueryRowx(query, args...).StructScan(&appointment)
	trackMetrics("Update", "facility_appointments", start, err)

	if err != nil {
		return nil, err
	}

	return &appointment, nil
}

// UpdateMany modifies multiple facility appointments based on the filter.
func (r *facilityAppointmentRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
	start := time.Now()

	setClauses := []string{}
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf(`
		UPDATE facility_appointments
		SET %s
		WHERE %s
	`, strings.Join(setClauses, ", "), strings.Join(whereClauses, " AND "))

	result, err := r.db.Exec(query, args...)
	trackMetrics("UpdateMany", "facility_appointments", start, err)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// Delete removes a facility appointment by ID and returns the deleted row.
func (r *facilityAppointmentRepository) Delete(id int64) (*models.FacilityAppointment, error) {
	start := time.Now()

	query := `
        DELETE FROM facility_appointments
        WHERE id = $1
        RETURNING id, patient_name, patient_contact, facility_id, doctor_id, appointment_time, status, reason_for_appointment
    `

	var appointment models.FacilityAppointment
	err := r.db.QueryRowx(query, id).StructScan(&appointment)
	trackMetrics("Delete", "facility_appointments", start, err)

	if err != nil {
		return nil, err
	}

	return &appointment, nil
}

// DeleteMany removes multiple facility appointments based on the filter and returns the deleted rows.
func (r *facilityAppointmentRepository) DeleteMany(filter map[string]interface{}) ([]models.FacilityAppointment, error) {
	start := time.Now()

	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf(`
        DELETE FROM facility_appointments
        WHERE %s
        RETURNING id, patient_name, patient_contact, facility_id, doctor_id, appointment_time, status, reason_for_appointment
    `, strings.Join(whereClauses, " AND "))

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		trackMetrics("DeleteMany", "facility_appointments", start, err)
		return nil, err
	}
	defer rows.Close()

	var deletedAppointments []models.FacilityAppointment
	for rows.Next() {
		var appointment models.FacilityAppointment
		if err := rows.StructScan(&appointment); err != nil {
			return nil, err
		}
		deletedAppointments = append(deletedAppointments, appointment)
	}

	trackMetrics("DeleteMany", "facility_appointments", start, nil)
	return deletedAppointments, nil
}
