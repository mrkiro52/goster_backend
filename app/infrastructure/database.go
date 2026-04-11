package gorm

import (
	"goster/config"
	"goster/domain"
	logging "goster/library/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDb() *gorm.DB {
	// Функция для создания подключения к бд
	logger := logging.GetLogger("ConnectToDb")

	db, err := gorm.Open(postgres.Open(config.DBuri), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Не подключиться к базе: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("Не получить sql.DB из gorm: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("База подключена")

	return db
}

func AutoMigrate(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"")

	return db.AutoMigrate(
		&domain.User{},
	)
}
