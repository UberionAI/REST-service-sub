package db

import (
	"REST-service-sub/internal/config"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DSN()
	gormLogger := logger.Default.LogMode(logger.Warn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed open gorm: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed get sqlDB: %w", err)
	}

	sqlDB.SetConnMaxLifetime(25)
	sqlDB.SetMaxIdleConns(10)

	log.Println("gorm connection successful...")
	return db, nil
}
