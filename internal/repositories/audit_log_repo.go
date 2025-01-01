package repositories

import (
	"fmt"
	"strings"

	"server/internal/models"

	"github.com/jmoiron/sqlx"
)

// AuditLogRepository defines read-only operations for the AuditLog model.
type AuditLogRepository interface {
	Find(id int64) (*models.AuditLog, error)
	FindMany(filter map[string]interface{}) ([]models.AuditLog, error)
}

// auditLogRepository is an implementation of AuditLogRepository.
type auditLogRepository struct {
	db *sqlx.DB
}

// NewAuditLogRepository initializes a new AuditLogRepository.
func NewAuditLogRepository(db *sqlx.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

// Find fetches an audit log entry by its ID.
func (r *auditLogRepository) Find(id int64) (*models.AuditLog, error) {
	var log models.AuditLog
	query := `SELECT * FROM audit_log WHERE id = $1`
	if err := r.db.Get(&log, query, id); err != nil {
		return nil, err
	}
	return &log, nil
}

// FindMany fetches audit log entries based on the filter.
func (r *auditLogRepository) FindMany(filter map[string]interface{}) ([]models.AuditLog, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for key, value := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := "SELECT * FROM audit_log"
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var logs []models.AuditLog
	if err := r.db.Select(&logs, query, args...); err != nil {
		return nil, err
	}
	return logs, nil
}
