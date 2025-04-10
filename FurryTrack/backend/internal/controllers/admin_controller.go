package controllers

import (
	"net/http"
	"gorm.io/gorm"
	"FurryTrack/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminController struct {
	service *services.AdminService
	authService *services.AuthService
}

func NewAdminController(adminService *services.AdminService, authService *services.AuthService) *AdminController {
    return &AdminController{
        service: adminService,
        authService: authService,
    }
}


func (c *AdminController) GetAllUsers(ctx *gin.Context) {
	users, err := c.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"data": users})
}

// Блокировка пользователя
func (c *AdminController) BanUser(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	adminID := ctx.MustGet("userID").(uuid.UUID)
	
	var request struct {
		Reason string `json:"reason" binding:"required"`
	}
	
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := c.service.BanUser(adminID, userID, request.Reason); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "User banned successfully"})
}

// UnbanUser разблокирует заблокированного пользователя
func (c *AdminController) UnbanUser(ctx *gin.Context) {
    userID, err := uuid.Parse(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
        return
    }

    adminID := ctx.MustGet("userID").(uuid.UUID)

    var request struct {
        Reason string `json:"reason" binding:"required,min=10,max=500"`
    }
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Вызов сервиса
    if err := c.service.UnbanUser(adminID, userID, request.Reason); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Пользователь успешно разблокирован",
        "details": gin.H{
            "user_id": userID,
            "by_admin": adminID,
        },
    })
}

// DeleteUser удаляет пользователя из системы
func (c *AdminController) DeleteUser(ctx *gin.Context) {
    userID, err := uuid.Parse(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
        return
    }

    adminID := ctx.MustGet("userID").(uuid.UUID)

    var request struct {
        Reason string `json:"reason" binding:"required,min=10,max=500"`
        HardDelete bool `json:"hard_delete"` // Флаг полного удаления
    }
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := c.service.DeleteUser(adminID, userID, request.Reason); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Пользователь успешно удален",
        "mode":    map[bool]string{true: "hard delete", false: "soft delete"}[request.HardDelete],
    })
}

// GetUserByID возвращает полную информацию о пользователе
func (c *AdminController) GetUserByID(ctx *gin.Context) {
    // 1. Парсинг ID пользователя
    userID, err := uuid.Parse(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
        return
    }

    // Проверка прав администратора
    adminID := ctx.MustGet("userID").(uuid.UUID)
    if !c.authService.IsAdmin(adminID) {
        ctx.JSON(http.StatusForbidden, gin.H{"error": "Требуются права администратора"})
        return
    }

    user, err := c.service.GetUserByID(userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
        } else {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    response := gin.H{
        "id":         user.ID,
        "email":      user.Email,
        "name":       user.Username,
        "role":       user.Role,
        "created_at": user.CreatedAt,
        "pets_count": len(user.Pets),
    }

    ctx.JSON(http.StatusOK, response)
}