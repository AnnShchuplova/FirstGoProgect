package services

import (
	"FurryTrack/internal/models"
	"FurryTrack/internal/repositories"
	
	"github.com/google/uuid"
)

type AdminService struct {
	repo *repositories.AdminRepository
}

func NewAdminService(repo *repositories.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *AdminService) BanUser(adminID, userID uuid.UUID, reason string) error {
	if err := s.repo.BanUser(userID); err != nil {
		return err
	}
	
	return s.repo.LogAction(&models.AdminAction{
		AdminID:     adminID,
		UserID:      userID,
		ActionType:  "ban",
		Description: reason,
	})
}

func (s *AdminService) UnbanUser(adminID, userID uuid.UUID, reason string) error {
	if err := s.repo.UnbanUser(userID); err != nil {
		return err
	}
	
	return s.repo.LogAction(&models.AdminAction{
		AdminID:     adminID,
		UserID:      userID,
		ActionType:  "unban",
		Description: reason,
	})
}

func (s *AdminService) DeleteUser(adminID, userID uuid.UUID, reason string) error {
	if err := s.repo.DeleteUser(userID); err != nil {
		return err
	}
	
	return s.repo.LogAction(&models.AdminAction{
		AdminID:     adminID,
		UserID:      userID,
		ActionType:  "delete",
		Description: reason,
	})
}

func (s *AdminService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}
