package models

import "time"

type AuditLog struct {
	ID        int64          `json:"id"`
	TableName string         `json:"table_name"`
	Operation string         `json:"operation"`
	OldData   map[string]any `json:"old_data"`
	NewData   map[string]any `json:"new_data"`
	ChangedAt time.Time      `json:"changed_at"`
}
