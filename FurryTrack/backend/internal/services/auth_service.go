package services

import (
	"FurryTrack/internal/models"
	"FurryTrack/internal/repositories"
	"FurryTrack/pkg/utils"
	"errors"
	"fmt"
	//"strings"
	//"log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Регистрация нового пользователя 
func (s *AuthService) Register(username, email, password string, isAdmin bool) (*models.User, error) {
    if exists, _ := s.userRepo.ExistsByEmail(email); exists {
        return nil, fmt.Errorf("user already exists")
    }
    // Хеширование пароля
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("password hashing failed")
    }
    user := &models.User{
        Username:     username,
        Email:        email,
        PasswordHash: string(hashedPassword),
        IsAdmin:      isAdmin, 
    }
    if err := s.userRepo.CreateUser(user); err != nil {
        return nil, fmt.Errorf("failed to create user")
    }
    return user, nil
}

// Login - выполняет вход пользователя (уже зарегистрированного)
func (s *AuthService) Login(email, password string) (string, *models.User, error) {
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return "", nil, errors.New("invalid credentials")
    }
    // Проверяем пароль с хэшем
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return "", nil, errors.New("invalid credentials")
    }
    token, err := utils.GenerateToken(user.ID, user.Role, s.jwtSecret)
    if err != nil {
        return "", nil, fmt.Errorf("failed to generate token: %w", err)
    }

    return token, user, nil
}

// GetUserProfile - получение профиля пользователя
func (s *AuthService) GetUserProfile(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return user, nil
}


func (s *AuthService) SetUserRole(userID uuid.UUID, role string) error {
    // Проверка допустимых ролей
    validRoles := map[string]bool{
        "USER":  true,
        "ADMIN": true,
        "VET":   true,
    }
    if !validRoles[role] {
        return fmt.Errorf("invalid role")
    }

    return s.userRepo.UpdateUser(userID, map[string]interface{}{
        "role":    role,
        "is_admin": role == "ADMIN", 
    })
}

func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
    return s.userRepo.FindByEmail(email)
}

func (s *AuthService) IsAdmin(userID uuid.UUID) bool {
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return false
    }
    return user.IsAdmin
}