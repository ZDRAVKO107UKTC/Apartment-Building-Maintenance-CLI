package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ZDRAVKO107UKTC/Apartment-Building-Maintenance-CLI/internal/model"
)

// DB is the process-wide GORM handle used by the repository layer.
var DB *gorm.DB

// Init connects to the configured SQLite database and keeps the schema in sync
// via GORM AutoMigrate. The database location is injected via the DB_PATH
// environment variable; nothing is hardcoded.
func Init() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH environment variable is not set")
	}

	gormLogger := logger.New(
		log.New(os.Stderr, "", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
		},
	)

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{Logger: gormLogger})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := DB.AutoMigrate(&model.Issue{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	fmt.Println("Database connected and migrated successfully")
}
