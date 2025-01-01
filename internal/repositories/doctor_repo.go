package repositories

import (
	"fmt"
	"strings"
	"time"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// DoctorRepository defines CRUD operations for the Doctor model.
type DoctorRepository interface {
	Repository[models.Doctor]
}

// doctorRepository is an implementation of DoctorRepository.
type doctorRepository struct {
	db *sqlx.DB
}

// NewDoctorRepository initializes a new DoctorRepository.
func NewDoctorRepository(db *sqlx.DB) DoctorRepository {
	return &doctorRepository{db: db}
}

// Find fetches a doctor by its ID.
func (r *doctorRepository) Find(id int64) (*models.Doctor, error) {
	start := time.Now()

	var doctor models.Doctor
	query := `SELECT * FROM doctors WHERE id = $1`
	err := r.db.Get(&doctor, query, id)

	trackMetrics("Find", "doctors", start, err)

	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

// FindMany fetches doctors based on the filter.
func (r *doctorRepository) FindMany(filter map[string]interface{}) ([]models.Doctor, error) {
	start := time.Now()

	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := "SELECT * FROM doctors"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var doctors []models.Doctor
	err := r.db.Select(&doctors, query, args...)

	trackMetrics("FindMany", "doctors", start, err)

	if err != nil {
		return nil, err
	}
	return doctors, nil
}

// Create adds a new doctor.
func (r *doctorRepository) Create(entity *models.Doctor) (*models.Doctor, error) {
	start := time.Now()

	query := `
		INSERT INTO doctors (name, specialty, primary_facility_id, contact_number, email)
		VALUES (:name, :specialty, :primary_facility_id, :contact_number, :email)
		RETURNING *`

	rows, err := r.db.NamedQuery(query, entity)
	if err != nil {
		trackMetrics("Create", "doctors", start, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(entity); err != nil {
			trackMetrics("Create", "doctors", start, err)
			return nil, err
		}
	}
	trackMetrics("Create", "doctors", start, nil)
	return entity, nil
}

// CreateMany adds multiple doctors.
func (r *doctorRepository) CreateMany(entities []models.Doctor) ([]models.Doctor, error) {
	start := time.Now()

	tx, err := r.db.Beginx()
	if err != nil {
		trackMetrics("CreateMany", "doctors", start, err)
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO doctors (name, specialty, primary_facility_id, contact_number, email)
		VALUES (:name, :specialty, :primary_facility_id, :contact_number, :email)
		RETURNING *`

	var results []models.Doctor
	for _, entity := range entities {
		rows, err := tx.NamedQuery(query, entity)
		if err != nil {
			trackMetrics("CreateMany", "doctors", start, err)
			return nil, err
		}
		defer rows.Close()

		var result models.Doctor
		if rows.Next() {
			if err := rows.StructScan(&result); err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	err = tx.Commit()
	trackMetrics("CreateMany", "doctors", start, err)

	if err != nil {
		return nil, err
	}
	return results, nil
}

// Update modifies a doctor by ID.
func (r *doctorRepository) Update(id int64, updates map[string]interface{}) (*models.Doctor, error) {
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
		UPDATE doctors
		SET %s
		WHERE id = $%d
		RETURNING *`, strings.Join(setClauses, ", "), argIndex)

	var doctor models.Doctor
	err := r.db.QueryRowx(query, args...).StructScan(&doctor)

	trackMetrics("Update", "doctors", start, err)

	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

// UpdateMany modifies multiple doctors based on the filter.
func (r *doctorRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
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
		UPDATE doctors
		SET %s
		WHERE %s`, strings.Join(setClauses, ", "), strings.Join(whereClauses, " AND "))

	result, err := r.db.Exec(query, args...)
	trackMetrics("UpdateMany", "doctors", start, err)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// Delete removes a doctor by ID and returns the deleted row.
func (r *doctorRepository) Delete(id int64) (*models.Doctor, error) {
	start := time.Now()

	query := `
		DELETE FROM doctors
		WHERE id = $1
		RETURNING *`

	var doctor models.Doctor
	err := r.db.QueryRowx(query, id).StructScan(&doctor)

	trackMetrics("Delete", "doctors", start, err)

	if err != nil {
		return nil, err
	}
	return &doctor, nil
}

// DeleteMany removes multiple doctors based on the filter and returns the deleted rows.
func (r *doctorRepository) DeleteMany(filter map[string]interface{}) ([]models.Doctor, error) {
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
		DELETE FROM doctors
		WHERE %s
		RETURNING *`, strings.Join(whereClauses, " AND "))

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		trackMetrics("DeleteMany", "doctors", start, err)
		return nil, err
	}
	defer rows.Close()

	var deletedDoctors []models.Doctor
	for rows.Next() {
		var doctor models.Doctor
		if err := rows.StructScan(&doctor); err != nil {
			return nil, err
		}
		deletedDoctors = append(deletedDoctors, doctor)
	}
	trackMetrics("DeleteMany", "doctors", start, nil)
	return deletedDoctors, nil
}
