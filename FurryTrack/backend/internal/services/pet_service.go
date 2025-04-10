package services

import (
	"FurryTrack/internal/models"
	"FurryTrack/internal/repositories"
	"errors"
	"fmt"
	"time"
	"gorm.io/gorm"

	"github.com/google/uuid"
	// "FurryTrack/pkg/database"
	"mime/multipart"
	"path/filepath"
	"os"
	"io"

)

var (
    ErrPetNotFound      = errors.New("pet not found")
    ErrPermissionDenied = errors.New("permission denied")
)

type PetService struct {
	petRepo           *repositories.PetRepository
	vaccineRepo       *repositories.VaccineRepository
	vaccineRecordRepo *repositories.VaccineRecordRepository
}

func NewPetService(petRepo *repositories.PetRepository) *PetService {
	return &PetService{
		petRepo: petRepo,
	}
}

// CreatePet создает нового питомца
func (s *PetService) CreatePet(pet *models.Pet) error {
	// Валидация на слишком короткое имя
	if len(pet.Name) < 2 {
		return errors.New("pet name too short")
	}

	return s.petRepo.Create(pet)
}

// GetUserPets возвращает питомцев пользователя
func (s *PetService) GetUserPets(userID uuid.UUID) ([]models.Pet, error) {
	return s.petRepo.FindByOwner(userID)
}

// GetPetVaccines — получение всех вакцин питомца
func (s *PetService) GetPetVaccines(petID uuid.UUID) ([]models.VaccineRecord, error) {
	return s.vaccineRecordRepo.GetRecordsByPet(petID)
}

func (s *PetService) GetPetByID(petID uuid.UUID, ownerID uuid.UUID) (*models.Pet, error) {
	pet, err := s.petRepo.FindByID(petID)
	if err != nil {
		return nil, err
	}
	// Проверка, что запрашивающий пользователь - владелец
	if pet.OwnerID != ownerID {
		return nil, fmt.Errorf("access denied: you are not the owner")
	}

	return pet, nil
}

// Добавляем запись о прививке
func (s *PetService) AddVaccine(
	vaccineID uuid.UUID,
	petID uuid.UUID,
	userID uuid.UUID,
	vaccineName string,
	date time.Time,
	clinic string,
	NextDate time.Time,
) (*models.VaccineRecord, error) {
	pet, err := s.petRepo.FindByID(petID)
	if err != nil {
		return nil, fmt.Errorf("pet not found: %w", err)
	}

	if pet.OwnerID != userID {
		return nil, fmt.Errorf("access denied: you are not the owner")
	}
	vaccine, err := s.vaccineRepo.FindByName(vaccineName)
	if err != nil {
		return nil, fmt.Errorf("vaccine not found: %w", err)
	}

	nextDate := date.AddDate(0, 0, vaccine.DurationDays)

	record := &models.VaccineRecord{
		PetID:     petID,
		VaccineID: vaccineID,
		Date:      date,
		Clinic:    clinic,
		NextDate:  nextDate,
	}

	if err := s.vaccineRecordRepo.AddRecord(record); err != nil {
		return nil, fmt.Errorf("failed to create vaccine record: %w", err)
	}

	return record, nil
}

// Обновление питомца
func (s *PetService) UpdatePet(
    petID uuid.UUID,
    userID uuid.UUID,
    name *string,
    breed *string,
    birthDate *time.Time, 
) (*models.Pet, error) {
    pet, err := s.petRepo.FindByID(petID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrPetNotFound
        }
        return nil, err
    }

    // Проверяем права доступа
    if pet.OwnerID != userID {
        return nil, ErrPermissionDenied
    }

    if name != nil {
        pet.Name = *name
    }
    if breed != nil {
        pet.Breed = *breed
    }
    if birthDate != nil {
        pet.BirthDate = *birthDate 
    }

    if err := s.petRepo.Update(pet); err != nil {
        return nil, err
    }

    return pet, nil
}

// Удаление питомца
func (s *PetService) DeletePet(petID, userID uuid.UUID) error {
	pet, err := s.petRepo.FindByID(petID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPetNotFound
		}
		return err
	}

	// Проверяем права доступа
	if pet.OwnerID != userID {
		return ErrPermissionDenied
	}

	// Удаляем питомца
	return s.petRepo.Delete(petID)
}

func (s *PetService) UploadPetPhoto(petID uuid.UUID, file *multipart.FileHeader) (string, error) {
	_, err := s.petRepo.FindByID(petID)
	if err != nil {
		return "", fmt.Errorf("pet not found: %w", err)
	}

	// Генерируем уникальное имя файла
	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("pet_%d_%d%s", petID, time.Now().Unix(), ext)
	uploadPath := filepath.Join("uploads", "pets", newFilename)

	if err := os.MkdirAll(filepath.Dir(uploadPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	if err := saveUploadedFile(file, uploadPath); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	photoURL := "/" + uploadPath
	if err := s.petRepo.UpdatePhoto(petID, photoURL); err != nil {
		os.Remove(uploadPath)
		return "", fmt.Errorf("failed to update pet photo: %w", err)
	}

	return photoURL, nil
}


func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}