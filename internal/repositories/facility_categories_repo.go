package repositories

import (
	"fmt"
	"strings"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// FacilityCategoriesRepository defines CRUD operations for the FacilityCategories model.
type FacilityCategoriesRepository interface {
	Repository[models.FacilityCategory]
}

// facilityCategoriesRepository is an implementation of FacilityCategoriesRepository.
type facilityCategoriesRepository struct {
	db *sqlx.DB
}

// NewFacilityCategoriesRepository initializes a new FacilityCategoriesRepository.
func NewFacilityCategoriesRepository(db *sqlx.DB) FacilityCategoriesRepository {
	return &facilityCategoriesRepository{db: db}
}

// Find fetches a facility category by its ID.
func (r *facilityCategoriesRepository) Find(id int64) (*models.FacilityCategory, error) {
	var category models.FacilityCategory
	query := `SELECT * FROM facility_categories WHERE id = $1`
	if err := r.db.Get(&category, query, id); err != nil {
		return nil, err
	}
	return &category, nil
}

// FindMany fetches facility categories based on the filter.
func (r *facilityCategoriesRepository) FindMany(filter map[string]interface{}) ([]models.FacilityCategory, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := "SELECT * FROM facility_categories"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var categories []models.FacilityCategory
	if err := r.db.Select(&categories, query, args...); err != nil {
		return nil, err
	}
	return categories, nil
}

// Create adds a new facility category.
func (r *facilityCategoriesRepository) Create(entity *models.FacilityCategory) (*models.FacilityCategory, error) {
	query := `
		INSERT INTO facility_categories (name, description, parent_id)
		VALUES (:name, :description, :parent_id)
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

// CreateMany adds multiple facility categories.
func (r *facilityCategoriesRepository) CreateMany(entities []models.FacilityCategory) ([]models.FacilityCategory, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO facility_categories (name, description, parent_id)
		VALUES (:name, :description, :parent_id)
		RETURNING *`

	var results []models.FacilityCategory
	for _, entity := range entities {
		rows, err := tx.NamedQuery(query, entity)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var result models.FacilityCategory
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

// Update modifies a facility category by ID.
func (r *facilityCategoriesRepository) Update(id int64, updates map[string]interface{}) (*models.FacilityCategory, error) {
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
		UPDATE facility_categories
		SET %s
		WHERE id = $%d
		RETURNING *`, strings.Join(setClauses, ", "), argIndex)

	var category models.FacilityCategory
	if err := r.db.QueryRowx(query, args...).StructScan(&category); err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateMany modifies multiple facility categories based on the filter.
func (r *facilityCategoriesRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
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
		UPDATE facility_categories
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

// Delete removes a facility category by ID and returns the deleted row.
func (r *facilityCategoriesRepository) Delete(id int64) (*models.FacilityCategory, error) {
	query := `
		DELETE FROM facility_categories
		WHERE id = $1
		RETURNING *`

	var category models.FacilityCategory
	if err := r.db.QueryRowx(query, id).StructScan(&category); err != nil {
		return nil, err
	}
	return &category, nil
}

// DeleteMany removes multiple facility categories based on the filter and returns the deleted rows.
func (r *facilityCategoriesRepository) DeleteMany(filter map[string]interface{}) ([]models.FacilityCategory, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf(`
		DELETE FROM facility_categories
		WHERE %s
		RETURNING *`, strings.Join(whereClauses, " AND "))

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deletedCategories []models.FacilityCategory
	for rows.Next() {
		var category models.FacilityCategory
		if err := rows.StructScan(&category); err != nil {
			return nil, err
		}
		deletedCategories = append(deletedCategories, category)
	}
	return deletedCategories, nil
}
