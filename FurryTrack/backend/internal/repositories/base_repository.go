package repositories

import (
	"gorm.io/gorm"
	// "FurryTrack/internal/models"
)

type BaseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}
