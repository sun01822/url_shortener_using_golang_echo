package database

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

var (
	dbname   = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=UTC",
		host,
		username,
		password,
		dbname,
		port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Printf("connected to PostgreSQL database")

	return db
}

func GetService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *Service) Health() map[string]string {
	sqlDB, err := s.db.DB()
	if err != nil {
		log.Error(err)
		return map[string]string{
			"status": "unhealthy",
			"error":  err.Error(),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = sqlDB.PingContext(ctx)
	if err != nil {
		log.Error(err)
		return map[string]string{
			"status": "unhealthy",
			"error":  err.Error(),
		}
	}

	return map[string]string{
		"status": "healthy",
	}
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *Service) Close() error {
	log.Printf("Disconnected from database: %s", dbname)
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %v", err)
	}
	return sqlDB.Close()
}

func (s *Service) Migrate() error {
	log.Printf("Running PostgreSQL migrations")

	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %v", err)
	}

	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id BIGSERIAL PRIMARY KEY,
			original_url VARCHAR(255) NOT NULL,
			short_url VARCHAR(8) NOT NULL UNIQUE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Printf("PostgreSQL migrations completed successfully.")
	return nil
}

func (s *Service) Clear() error {
	log.Printf("Clearing all urls")

	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %v", err)
	}

	_, err = sqlDB.Exec(`DELETE From URLS WHERE 1 = 1`)
	if err != nil {
		return fmt.Errorf("failed to delete")
	}

	log.Printf("Clearing all short_urls")
	return nil
}
