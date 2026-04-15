package database

import (
	"log"

	"job-portal/config"
	"job-portal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect initializes the database connection and runs auto-migration
func Connect(cfg *config.Config) {
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if cfg.AppEnv == "production" {
		gormCfg.Logger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), gormCfg)
	if err != nil {
		log.Fatalf("[FATAL] Failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[FATAL] Failed to get underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	DB = db
	log.Println("[INFO] Database connected successfully")

	migrate(db)
}

// migrate runs GORM auto-migration for all models
func migrate(db *gorm.DB) {
	log.Println("[INFO] Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Job{},
		&models.Application{},
	)
	if err != nil {
		log.Fatalf("[FATAL] Migration failed: %v", err)
	}

	log.Println("[INFO] Database migration completed")
}
