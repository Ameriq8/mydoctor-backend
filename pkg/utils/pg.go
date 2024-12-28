package pg

import (
	"fmt"
	"server/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB initializes and returns a new database connection using sqlx
func NewDB() (*sqlx.DB, error) {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Connect(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(0)

	return db, nil
}
