package controllers

import (
	//"FurryTrack/internal/models"
	"FurryTrack/internal/services"
	"FurryTrack/pkg/utils"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: *authService}
}

// RegisterUser - регистрация нового пользователя
func (c *AuthController) RegisterUser(ctx *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		IsAdmin  bool    `gorm:"default:false"`
	}

	// Валидация входных данных
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid input format" + err.Error())
		return
	}

	// Создаем пользователя через сервис
	user, err := c.authService.Register(
		input.Username,
		input.Email,
		input.Password, 
		input.IsAdmin,
	)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user already exists" {
			status = http.StatusConflict
		}
		utils.RespondWithError(ctx, status, err.Error())
		return
	}

	utils.RespondWithJSON(ctx, http.StatusCreated, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	})
}

// LoginUser - аутентификация пользователя
func (c *AuthController) LoginUser(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid input format")
		return
	}

	token, user, err := c.authService.Login(input.Email, input.Password)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// GetUserProfile - получение профиля пользователя
func (c *AuthController) GetUserProfile(ctx *gin.Context) {
	// Получаем userID из контекста
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.RespondWithError(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.RespondWithError(ctx, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Получаем профиль через сервис
	user, err := c.authService.GetUserProfile(userUUID)
	if err != nil {
		utils.RespondWithError(ctx, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondWithJSON(ctx, http.StatusOK, user)
}


func (c *AuthController) GetUserByEmail(ctx *gin.Context) {
    email := ctx.Param("email")
    if email == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter is required"})
        return
    }

    user, err := c.authService.GetUserByEmail(email)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    response := gin.H{
        "data": gin.H{
            "ID":        user.ID,
            "email":     user.Email,
            "username": user.Username,
            "createdAt": user.CreatedAt,
        },
    }

    ctx.JSON(http.StatusOK, response)
}