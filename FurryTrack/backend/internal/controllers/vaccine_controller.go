package controllers

import (
	"FurryTrack/internal/models"
	"FurryTrack/internal/services"
	"net/http"
	"FurryTrack/pkg/middleware"
	//"log"
	"github.com/gin-gonic/gin"
)

type VaccineController struct {
	service services.VaccineService
}

func NewVaccineController(service services.VaccineService) *VaccineController {
	return &VaccineController{service: service}
}

// CreateVaccine создает новую вакцину в справочнике
func (c *VaccineController) CreateVaccine(ctx *gin.Context) {

	// Проверяем, что пользователь админ
	userRole, exists := ctx.Get(middleware.RoleKey)
	if !exists || userRole != middleware.RoleAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Only admin can create vaccines"})
		return
	}

	// Парсим входные данные
	var vaccine models.Vaccine
	if err := ctx.ShouldBindJSON(&vaccine); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Создаем вакцину
	if err := c.service.AddVaccine(&vaccine); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем успешный ответ
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   vaccine,
	})
}


func (c *VaccineController) GetAllVaccines(ctx *gin.Context) {
	// Проверяем, что пользователь админ
	userRole, exists := ctx.Get(middleware.RoleKey)
	if !exists || userRole != middleware.RoleAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Only admin can create vaccines"})
		return
	}
	vaccines, err := c.service.GetAllVaccines()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, vaccines)
}
