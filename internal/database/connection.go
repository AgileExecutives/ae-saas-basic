package database

import (
	"fmt"
	"log"

	"github.com/ae-saas-basic/ae-saas-basic/internal/models"
	"golang.org/x/crypto/bcrypt"
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

// CreateDatabaseIfNotExists creates the database if it doesn't exist
func CreateDatabaseIfNotExists(config Config) error {
	// Connect to PostgreSQL without specifying database (connect to 'postgres' database)
	adminDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.SSLMode)

	adminDB, err := gorm.Open(postgres.Open(adminDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent for admin operations
	})
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL server: %w", err)
	}

	// Get the underlying SQL database
	sqlDB, err := adminDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database: %w", err)
	}
	defer sqlDB.Close()

	// Check if database exists
	var exists bool
	query := `SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)`
	err = adminDB.Raw(query, config.DBName).Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		log.Printf("Creating database '%s'...", config.DBName)
		createQuery := fmt.Sprintf("CREATE DATABASE %s", config.DBName)
		err = adminDB.Exec(createQuery).Error
		if err != nil {
			return fmt.Errorf("failed to create database '%s': %w", config.DBName, err)
		}
		log.Printf("Database '%s' created successfully", config.DBName)
	} else {
		log.Printf("Database '%s' already exists", config.DBName)
	}

	return nil
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

// ConnectWithAutoCreate creates the database if it doesn't exist, then connects to it
func ConnectWithAutoCreate(config Config) (*gorm.DB, error) {
	// First, ensure the database exists
	if err := CreateDatabaseIfNotExists(config); err != nil {
		return nil, err
	}

	// Then connect to the database
	return Connect(config)
}

// Migrate runs database migrations with table existence checking
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// List of all models to migrate
	models := []interface{}{
		&models.Tenant{},
		&models.Plan{},
		&models.User{},
		&models.Customer{},
		&models.Contact{},
		&models.Email{},
		&models.Newsletter{},
		&models.UserSettings{},
		&models.TokenBlacklist{},
	}

	// Use GORM's AutoMigrate which handles dependencies automatically
	err := db.AutoMigrate(models...)
	if err != nil {
		log.Printf("ERROR: Migration failed: %v", err)
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Seed adds initial data to the database
func Seed(db *gorm.DB) error {
	log.Println("Seeding database with initial data...")

	// Create default tenant if it doesn't exist
	var tenantCount int64
	db.Model(&models.Tenant{}).Count(&tenantCount)
	if tenantCount == 0 {
		defaultTenant := models.Tenant{
			Name: "Default Tenant",
			Slug: "default-tenant",
		}
		if err := db.Create(&defaultTenant).Error; err != nil {
			return fmt.Errorf("failed to create default tenant: %w", err)
		}
		log.Println("Created default tenant")
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

	// Create test user if it doesn't exist
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		// Get the first tenant
		var defaultTenant models.Tenant
		if err := db.First(&defaultTenant).Error; err != nil {
			return fmt.Errorf("failed to find default tenant: %w", err)
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("newpass123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		testUser := models.User{
			Username:     "testuser",
			Email:        "testuser@example.com",
			PasswordHash: string(hashedPassword),
			FirstName:    "Test",
			LastName:     "User",
			TenantID:     defaultTenant.ID,
			Role:         "admin",
			Active:       true,
		}

		if err := db.Create(&testUser).Error; err != nil {
			return fmt.Errorf("failed to create test user: %w", err)
		}
		log.Println("Created test user")
	}

	log.Println("Database seeding completed successfully! 🎉")
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
