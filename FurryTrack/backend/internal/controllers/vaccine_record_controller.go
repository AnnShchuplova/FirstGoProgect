package controllers

import (
	"FurryTrack/internal/models"
	"FurryTrack/internal/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VaccineRecordController struct {
	service services.VaccineRecordService
}

func NewVaccineRecordController(service services.VaccineRecordService) *VaccineRecordController {
	return &VaccineRecordController{service: service}
}


func (c *VaccineRecordController) AddVaccineRecord(ctx *gin.Context) {
	// Получаем pet_id из URL
	petID, err := uuid.Parse(ctx.Param("pet_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	// Определяем структуру запроса
	var request struct {
		VaccineName string `json:"vaccine_name" binding:"required"`
		Date        string `json:"date" binding:"required"`
		Clinic      string `json:"clinic,omitempty"`
	}

	// Парсим JSON тело запроса
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Парсим дату
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":    "Invalid date format",
			"expected": "YYYY-MM-DD",
			"received": request.Date,
			"example":  "2023-12-31",
		})
		return
	}

	// Получаем userID из контекста
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userRole := ctx.MustGet("role").(models.Role) 
	log.Printf("PetID: %v", petID)

	// Вызываем сервис
	record, err := c.service.AddVaccineRecord(
		userID.(uuid.UUID),
		petID,
		request.VaccineName,
		date,
		request.Clinic,
		userRole,
	)

	// Обрабатываем ошибки
	if err != nil {
		switch err.Error() {
		case "access denied":
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case "pet not found", "vaccine not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Возвращаем ответ
	ctx.JSON(http.StatusCreated, record)
}

func (c *VaccineRecordController) GetPetVaccineHistory(ctx *gin.Context) {
	petID, err := uuid.Parse(ctx.Param("pet_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	records, err := c.service.GetPetVaccinationHistory(petID)
	if err != nil {
		switch err.Error() {
		case "access denied":
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case "pet not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, records)
}
