package repositories

import (
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
	var facility models.Facility
	query := `SELECT * FROM facilities WHERE id = $1`
	err := r.db.Get(&facility, query, id) // Use db.Get instead of db.Select
	if err != nil {
		return nil, err
	}
	return &facility, nil
}

// FindMany fetches facilities based on the filter.
func (r *facilityRepository) FindMany(filter map[string]interface{}) ([]models.Facility, error) {
	// Implement dynamic filtering logic if required
	var facilities []models.Facility
	query := "SELECT * FROM facilities"
	err := r.db.Select(&facilities, query)
	if err != nil {
		return nil, err
	}
	return facilities, nil
}

// Create adds a new facility.
func (r *facilityRepository) Create(entity *models.Facility) (*models.Facility, error) {
	query := `
		INSERT INTO facilities (name, type, category_id, city_id, location, coordinates, phone, emergency_phone, email, website, rating, bed_capacity, is_24_hours, has_emergency, has_parking, has_ambulance, accepts_insurance, description, image_url, amenities, accreditations, meta_data)
		VALUES (:name, :type, :category_id, :city_id, :location, :coordinates, :phone, :emergency_phone, :email, :website, :rating, :bed_capacity, :is_24_hours, :has_emergency, :has_parking, :has_ambulance, :accepts_insurance, :description, :image_url, :amenities, :accreditations, :meta_data)
		RETURNING *
	`
	rows, err := r.db.NamedQuery(query, entity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		rows.StructScan(entity)
	}
	return entity, nil
}

// CreateMany adds multiple facilities.
func (r *facilityRepository) CreateMany(entities []models.Facility) ([]models.Facility, error) {
	// Implement bulk insert if supported by the database
	return nil, nil // Placeholder
}

// Update modifies a facility by ID.
func (r *facilityRepository) Update(id int64, updates map[string]interface{}) (*models.Facility, error) {
	// Implement update logic
	return nil, nil // Placeholder
}

// UpdateMany modifies multiple facilities based on the filter.
func (r *facilityRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
	// Implement update logic
	return 0, nil // Placeholder
}

// Delete removes a facility by ID.
func (r *facilityRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM facilities WHERE id = $1", id)
	return err
}

// DeleteMany removes multiple facilities based on the filter.
func (r *facilityRepository) DeleteMany(filter map[string]interface{}) (int64, error) {
	// Implement delete logic
	return 0, nil // Placeholder
}
