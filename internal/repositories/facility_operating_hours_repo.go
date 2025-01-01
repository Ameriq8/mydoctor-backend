package repositories

import (
	"fmt"
	"strings"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// FacilityOperatingHoursRepository defines CRUD operations for the FacilityOperatingHours model.
type FacilityOperatingHoursRepository interface {
	Repository[models.FacilityOperatingHours]
}

// facilityOperatingHoursRepository is an implementation of FacilityOperatingHoursRepository.
type facilityOperatingHoursRepository struct {
	db *sqlx.DB
}

// NewFacilityOperatingHoursRepository initializes a new FacilityOperatingHoursRepository.
func NewFacilityOperatingHoursRepository(db *sqlx.DB) FacilityOperatingHoursRepository {
	return &facilityOperatingHoursRepository{db: db}
}

// Find fetches a facility operating hours record by its ID.
func (r *facilityOperatingHoursRepository) Find(id int64) (*models.FacilityOperatingHours, error) {
	var hours models.FacilityOperatingHours
	query := `SELECT * FROM facility_operating_hours WHERE id = $1`
	if err := r.db.Get(&hours, query, id); err != nil {
		return nil, err
	}

	return &hours, nil
}

// FindMany fetches facility operating hours based on the filter.
func (r *facilityOperatingHoursRepository) FindMany(filter map[string]interface{}) ([]models.FacilityOperatingHours, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := "SELECT * FROM facility_operating_hours"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var hours []models.FacilityOperatingHours
	if err := r.db.Select(&hours, query, args...); err != nil {
		return nil, err
	}

	return hours, nil
}

// Create adds a new facility operating hours record.
func (r *facilityOperatingHoursRepository) Create(entity *models.FacilityOperatingHours) (*models.FacilityOperatingHours, error) {
	query := `
		INSERT INTO facility_operating_hours (facility_id, department_id, day_of_week, start_time, end_time, is_closed)
		VALUES (:facility_id, :department_id, :day_of_week, :start_time, :end_time, :is_closed)
		RETURNING *`

	rows, err := r.db.NamedQuery(query, entity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(entity); err != nil {
			return nil, err
		}
	}

	return entity, nil
}

// CreateMany adds multiple facility operating hours records.
func (r *facilityOperatingHoursRepository) CreateMany(entities []models.FacilityOperatingHours) ([]models.FacilityOperatingHours, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO facility_operating_hours (facility_id, department_id, day_of_week, start_time, end_time, is_closed)
		VALUES (:facility_id, :department_id, :day_of_week, :start_time, :end_time, :is_closed)
		RETURNING *`

	var results []models.FacilityOperatingHours
	for _, entity := range entities {
		rows, err := tx.NamedQuery(query, entity)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var result models.FacilityOperatingHours
		if rows.Next() {
			if err := rows.StructScan(&result); err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return results, nil
}

// Update modifies a facility operating hours record by ID.
func (r *facilityOperatingHoursRepository) Update(id int64, updates map[string]interface{}) (*models.FacilityOperatingHours, error) {
	setClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE facility_operating_hours
		SET %s
		WHERE id = $%d
		RETURNING *`, strings.Join(setClauses, ", "), argIndex)

	var hours models.FacilityOperatingHours
	if err := r.db.QueryRowx(query, args...).StructScan(&hours); err != nil {
		return nil, err
	}

	return &hours, nil
}

// UpdateMany modifies multiple facility operating hours records based on the filter.
func (r *facilityOperatingHoursRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
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
		UPDATE facility_operating_hours
		SET %s
		WHERE %s`, strings.Join(setClauses, ", "), strings.Join(whereClauses, " AND "))

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// Delete removes a facility operating hours record by ID and returns the deleted row.
func (r *facilityOperatingHoursRepository) Delete(id int64) (*models.FacilityOperatingHours, error) {
	query := `
		DELETE FROM facility_operating_hours
		WHERE id = $1
		RETURNING *`

	var hours models.FacilityOperatingHours
	if err := r.db.QueryRowx(query, id).StructScan(&hours); err != nil {
		return nil, err
	}

	return &hours, nil
}

// DeleteMany removes multiple facility operating hours records based on the filter and returns the deleted rows.
func (r *facilityOperatingHoursRepository) DeleteMany(filter map[string]interface{}) ([]models.FacilityOperatingHours, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf(`
		DELETE FROM facility_operating_hours
		WHERE %s
		RETURNING *`, strings.Join(whereClauses, " AND "))

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deletedHours []models.FacilityOperatingHours
	for rows.Next() {
		var hours models.FacilityOperatingHours
		if err := rows.StructScan(&hours); err != nil {
			return nil, err
		}
		deletedHours = append(deletedHours, hours)
	}

	return deletedHours, nil
}
