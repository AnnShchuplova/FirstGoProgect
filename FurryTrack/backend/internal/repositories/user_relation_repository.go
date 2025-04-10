package repositories

import (
    "FurryTrack/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type UserRelationRepository struct {
    db *gorm.DB
}

func NewUserRelationRepository(db *gorm.DB) *UserRelationRepository {
    return &UserRelationRepository{db: db}
}

// Создание подписки
func (r *UserRelationRepository) Follow(followerID, followingID uuid.UUID) error {
    relation := models.UserRelation{
        FollowerID:  followerID,
        FollowingID: followingID,
    }
    result := r.db.Create(&relation)
    return result.Error
}

// GetFollowing - получения списка подписок
func (r *UserRelationRepository) GetFollowing(userID uuid.UUID) ([]models.UserRelation, error) {
    var relations []models.UserRelation
    err := r.db.Where("follower_id = ?", userID).Find(&relations).Error
    return relations, err
}

// GetFollowers - получения списка подписчиков
func (r *UserRelationRepository) GetFollowers(userID uuid.UUID) ([]models.UserRelation, error) {
    var relations []models.UserRelation
    err := r.db.Where("following_id = ?", userID).Find(&relations).Error
    return relations, err
}