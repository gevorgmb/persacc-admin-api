package data

import (
	"fmt"
	"log"
	"os"

	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// Fallback for local development or raise error
		dsn = "host=localhost user=postgres password=postgres dbname=persacc port=5452 sslmode=disable"
		log.Println("DB_DSN not set, using default...")
	}

	// Clean DSN: remove leading/trailing whitespace and quotes
	dsn = strings.TrimSpace(dsn)
	dsn = strings.Trim(dsn, "\"'")

	log.Printf("DB_DSN (quoted): %q", dsn)
	log.Printf("DB_DSN (bytes): %v", []byte(dsn))

	log.Printf("Attempting DSN: %s", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
