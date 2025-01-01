package repositories

import (
	"fmt"
	"strings"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// FacilityInsuranceProvidersRepository defines CRUD operations for the FacilityInsuranceProviders model.
type FacilityInsuranceProvidersRepository interface {
	Repository[models.FacilityInsuranceProvider]
}

// facilityInsuranceProvidersRepository is an implementation of FacilityInsuranceProvidersRepository.
type facilityInsuranceProvidersRepository struct {
	db *sqlx.DB
}

// NewFacilityInsuranceProvidersRepository initializes a new FacilityInsuranceProvidersRepository.
func NewFacilityInsuranceProvidersRepository(db *sqlx.DB) FacilityInsuranceProvidersRepository {
	return &facilityInsuranceProvidersRepository{db: db}
}

// Find fetches a facility insurance provider record by facility_id and insurance_provider_id.
func (r *facilityInsuranceProvidersRepository) Find(id int64) (*models.FacilityInsuranceProvider, error) {
	query := `
		SELECT * FROM facility_insurance_providers
		WHERE id = $1`

	var provider models.FacilityInsuranceProvider
	if err := r.db.Get(&provider, query, id); err != nil {
		return nil, err
	}
	return &provider, nil
}

// FindMany fetches facility insurance providers based on the filter.
func (r *facilityInsuranceProvidersRepository) FindMany(filter map[string]interface{}) ([]models.FacilityInsuranceProvider, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := "SELECT * FROM facility_insurance_providers"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var providers []models.FacilityInsuranceProvider
	if err := r.db.Select(&providers, query, args...); err != nil {
		return nil, err
	}
	return providers, nil
}

// Create adds a new facility insurance provider.
func (r *facilityInsuranceProvidersRepository) Create(entity *models.FacilityInsuranceProvider) (*models.FacilityInsuranceProvider, error) {
	query := `
		INSERT INTO facility_insurance_providers (facility_id, insurance_provider_id, coverage_details)
		VALUES (:facility_id, :insurance_provider_id, :coverage_details)
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

// CreateMany adds multiple facility insurance provider records.
func (r *facilityInsuranceProvidersRepository) CreateMany(entities []models.FacilityInsuranceProvider) ([]models.FacilityInsuranceProvider, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO facility_insurance_providers (facility_id, insurance_provider_id, coverage_details)
		VALUES (:facility_id, :insurance_provider_id, :coverage_details)
		RETURNING *`

	var results []models.FacilityInsuranceProvider
	for _, entity := range entities {
		rows, err := tx.NamedQuery(query, entity)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var result models.FacilityInsuranceProvider
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

// Update modifies a facility insurance provider record.
func (r *facilityInsuranceProvidersRepository) Update(id int64, updates map[string]interface{}) (*models.FacilityInsuranceProvider, error) {
	setClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	// Build SET clause
	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE facility_insurance_providers
		SET %s
		WHERE id = $%d
		RETURNING *`, strings.Join(setClauses, ", "), argIndex)

	var provider models.FacilityInsuranceProvider
	if err := r.db.QueryRowx(query, args...).StructScan(&provider); err != nil {
		return nil, err
	}
	return &provider, nil
}

// UpdateMany modifies multiple facility insurance provider records based on the filter.
func (r *facilityInsuranceProvidersRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
	setClauses := []string{}
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	// Build SET clause
	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	// Build WHERE clause
	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf(`
		UPDATE facility_insurance_providers
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

// Delete removes a facility insurance provider record and returns the deleted row.
func (r *facilityInsuranceProvidersRepository) Delete(id int64) (*models.FacilityInsuranceProvider, error) {
	query := `
		DELETE FROM facility_insurance_providers
		WHERE id = $1
		RETURNING *`

	var provider models.FacilityInsuranceProvider
	if err := r.db.QueryRowx(query, id).StructScan(&provider); err != nil {
		return nil, err
	}
	return &provider, nil
}

// DeleteMany removes multiple facility insurance provider records based on the filter and returns the deleted rows.
func (r *facilityInsuranceProvidersRepository) DeleteMany(filter map[string]interface{}) ([]models.FacilityInsuranceProvider, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf(`
		DELETE FROM facility_insurance_providers
		WHERE %s
		RETURNING *`, strings.Join(whereClauses, " AND "))

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deletedProviders []models.FacilityInsuranceProvider
	for rows.Next() {
		var provider models.FacilityInsuranceProvider
		if err := rows.StructScan(&provider); err != nil {
			return nil, err
		}
		deletedProviders = append(deletedProviders, provider)
	}
	return deletedProviders, nil
}
