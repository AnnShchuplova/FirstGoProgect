package services

import (
	"FurryTrack/internal/models"
	"FurryTrack/internal/repositories"
	"github.com/google/uuid"
)

type UserRelationService struct {
	repo *repositories.UserRelationRepository
}

// Обращения к репозиторию
func NewUserRelationService(repo *repositories.UserRelationRepository) *UserRelationService {
	return &UserRelationService{repo: repo}
}

func (s *UserRelationService) Follow(followerID, followingID uuid.UUID) error {
	return s.repo.Follow(followerID, followingID)
}

func (s *UserRelationService) GetFollowing(userID uuid.UUID) ([]models.UserRelation, error) {
	return s.repo.GetFollowing(userID)
}

func (s *UserRelationService) GetFollowers(userID uuid.UUID) ([]models.UserRelation, error) {
	return s.repo.GetFollowers(userID)
}