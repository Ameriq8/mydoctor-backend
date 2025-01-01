package repositories

import (
	"fmt"
	"strings"
	"time"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// PlansRepository defines CRUD operations for the Plans model.
type PlansRepository interface {
	Repository[models.Plan]
}

// plansRepository is an implementation of PlansRepository.
type plansRepository struct {
	db *sqlx.DB
}

// NewPlansRepository initializes a new PlansRepository.
func NewPlansRepository(db *sqlx.DB) PlansRepository {
	return &plansRepository{db: db}
}

// Find fetches a plan by its ID.
func (r *plansRepository) Find(id int64) (*models.Plan, error) {
	start := time.Now()

	var plan models.Plan
	query := `SELECT * FROM plans WHERE id = $1`
	err := r.db.Get(&plan, query, id)

	trackMetrics("Find", "plans", start, err)

	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindMany fetches plans based on the filter.
func (r *plansRepository) FindMany(filter map[string]interface{}) ([]models.Plan, error) {
	start := time.Now()

	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := "SELECT * FROM plans"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var plans []models.Plan
	err := r.db.Select(&plans, query, args...)

	trackMetrics("FindMany", "plans", start, err)

	if err != nil {
		return nil, err
	}
	return plans, nil
}

// Create adds a new plan.
func (r *plansRepository) Create(entity *models.Plan) (*models.Plan, error) {
	start := time.Now()

	query := `
		INSERT INTO plans (name, monthly_price, yearly_price, description, features)
		VALUES (:name, :monthly_price, :yearly_price, :description, :features)
		RETURNING *`

	rows, err := r.db.NamedQuery(query, entity)
	if err != nil {
		trackMetrics("Create", "plans", start, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(entity); err != nil {
			trackMetrics("Create", "plans", start, err)
			return nil, err
		}
	}
	trackMetrics("Create", "plans", start, nil)
	return entity, nil
}

// CreateMany adds multiple plans.
func (r *plansRepository) CreateMany(entities []models.Plan) ([]models.Plan, error) {
	start := time.Now()

	tx, err := r.db.Beginx()
	if err != nil {
		trackMetrics("CreateMany", "plans", start, err)
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO plans (name, monthly_price, yearly_price, description, features)
		VALUES (:name, :monthly_price, :yearly_price, :description, :features)
		RETURNING *`

	var results []models.Plan
	for _, entity := range entities {
		rows, err := tx.NamedQuery(query, entity)
		if err != nil {
			trackMetrics("CreateMany", "plans", start, err)
			return nil, err
		}
		defer rows.Close()

		var result models.Plan
		if rows.Next() {
			if err := rows.StructScan(&result); err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	if err := tx.Commit(); err != nil {
		trackMetrics("CreateMany", "plans", start, err)
		return nil, err
	}
	trackMetrics("CreateMany", "plans", start, nil)
	return results, nil
}

// Update modifies a plan by ID.
func (r *plansRepository) Update(id int64, updates map[string]interface{}) (*models.Plan, error) {
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
		UPDATE plans
		SET %s
		WHERE id = $%d
		RETURNING *`, strings.Join(setClauses, ", "), argIndex)

	var plan models.Plan
	err := r.db.QueryRowx(query, args...).StructScan(&plan)

	trackMetrics("Update", "plans", start, err)

	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// UpdateMany modifies multiple plans based on the filter.
func (r *plansRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
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
		UPDATE plans
		SET %s
		WHERE %s`, strings.Join(setClauses, ", "), strings.Join(whereClauses, " AND "))

	result, err := r.db.Exec(query, args...)

	trackMetrics("UpdateMany", "plans", start, err)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// Delete removes a plan by ID and returns the deleted row.
func (r *plansRepository) Delete(id int64) (*models.Plan, error) {
	start := time.Now()

	query := `
		DELETE FROM plans
		WHERE id = $1
		RETURNING *`

	var plan models.Plan
	err := r.db.QueryRowx(query, id).StructScan(&plan)

	trackMetrics("Delete", "plans", start, err)

	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// DeleteMany removes multiple plans based on the filter and returns the deleted rows.
func (r *plansRepository) DeleteMany(filter map[string]interface{}) ([]models.Plan, error) {
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
		DELETE FROM plans
		WHERE %s
		RETURNING *`, strings.Join(whereClauses, " AND "))

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		trackMetrics("DeleteMany", "plans", start, err)
		return nil, err
	}
	defer rows.Close()

	var deletedPlans []models.Plan
	for rows.Next() {
		var plan models.Plan
		if err := rows.StructScan(&plan); err != nil {
			return nil, err
		}
		deletedPlans = append(deletedPlans, plan)
	}
	trackMetrics("DeleteMany", "plans", start, nil)
	return deletedPlans, nil
}
