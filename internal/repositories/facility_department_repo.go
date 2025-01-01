package repositories

import (
	"fmt"
	"strings"
	"time"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// FacilityDepartmentRepository defines CRUD operations for the FacilityDepartment model.
type FacilityDepartmentRepository interface {
	Repository[models.FacilityDepartment]
}

// facilityDepartmentRepository is an implementation of FacilityDepartmentRepository.
type facilityDepartmentRepository struct {
	db *sqlx.DB
}

// NewFacilityDepartmentRepository initializes a new FacilityDepartmentRepository.
func NewFacilityDepartmentRepository(db *sqlx.DB) FacilityDepartmentRepository {
	return &facilityDepartmentRepository{db: db}
}

// Find fetches a facility department by its ID.
func (r *facilityDepartmentRepository) Find(id int64) (*models.FacilityDepartment, error) {
	start := time.Now()

	var department models.FacilityDepartment
	query := `SELECT * FROM facility_departments WHERE id = $1`
	err := r.db.Get(&department, query, id)

	trackMetrics("Find", "facility_departments", start, err)

	if err != nil {
		return nil, err
	}
	return &department, nil
}

// FindMany fetches facility departments based on the filter.
func (r *facilityDepartmentRepository) FindMany(filter map[string]interface{}) ([]models.FacilityDepartment, error) {
	start := time.Now()

	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := "SELECT * FROM facility_departments"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var departments []models.FacilityDepartment
	err := r.db.Select(&departments, query, args...)

	trackMetrics("FindMany", "facility_departments", start, err)

	if err != nil {
		return nil, err
	}
	return departments, nil
}

// Create adds a new facility department.
func (r *facilityDepartmentRepository) Create(entity *models.FacilityDepartment) (*models.FacilityDepartment, error) {
	start := time.Now()

	query := `
		INSERT INTO facility_departments (facility_id, name, description, floor_number, head_doctor_id, contact_number)
		VALUES (:facility_id, :name, :description, :floor_number, :head_doctor_id, :contact_number)
		RETURNING *`

	rows, err := r.db.NamedQuery(query, entity)
	if err != nil {
		trackMetrics("Create", "facility_departments", start, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(entity); err != nil {
			trackMetrics("Create", "facility_departments", start, err)
			return nil, err
		}
	}
	trackMetrics("Create", "facility_departments", start, nil)
	return entity, nil
}

// CreateMany adds multiple facility departments.
func (r *facilityDepartmentRepository) CreateMany(entities []models.FacilityDepartment) ([]models.FacilityDepartment, error) {
	start := time.Now()

	tx, err := r.db.Beginx()
	if err != nil {
		trackMetrics("CreateMany", "facility_departments", start, err)
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO facility_departments (facility_id, name, description, floor_number, head_doctor_id, contact_number)
		VALUES (:facility_id, :name, :description, :floor_number, :head_doctor_id, :contact_number)
		RETURNING *`

	var results []models.FacilityDepartment
	for _, entity := range entities {
		rows, err := tx.NamedQuery(query, entity)
		if err != nil {
			trackMetrics("CreateMany", "facility_departments", start, err)
			return nil, err
		}
		defer rows.Close()

		var result models.FacilityDepartment
		if rows.Next() {
			if err := rows.StructScan(&result); err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	err = tx.Commit()
	trackMetrics("CreateMany", "facility_departments", start, err)

	if err != nil {
		return nil, err
	}
	return results, nil
}

// Update modifies a facility department by ID.
func (r *facilityDepartmentRepository) Update(id int64, updates map[string]interface{}) (*models.FacilityDepartment, error) {
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

	query := fmt.Sprintf(`
		UPDATE facility_departments
		SET %s
		WHERE id = $%d
		RETURNING *`, strings.Join(setClauses, ", "), argIndex)

	var department models.FacilityDepartment
	err := r.db.QueryRowx(query, args...).StructScan(&department)

	trackMetrics("Update", "facility_departments", start, err)

	if err != nil {
		return nil, err
	}
	return &department, nil
}

// UpdateMany modifies multiple facility departments based on the filter.
func (r *facilityDepartmentRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
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
		UPDATE facility_departments
		SET %s
		WHERE %s`, strings.Join(setClauses, ", "), strings.Join(whereClauses, " AND "))

	result, err := r.db.Exec(query, args...)
	trackMetrics("UpdateMany", "facility_departments", start, err)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// Delete removes a facility department by ID and returns the deleted row.
func (r *facilityDepartmentRepository) Delete(id int64) (*models.FacilityDepartment, error) {
	start := time.Now()

	query := `
		DELETE FROM facility_departments
		WHERE id = $1
		RETURNING *`

	var department models.FacilityDepartment
	err := r.db.QueryRowx(query, id).StructScan(&department)

	trackMetrics("Delete", "facility_departments", start, err)

	if err != nil {
		return nil, err
	}
	return &department, nil
}

// DeleteMany removes multiple facility departments based on the filter and returns the deleted rows.
func (r *facilityDepartmentRepository) DeleteMany(filter map[string]interface{}) ([]models.FacilityDepartment, error) {
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
		DELETE FROM facility_departments
		WHERE %s
		RETURNING *`, strings.Join(whereClauses, " AND "))

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		trackMetrics("DeleteMany", "facility_departments", start, err)
		return nil, err
	}
	defer rows.Close()

	var deletedDepartments []models.FacilityDepartment
	for rows.Next() {
		var department models.FacilityDepartment
		if err := rows.StructScan(&department); err != nil {
			return nil, err
		}
		deletedDepartments = append(deletedDepartments, department)
	}
	trackMetrics("DeleteMany", "facility_departments", start, nil)
	return deletedDepartments, nil
}
