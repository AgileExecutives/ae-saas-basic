package database

import (
	"fmt"
	"log"

	"github.com/ae-saas-basic/ae-saas-basic/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Connect creates a database connection
func Connect(config Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Auto-migrate all models
	err := db.AutoMigrate(
		&models.Organization{},
		&models.Plan{},
		&models.User{},
		&models.Customer{},
		&models.Contact{},
		&models.Email{},
		&models.UserSettings{},
		&models.TokenBlacklist{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Seed adds initial data to the database
func Seed(db *gorm.DB) error {
	log.Println("Seeding database with initial data...")

	// Create default organization if it doesn't exist
	var orgCount int64
	db.Model(&models.Organization{}).Count(&orgCount)
	if orgCount == 0 {
		defaultOrg := models.Organization{
			Name: "Default Organization",
			Slug: "default-org",
		}
		if err := db.Create(&defaultOrg).Error; err != nil {
			return fmt.Errorf("failed to create default organization: %w", err)
		}
		log.Println("Created default organization")
	}

	// Create default plan if it doesn't exist
	var planCount int64
	db.Model(&models.Plan{}).Count(&planCount)
	if planCount == 0 {
		defaultPlan := models.Plan{
			Name:          "Basic Plan",
			Slug:          "basic",
			Description:   "Basic SaaS plan with essential features",
			Price:         29.99,
			Currency:      "EUR",
			InvoicePeriod: "monthly",
			MaxUsers:      10,
			MaxClients:    100,
			Features:      `{"users": 10, "clients": 100, "support": "email"}`,
			Active:        true,
		}
		if err := db.Create(&defaultPlan).Error; err != nil {
			return fmt.Errorf("failed to create default plan: %w", err)
		}
		log.Println("Created default plan")
	}

	log.Println("Database seeding completed")
	return nil
}

// GetDefaultConfig returns default database configuration
func GetDefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "password",
		DBName:   "ae_saas_basic",
		SSLMode:  "disable",
	}
}
