package database

import (
	"testing"

	"gorm.io/driver/sqlite" // Используем SQLite для тестов
	"gorm.io/gorm"
)

// ConnectForTesting создает тестовое подключение к БД (SQLite в памяти)
func ConnectForTesting(t *testing.T) (*gorm.DB, error) {
	// Используем SQLite в памяти для тестов
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
		return nil, err
	}

	return db, nil
}

// SetupTestDB создает и мигрирует тестовую БД
func SetupTestDB(t *testing.T, models ...interface{}) *gorm.DB {
	db, err := ConnectForTesting(t)
	if err != nil {
		t.Fatal(err)
	}

	// Автомиграция для тестовых моделей
	if err := db.AutoMigrate(models...); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}