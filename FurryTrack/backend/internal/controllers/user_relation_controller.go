package controllers

import (
	"net/http"
	//"FurryTrack/internal/models"
	"FurryTrack/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	//"log"
	"strings"
	"fmt"
)

type UserRelationController struct {
	service *services.UserRelationService
}

func NewUserRelationController(service *services.UserRelationService) *UserRelationController {
	return &UserRelationController{service: service}
}

func (c *UserRelationController) Follow(ctx *gin.Context) {
    // Получаем ID текущего пользователя из контекста
    currentUserID, ok := ctx.Get("userID")
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    followerID, ok := currentUserID.(uuid.UUID)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
        return
    }

    // Получаем ID целевого пользователя из параметров URL
    targetUserIDStr := ctx.Param("id")
    followingID, err := uuid.Parse(targetUserIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID format"})
        return
    }

    // Проверяем, что пользователь не пытается подписаться на себя
    if followerID == followingID {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
        return
    }

    if err := c.service.Follow(followerID, followingID); err != nil {
        // Проверяем на дубликат подписки
        if strings.Contains(err.Error(), "duplicate") {
            ctx.JSON(http.StatusConflict, gin.H{"error": "Already following this user"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Followed successfully",
        "data": gin.H{
            "follower_id":  followerID,
            "following_id": followingID,
        },
    })
}

func (c *UserRelationController) GetFollowing(ctx *gin.Context) {
    // Получаем ID пользователя из контекста
    userID, ok := ctx.Get("userID")
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "error": "User not authenticated",
            "details": "userID missing in request context",
        })
        return
    }

    uuidUserID, ok := userID.(uuid.UUID)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Invalid user ID type",
            "expected": "uuid.UUID",
            "got": fmt.Sprintf("%T", userID),
        })
        return
    }

    
    // Получаем данные из сервиса
    following, err := c.service.GetFollowing(uuidUserID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get following list",
            "details": err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "following": following,
            "count": len(following),
        },
        "meta": gin.H{
            "user_id": uuidUserID,
        },
    })
}

func (c *UserRelationController) GetFollowers(ctx *gin.Context) {
    // Получаем ID пользователя из контекста
    userID, ok := ctx.Get("userID")
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "error": "User not authenticated",
            "details": "userID missing in request context",
        })
        return
    }

    uuidUserID, ok := userID.(uuid.UUID)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Invalid user ID type",
            "expected": "uuid.UUID",
            "got": fmt.Sprintf("%T", userID),
        })
        return
    }

    
    // Получаем данные из сервиса
    followers, err := c.service.GetFollowers(uuidUserID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get followers list",
            "details": err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "followers": followers,
            "count": len(followers),
        },
        "meta": gin.H{
            "user_id": uuidUserID,
        },
    })
}