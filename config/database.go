package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// ConnectDatabase establishes a connection to the PostgreSQL database
func ConnectDatabase() error {
	var err error

	DB, err = sql.Open("postgres", AppConfig.DatabaseURL)
	if err != nil {
		return err
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return err
	}

	// Configure connection pool
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	log.Println("Database connection established successfully")
	return nil
}

// CloseDatabase gracefully closes the database connection
func CloseDatabase() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}
