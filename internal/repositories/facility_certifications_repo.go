package repositories

import (
	"fmt"
	"strings"
	"time"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// FacilityCertificationsRepository defines CRUD operations for the FacilityCertifications model.
type FacilityCertificationsRepository interface {
	Repository[models.FacilityCertification]
}

// facilityCertificationsRepository is an implementation of FacilityCertificationsRepository.
type facilityCertificationsRepository struct {
	db *sqlx.DB
}

// NewFacilityCertificationsRepository initializes a new FacilityCertificationsRepository.
func NewFacilityCertificationsRepository(db *sqlx.DB) FacilityCertificationsRepository {
	return &facilityCertificationsRepository{db: db}
}

// Find fetches a facility certification by its ID.
func (r *facilityCertificationsRepository) Find(id int64) (*models.FacilityCertification, error) {
	start := time.Now()

	var certification models.FacilityCertification
	query := `SELECT * FROM facility_certifications WHERE id = $1`
	err := r.db.Get(&certification, query, id)

	trackMetrics("Find", "facility_certifications", start, err)

	if err != nil {
		return nil, err
	}
	return &certification, nil
}

// FindMany fetches facility certifications based on the filter.
func (r *facilityCertificationsRepository) FindMany(filter map[string]interface{}) ([]models.FacilityCertification, error) {
	start := time.Now()

	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := "SELECT * FROM facility_certifications"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var certifications []models.FacilityCertification
	err := r.db.Select(&certifications, query, args...)

	trackMetrics("FindMany", "facility_certifications", start, err)

	if err != nil {
		return nil, err
	}
	return certifications, nil
}

// Create adds a new facility certification.
func (r *facilityCertificationsRepository) Create(entity *models.FacilityCertification) (*models.FacilityCertification, error) {
	start := time.Now()

	query := `
		INSERT INTO facility_certifications (facility_id, name, issuing_authority, issue_date, expiry_date, status, document_url)
		VALUES (:facility_id, :name, :issuing_authority, :issue_date, :expiry_date, :status, :document_url)
		RETURNING *`

	rows, err := r.db.NamedQuery(query, entity)
	if err != nil {
		trackMetrics("Create", "facility_certifications", start, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(entity); err != nil {
			trackMetrics("Create", "facility_certifications", start, err)
			return nil, err
		}
	}
	trackMetrics("Create", "facility_certifications", start, nil)
	return entity, nil
}

// CreateMany adds multiple facility certifications.
func (r *facilityCertificationsRepository) CreateMany(entities []models.FacilityCertification) ([]models.FacilityCertification, error) {
	start := time.Now()

	tx, err := r.db.Beginx()
	if err != nil {
		trackMetrics("CreateMany", "facility_certifications", start, err)
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO facility_certifications (facility_id, name, issuing_authority, issue_date, expiry_date, status, document_url)
		VALUES (:facility_id, :name, :issuing_authority, :issue_date, :expiry_date, :status, :document_url)
		RETURNING *`

	var results []models.FacilityCertification
	for _, entity := range entities {
		rows, err := tx.NamedQuery(query, entity)
		if err != nil {
			trackMetrics("CreateMany", "facility_certifications", start, err)
			return nil, err
		}
		defer rows.Close()

		var result models.FacilityCertification
		if rows.Next() {
			if err := rows.StructScan(&result); err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	err = tx.Commit()
	trackMetrics("CreateMany", "facility_certifications", start, err)

	if err != nil {
		return nil, err
	}
	return results, nil
}

// Update modifies a facility certification by ID.
func (r *facilityCertificationsRepository) Update(id int64, updates map[string]interface{}) (*models.FacilityCertification, error) {
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
		UPDATE facility_certifications
		SET %s
		WHERE id = $%d
		RETURNING *`, strings.Join(setClauses, ", "), argIndex)

	var certification models.FacilityCertification
	err := r.db.QueryRowx(query, args...).StructScan(&certification)

	trackMetrics("Update", "facility_certifications", start, err)

	if err != nil {
		return nil, err
	}
	return &certification, nil
}

// UpdateMany modifies multiple facility certifications based on the filter.
func (r *facilityCertificationsRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
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
		UPDATE facility_certifications
		SET %s
		WHERE %s`, strings.Join(setClauses, ", "), strings.Join(whereClauses, " AND "))

	result, err := r.db.Exec(query, args...)
	trackMetrics("UpdateMany", "facility_certifications", start, err)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// Delete removes a facility certification by ID and returns the deleted row.
func (r *facilityCertificationsRepository) Delete(id int64) (*models.FacilityCertification, error) {
	start := time.Now()

	query := `
		DELETE FROM facility_certifications
		WHERE id = $1
		RETURNING *`

	var certification models.FacilityCertification
	err := r.db.QueryRowx(query, id).StructScan(&certification)

	trackMetrics("Delete", "facility_certifications", start, err)

	if err != nil {
		return nil, err
	}
	return &certification, nil
}

// DeleteMany removes multiple facility certifications based on the filter and returns the deleted rows.
func (r *facilityCertificationsRepository) DeleteMany(filter map[string]interface{}) ([]models.FacilityCertification, error) {
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
		DELETE FROM facility_certifications
		WHERE %s
		RETURNING *`, strings.Join(whereClauses, " AND "))

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		trackMetrics("DeleteMany", "facility_certifications", start, err)
		return nil, err
	}
	defer rows.Close()

	var deletedCertifications []models.FacilityCertification
	for rows.Next() {
		var certification models.FacilityCertification
		if err := rows.StructScan(&certification); err != nil {
			return nil, err
		}
		deletedCertifications = append(deletedCertifications, certification)
	}
	trackMetrics("DeleteMany", "facility_certifications", start, nil)
	return deletedCertifications, nil
}
