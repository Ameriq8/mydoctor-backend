package repositories

import (
	"fmt"
	"strings"
	"time"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// FacilityRepository defines CRUD operations for the Facility model.
type FacilityRepository interface {
	Repository[models.Facility]
}

// facilityRepository is an implementation of FacilityRepository.
type facilityRepository struct {
	db *sqlx.DB
}

// NewFacilityRepository initializes a new FacilityRepository.
func NewFacilityRepository(db *sqlx.DB) FacilityRepository {
	return &facilityRepository{db: db}
}

// Find fetches a facility by its ID.
func (r *facilityRepository) Find(id int64) (*models.Facility, error) {
	start := time.Now()

	var facility models.Facility
	query := `SELECT * FROM facilities WHERE id = $1`
	err := r.db.Get(&facility, query, id)

	trackMetrics("Find", "facilities", start, err)

	if err != nil {
		return nil, err
	}
	return &facility, nil
}

// FindMany fetches facilities based on the filter.
func (r *facilityRepository) FindMany(filter map[string]interface{}) ([]models.Facility, error) {
	start := time.Now()

	var facilities []models.Facility
	query := "SELECT * FROM facilities"
	err := r.db.Select(&facilities, query)

	trackMetrics("FindMany", "facilities", start, err)

	if err != nil {
		return nil, err
	}
	return facilities, nil
}

// Create adds a new facility.
func (r *facilityRepository) Create(entity *models.Facility) (*models.Facility, error) {
	start := time.Now()

	query := `
		INSERT INTO facilities (name, type, category_id, city_id, location, coordinates, phone, emergency_phone, email, website, rating, bed_capacity, is_24_hours, has_emergency, has_parking, has_ambulance, accepts_insurance, description, image_url, amenities, accreditations, meta_data)
		VALUES (:name, :type, :category_id, :city_id, :location, :coordinates, :phone, :emergency_phone, :email, :website, :rating, :bed_capacity, :is_24_hours, :has_emergency, :has_parking, :has_ambulance, :accepts_insurance, :description, :image_url, :amenities, :accreditations, :meta_data)
		RETURNING *
	`
	rows, err := r.db.NamedQuery(query, entity)
	if err != nil {
		trackMetrics("Create", "facilities", start, err)
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		rows.StructScan(entity)
	}
	trackMetrics("Create", "facilities", start, nil)
	return entity, nil
}

// CreateMany adds multiple facilities.
func (r *facilityRepository) CreateMany(entities []models.Facility) ([]models.Facility, error) {
	start := time.Now()

	query := `
    INSERT INTO facilities (name, type, category_id, city_id, location, coordinates, phone, emergency_phone, email, website, rating, bed_capacity, is_24_hours, has_emergency, has_parking, has_ambulance, accepts_insurance, description, image_url, amenities, accreditations, meta_data)
    VALUES (:name, :type, :category_id, :city_id, :location, :coordinates, :phone, :emergency_phone, :email, :website, :rating, :bed_capacity, :is_24_hours, :has_emergency, :has_parking, :has_ambulance, :accepts_insurance, :description, :image_url, :amenities, :accreditations, :meta_data)
    RETURNING *;`
	tx, err := r.db.Beginx()
	if err != nil {
		trackMetrics("CreateMany", "facilities", start, err)
		return nil, err
	}
	defer tx.Rollback()

	var results []models.Facility
	for _, entity := range entities {
		rows, err := r.db.NamedQuery(query, entity)
		if err != nil {
			trackMetrics("CreateMany", "facilities", start, err)
			return nil, err
		}
		defer rows.Close()

		var result models.Facility
		if rows.Next() {
			if err := rows.StructScan(&result); err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	if err := tx.Commit(); err != nil {
		trackMetrics("CreateMany", "facilities", start, err)
		return nil, err
	}

	trackMetrics("CreateMany", "facilities", start, nil)
	return results, nil
}

// Update modifies a facility by ID.
func (r *facilityRepository) Update(id int64, updates map[string]interface{}) (*models.Facility, error) {
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

	query := fmt.Sprintf(`UPDATE facilities SET %s WHERE id = $%d RETURNING *`, strings.Join(setClauses, ", "), argIndex)

	var facility models.Facility
	err := r.db.QueryRowx(query, args...).StructScan(&facility)
	trackMetrics("Update", "facilities", start, err)

	if err != nil {
		return nil, err
	}

	return &facility, nil
}

// UpdateMany modifies multiple facilities based on the filter.
func (r *facilityRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
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
		UPDATE facilities
		SET %s
		WHERE %s
	`, strings.Join(setClauses, ", "), strings.Join(whereClauses, " AND "))

	result, err := r.db.Exec(query, args...)
	trackMetrics("UpdateMany", "facilities", start, err)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// Delete removes a facility by ID and returns the deleted row.
func (r *facilityRepository) Delete(id int64) (*models.Facility, error) {
	start := time.Now()

	query := `
        DELETE FROM facilities
        WHERE id = $1
        RETURNING id, name, type, category_id, city_id, location, coordinates, phone, emergency_phone, email, website, rating, bed_capacity, is_24_hours, has_emergency, has_parking, has_ambulance, accepts_insurance, description, image_url, amenities, accreditations, meta_data
    `

	var facility models.Facility
	err := r.db.QueryRowx(query, id).StructScan(&facility)
	trackMetrics("Delete", "facilities", start, err)

	if err != nil {
		return nil, err
	}

	return &facility, nil
}

// DeleteMany removes multiple facilities based on the filter and returns the deleted rows.
func (r *facilityRepository) DeleteMany(filter map[string]interface{}) ([]models.Facility, error) {
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
        DELETE FROM facilities
        WHERE %s
        RETURNING id, name, type, category_id, city_id, location, coordinates, phone, emergency_phone, email, website, rating, bed_capacity, is_24_hours, has_emergency, has_parking, has_ambulance, accepts_insurance, description, image_url, amenities, accreditations, meta_data
    `, strings.Join(whereClauses, " AND "))

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		trackMetrics("DeleteMany", "facilities", start, err)
		return nil, err
	}
	defer rows.Close()

	var deletedFacilities []models.Facility
	for rows.Next() {
		var facility models.Facility
		if err := rows.StructScan(&facility); err != nil {
			return nil, err
		}
		deletedFacilities = append(deletedFacilities, facility)
	}

	trackMetrics("DeleteMany", "facilities", start, nil)
	return deletedFacilities, nil
}
