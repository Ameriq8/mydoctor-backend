package repositories

import (
	"fmt"
	"strings"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// FacilityEquipmentRepository defines CRUD operations for the FacilityEquipment model.
type FacilityEquipmentRepository interface {
	Repository[models.FacilityEquipment]
}

// facilityEquipmentRepository is an implementation of FacilityEquipmentRepository.
type facilityEquipmentRepository struct {
	db *sqlx.DB
}

// NewFacilityEquipmentRepository initializes a new FacilityEquipmentRepository.
func NewFacilityEquipmentRepository(db *sqlx.DB) FacilityEquipmentRepository {
	return &facilityEquipmentRepository{db: db}
}

// Find fetches a facility equipment by its ID.
func (r *facilityEquipmentRepository) Find(id int64) (*models.FacilityEquipment, error) {
	var equipment models.FacilityEquipment
	query := `SELECT * FROM facility_equipment WHERE id = $1`
	if err := r.db.Get(&equipment, query, id); err != nil {
		return nil, err
	}
	return &equipment, nil
}

// FindMany fetches facility equipment based on the filter.
func (r *facilityEquipmentRepository) FindMany(filter map[string]interface{}) ([]models.FacilityEquipment, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := "SELECT * FROM facility_equipment"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var equipment []models.FacilityEquipment
	if err := r.db.Select(&equipment, query, args...); err != nil {
		return nil, err
	}
	return equipment, nil
}

// Create adds a new facility equipment.
func (r *facilityEquipmentRepository) Create(entity *models.FacilityEquipment) (*models.FacilityEquipment, error) {
	query := `
		INSERT INTO facility_equipment (facility_id, department_id, name, model, manufacturer, purchase_date, last_maintenance_date, next_maintenance_date, status)
		VALUES (:facility_id, :department_id, :name, :model, :manufacturer, :purchase_date, :last_maintenance_date, :next_maintenance_date, :status)
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

// CreateMany adds multiple facility equipment records.
func (r *facilityEquipmentRepository) CreateMany(entities []models.FacilityEquipment) ([]models.FacilityEquipment, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO facility_equipment (facility_id, department_id, name, model, manufacturer, purchase_date, last_maintenance_date, next_maintenance_date, status)
		VALUES (:facility_id, :department_id, :name, :model, :manufacturer, :purchase_date, :last_maintenance_date, :next_maintenance_date, :status)
		RETURNING *`

	var results []models.FacilityEquipment
	for _, entity := range entities {
		rows, err := tx.NamedQuery(query, entity)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var result models.FacilityEquipment
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

// Update modifies a facility equipment record by ID.
func (r *facilityEquipmentRepository) Update(id int64, updates map[string]interface{}) (*models.FacilityEquipment, error) {
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
		UPDATE facility_equipment
		SET %s
		WHERE id = $%d
		RETURNING *`, strings.Join(setClauses, ", "), argIndex)

	var equipment models.FacilityEquipment
	if err := r.db.QueryRowx(query, args...).StructScan(&equipment); err != nil {
		return nil, err
	}
	return &equipment, nil
}

// UpdateMany modifies multiple facility equipment records based on the filter.
func (r *facilityEquipmentRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
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
		UPDATE facility_equipment
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

// Delete removes a facility equipment record by ID and returns the deleted row.
func (r *facilityEquipmentRepository) Delete(id int64) (*models.FacilityEquipment, error) {
	query := `
		DELETE FROM facility_equipment
		WHERE id = $1
		RETURNING *`

	var equipment models.FacilityEquipment
	if err := r.db.QueryRowx(query, id).StructScan(&equipment); err != nil {
		return nil, err
	}
	return &equipment, nil
}

// DeleteMany removes multiple facility equipment records based on the filter and returns the deleted rows.
func (r *facilityEquipmentRepository) DeleteMany(filter map[string]interface{}) ([]models.FacilityEquipment, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf(`
		DELETE FROM facility_equipment
		WHERE %s
		RETURNING *`, strings.Join(whereClauses, " AND "))

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deletedEquipment []models.FacilityEquipment
	for rows.Next() {
		var equipment models.FacilityEquipment
		if err := rows.StructScan(&equipment); err != nil {
			return nil, err
		}
		deletedEquipment = append(deletedEquipment, equipment)
	}
	return deletedEquipment, nil
}
