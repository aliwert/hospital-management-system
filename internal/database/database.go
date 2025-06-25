package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aliwert/go-hospital-management/internal/models"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var SqlDB *sql.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	// Initialize GORM DB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize raw SQL DB
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database with raw SQL:", err)
	}

	// Test raw SQL connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database with raw SQL:", err)
	}

	// Set connection pool
	gormSqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	gormSqlDB.SetMaxIdleConns(10)
	gormSqlDB.SetMaxOpenConns(100)
	gormSqlDB.SetConnMaxLifetime(time.Hour)

	// Set connection pool for raw SQL
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.AutoMigrate(
		&models.User{},
		&models.Doctor{},
		&models.Supplier{},
		&models.Inventory{},
		&models.Department{},
		&models.MedicalRecord{},
		&models.Appointment{},
		&models.DoctorSchedule{},
		&models.Patient{},
		&models.Prescription{},
		&models.TestResult{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB = db
	SqlDB = sqlDB
	log.Println("Database connection established successfully")
}

func GetDB() *gorm.DB {
	return DB
}

func GetSqlDB() *sql.DB {
	return SqlDB
}

// HealthCheck performs a health check on the database
func HealthCheck() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
