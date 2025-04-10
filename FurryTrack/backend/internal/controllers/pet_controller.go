package controllers

import (
	"FurryTrack/internal/models"
	"FurryTrack/internal/services"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"errors"
	"time"

)

type PetController struct {
	petService     services.PetService
	vaccineService services.VaccineService
}

func NewPetController(petService services.PetService, vaccineService services.VaccineService) *PetController {
	return &PetController{
		petService:     petService,
		vaccineService: vaccineService,
	}
}

// CreatePet — добавление нового питомца
func (c *PetController) CreatePet(ctx *gin.Context) {
	var pet models.Pet
	if err := ctx.ShouldBindJSON(&pet); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Получаем userID из контекста
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	pet.OwnerID = userID.(uuid.UUID)

	if err := c.petService.CreatePet(&pet); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, pet)
}

// GetPet — получение информации о питомце
func (c *PetController) GetPet(ctx *gin.Context) {
	petIDStr := ctx.Param("id")
	petID, err := uuid.Parse(petIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	ownerID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	pet, err := c.petService.GetPetByID(petID, ownerID.(uuid.UUID))
	if err != nil {
		switch err.Error() {
		case "pet not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "access denied":
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, pet)
}

func (c *PetController) GetUserPets(ctx *gin.Context) {
    // Получаем userID из контекста
    userIDVal, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    userID, ok := userIDVal.(uuid.UUID) 
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
        return
    }

    pets, err := c.petService.GetUserPets(userID) 
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"data": pets})
}

func (c *PetController) UpdatePet(ctx *gin.Context) {
    // Получаем ID питомца из URL
    petIDstr := ctx.Param("pet_id")
    
	petID, err := uuid.Parse(petIDstr);
    // Проверяем валидность UUID
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID format"})
        return
    }

    // Получаем ID пользователя из токена
    userID, ok := ctx.Get("userID")
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Получаем данные для обновления
    var updateData struct {
        Name      *string `json:"name,omitempty"`
        Breed     *string `json:"breed,omitempty"`
        BirthDate *time.Time `json:"birth_date,omitempty"`
    }
    
    if err := ctx.ShouldBindJSON(&updateData); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    //  Используем сервис для проверки прав и обновления
    updatedPet, err := c.petService.UpdatePet(
        petID,
        userID.(uuid.UUID),
        updateData.Name,
        updateData.Breed,
        updateData.BirthDate,
    )
    
    if err != nil {
        if errors.Is(err, services.ErrPetNotFound) {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
        } else if errors.Is(err, services.ErrPermissionDenied) {
            ctx.JSON(http.StatusForbidden, gin.H{"error": "No permission to update this pet"})
        } else {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"data": updatedPet})
}

func handleDeleteError(ctx *gin.Context, err error) {
	if errors.Is(err, services.ErrPetNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Питомец не найден"})
	} else if errors.Is(err, services.ErrPermissionDenied) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Нет прав для удаления"})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
	}
}

func (c *PetController) DeletePet(ctx *gin.Context) {
	//  Получаем и валидируем ID питомца
	petID, err := uuid.Parse(ctx.Param("pet_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID питомца"})
		return
	}

	// Получаем ID пользователя из аутентификации
	userID, ok := ctx.Get("userID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}

	// Вызываем сервис
	if err := c.petService.DeletePet(petID, userID.(uuid.UUID)); err != nil {
		handleDeleteError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// UploadPetPhoto загружает фото для питомца
func (pc *PetController) UploadPetPhoto(c *gin.Context) {
	// Парсим UUID питомца
	petID, err := uuid.Parse(c.Param("pet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID format"})
		return
	}

	// Получаем файл
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызываем сервис
	photoURL, err := pc.petService.UploadPetPhoto(petID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"photo_url": photoURL,
		"message":   "Photo uploaded successfully",
	})
}