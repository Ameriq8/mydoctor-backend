package pg

import (
	"database/sql"
	"server/config"

	_ "github.com/lib/pq"
)

// Connect to the database
func Connect() (*sql.DB, error) {
	// get the environment variables
	host := config.LoadConfig().DBHost
	port := config.LoadConfig().DBPort
	user := config.LoadConfig().DBUser
	password := config.LoadConfig().DBPassword
	database := config.LoadConfig().DBName

	// build the connection string
	connStr := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " sslmode=disable"

	// connect to the database
	db, err := sql.Open(database, connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	return db, nil
}
