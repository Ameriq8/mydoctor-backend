package repositories

import (
	"fmt"
	"strings"
	"time"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// FacilityPlansRepository defines CRUD operations for the FacilityPlans model.
type FacilityPlansRepository interface {
	Repository[models.FacilityPlan]
}

// facilityPlansRepository is an implementation of FacilityPlansRepository.
type facilityPlansRepository struct {
	db *sqlx.DB
}

// NewFacilityPlansRepository initializes a new FacilityPlansRepository.
func NewFacilityPlansRepository(db *sqlx.DB) FacilityPlansRepository {
	return &facilityPlansRepository{db: db}
}

// Find fetches a facility plan by its ID.
func (r *facilityPlansRepository) Find(id int64) (*models.FacilityPlan, error) {
	start := time.Now()

	var plan models.FacilityPlan
	query := `SELECT * FROM facility_plans WHERE id = $1`
	err := r.db.Get(&plan, query, id)

	trackMetrics("Find", "facility_plans", start, err)

	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindMany fetches facility plans based on the filter.
func (r *facilityPlansRepository) FindMany(filter map[string]interface{}) ([]models.FacilityPlan, error) {
	start := time.Now()

	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := "SELECT * FROM facility_plans"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var plans []models.FacilityPlan
	err := r.db.Select(&plans, query, args...)

	trackMetrics("FindMany", "facility_plans", start, err)

	if err != nil {
		return nil, err
	}
	return plans, nil
}

// Create adds a new facility plan.
func (r *facilityPlansRepository) Create(entity *models.FacilityPlan) (*models.FacilityPlan, error) {
	start := time.Now()

	query := `
		INSERT INTO facility_plans (facility_id, plan_id, start_date, end_date, is_active)
		VALUES (:facility_id, :plan_id, :start_date, :end_date, :is_active)
		RETURNING *`

	rows, err := r.db.NamedQuery(query, entity)
	if err != nil {
		trackMetrics("Create", "facility_plans", start, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(entity); err != nil {
			trackMetrics("Create", "facility_plans", start, err)
			return nil, err
		}
	}
	trackMetrics("Create", "facility_plans", start, nil)
	return entity, nil
}

// CreateMany adds multiple facility plans.
func (r *facilityPlansRepository) CreateMany(entities []models.FacilityPlan) ([]models.FacilityPlan, error) {
	start := time.Now()

	tx, err := r.db.Beginx()
	if err != nil {
		trackMetrics("CreateMany", "facility_plans", start, err)
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO facility_plans (facility_id, plan_id, start_date, end_date, is_active)
		VALUES (:facility_id, :plan_id, :start_date, :end_date, :is_active)
		RETURNING *`

	var results []models.FacilityPlan
	for _, entity := range entities {
		rows, err := tx.NamedQuery(query, entity)
		if err != nil {
			trackMetrics("CreateMany", "facility_plans", start, err)
			return nil, err
		}
		defer rows.Close()

		var result models.FacilityPlan
		if rows.Next() {
			if err := rows.StructScan(&result); err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	err = tx.Commit()
	trackMetrics("CreateMany", "facility_plans", start, err)

	if err != nil {
		return nil, err
	}
	return results, nil
}

// Update modifies a facility plan by ID.
func (r *facilityPlansRepository) Update(id int64, updates map[string]interface{}) (*models.FacilityPlan, error) {
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
		UPDATE facility_plans
		SET %s
		WHERE id = $%d
		RETURNING *`, strings.Join(setClauses, ", "), argIndex)

	var plan models.FacilityPlan
	err := r.db.QueryRowx(query, args...).StructScan(&plan)

	trackMetrics("Update", "facility_plans", start, err)

	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// UpdateMany modifies multiple facility plans based on the filter.
func (r *facilityPlansRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
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
		UPDATE facility_plans
		SET %s
		WHERE %s`, strings.Join(setClauses, ", "), strings.Join(whereClauses, " AND "))

	result, err := r.db.Exec(query, args...)
	trackMetrics("UpdateMany", "facility_plans", start, err)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// Delete removes a facility plan by ID and returns the deleted row.
func (r *facilityPlansRepository) Delete(id int64) (*models.FacilityPlan, error) {
	start := time.Now()

	query := `
		DELETE FROM facility_plans
		WHERE id = $1
		RETURNING *`

	var plan models.FacilityPlan
	err := r.db.QueryRowx(query, id).StructScan(&plan)

	trackMetrics("Delete", "facility_plans", start, err)

	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// DeleteMany removes multiple facility plans based on the filter and returns the deleted rows.
func (r *facilityPlansRepository) DeleteMany(filter map[string]interface{}) ([]models.FacilityPlan, error) {
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
		DELETE FROM facility_plans
		WHERE %s
		RETURNING *`, strings.Join(whereClauses, " AND "))

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		trackMetrics("DeleteMany", "facility_plans", start, err)
		return nil, err
	}
	defer rows.Close()

	var deletedPlans []models.FacilityPlan
	for rows.Next() {
		var plan models.FacilityPlan
		if err := rows.StructScan(&plan); err != nil {
			return nil, err
		}
		deletedPlans = append(deletedPlans, plan)
	}
	trackMetrics("DeleteMany", "facility_plans", start, nil)
	return deletedPlans, nil
}
