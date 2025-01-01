package repositories

import (
	"fmt"
	"strings"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// CitiesRepository defines CRUD operations for the Cities model.
type CitiesRepository interface {
	Repository[models.City]
}

// citiesRepository is an implementation of CitiesRepository.
type citiesRepository struct {
	db *sqlx.DB
}

// NewCitiesRepository initializes a new CitiesRepository.
func NewCitiesRepository(db *sqlx.DB) CitiesRepository {
	return &citiesRepository{db: db}
}

// Find fetches a city by its ID.
func (r *citiesRepository) Find(id int64) (*models.City, error) {
	var city models.City
	query := `SELECT * FROM cities WHERE id = $1`
	if err := r.db.Get(&city, query, id); err != nil {
		return nil, err
	}
	return &city, nil
}

// FindMany fetches cities based on the filter.
func (r *citiesRepository) FindMany(filter map[string]interface{}) ([]models.City, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := "SELECT * FROM cities"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var cities []models.City
	if err := r.db.Select(&cities, query, args...); err != nil {
		return nil, err
	}
	return cities, nil
}

// Create adds a new city.
func (r *citiesRepository) Create(entity *models.City) (*models.City, error) {
	query := `
		INSERT INTO cities (name, population, image_url, timezone)
		VALUES (:name, :population, :image_url, :timezone)
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

// CreateMany adds multiple cities.
func (r *citiesRepository) CreateMany(entities []models.City) ([]models.City, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO cities (name, population, image_url, timezone)
		VALUES (:name, :population, :image_url, :timezone)
		RETURNING *`

	var results []models.City
	for _, entity := range entities {
		rows, err := tx.NamedQuery(query, entity)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var result models.City
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

// Update modifies a city by ID.
func (r *citiesRepository) Update(id int64, updates map[string]interface{}) (*models.City, error) {
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
		UPDATE cities
		SET %s
		WHERE id = $%d
		RETURNING *`, strings.Join(setClauses, ", "), argIndex)

	var city models.City
	if err := r.db.QueryRowx(query, args...).StructScan(&city); err != nil {
		return nil, err
	}
	return &city, nil
}

// UpdateMany modifies multiple city records based on the filter.
func (r *citiesRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
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
		UPDATE cities
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

// Delete removes a city by ID and returns the deleted row.
func (r *citiesRepository) Delete(id int64) (*models.City, error) {
	query := `
		DELETE FROM cities
		WHERE id = $1
		RETURNING *`

	var city models.City
	if err := r.db.QueryRowx(query, id).StructScan(&city); err != nil {
		return nil, err
	}
	return &city, nil
}

// DeleteMany removes multiple cities based on the filter and returns the deleted rows.
func (r *citiesRepository) DeleteMany(filter map[string]interface{}) ([]models.City, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf(`
		DELETE FROM cities
		WHERE %s
		RETURNING *`, strings.Join(whereClauses, " AND "))

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deletedCities []models.City
	for rows.Next() {
		var city models.City
		if err := rows.StructScan(&city); err != nil {
			return nil, err
		}
		deletedCities = append(deletedCities, city)
	}
	return deletedCities, nil
}
