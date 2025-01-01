package repositories

import (
	"fmt"
	"strings"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// ReviewsRepository defines CRUD operations for the Reviews model.
type ReviewsRepository interface {
	Repository[models.Review]
}

// reviewsRepository is an implementation of ReviewsRepository.
type reviewsRepository struct {
	db *sqlx.DB
}

// NewReviewsRepository initializes a new ReviewsRepository.
func NewReviewsRepository(db *sqlx.DB) ReviewsRepository {
	return &reviewsRepository{db: db}
}

// Find fetches a review by its ID.
func (r *reviewsRepository) Find(id int64) (*models.Review, error) {
	var review models.Review
	query := `SELECT * FROM reviews WHERE id = $1`
	if err := r.db.Get(&review, query, id); err != nil {
		return nil, err
	}
	return &review, nil
}

// FindMany fetches reviews based on the filter.
func (r *reviewsRepository) FindMany(filter map[string]interface{}) ([]models.Review, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := "SELECT * FROM reviews"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var reviews []models.Review
	if err := r.db.Select(&reviews, query, args...); err != nil {
		return nil, err
	}
	return reviews, nil
}

// Create adds a new review.
func (r *reviewsRepository) Create(entity *models.Review) (*models.Review, error) {
	query := `
		INSERT INTO reviews (entity_type, entity_id, user_id, rating, comment)
		VALUES (:entity_type, :entity_id, :user_id, :rating, :comment)
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

// CreateMany adds multiple reviews.
func (r *reviewsRepository) CreateMany(entities []models.Review) ([]models.Review, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO reviews (entity_type, entity_id, user_id, rating, comment)
		VALUES (:entity_type, :entity_id, :user_id, :rating, :comment)
		RETURNING *`

	var results []models.Review
	for _, entity := range entities {
		rows, err := tx.NamedQuery(query, entity)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var result models.Review
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

// Update modifies a review by ID.
func (r *reviewsRepository) Update(id int64, updates map[string]interface{}) (*models.Review, error) {
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
		UPDATE reviews
		SET %s
		WHERE id = $%d
		RETURNING *`, strings.Join(setClauses, ", "), argIndex)

	var review models.Review
	if err := r.db.QueryRowx(query, args...).StructScan(&review); err != nil {
		return nil, err
	}
	return &review, nil
}

// UpdateMany modifies multiple reviews based on the filter.
func (r *reviewsRepository) UpdateMany(filter map[string]interface{}, updates map[string]interface{}) (int64, error) {
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
		UPDATE reviews
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

// Delete removes a review by ID and returns the deleted row.
func (r *reviewsRepository) Delete(id int64) (*models.Review, error) {
	query := `
		DELETE FROM reviews
		WHERE id = $1
		RETURNING *`

	var review models.Review
	if err := r.db.QueryRowx(query, id).StructScan(&review); err != nil {
		return nil, err
	}
	return &review, nil
}

// DeleteMany removes multiple reviews based on the filter and returns the deleted rows.
func (r *reviewsRepository) DeleteMany(filter map[string]interface{}) ([]models.Review, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf(`
		DELETE FROM reviews
		WHERE %s
		RETURNING *`, strings.Join(whereClauses, " AND "))

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deletedReviews []models.Review
	for rows.Next() {
		var review models.Review
		if err := rows.StructScan(&review); err != nil {
			return nil, err
		}
		deletedReviews = append(deletedReviews, review)
	}
	return deletedReviews, nil
}
