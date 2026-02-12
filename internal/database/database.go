package database

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Info("connected to database")

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
	log.Printf("Migrating database: %s", dbname)
	sqlDB, err := s.db.DB()
	_, err = sqlDB.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}
	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id INT AUTO_INCREMENT PRIMARY KEY,
			original_url VARCHAR(255) NOT NULL,
		    short_url VARCHAR(8) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP NULL DEFAULT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	log.Printf("Database migrations completed successfully.")
	return nil
}
