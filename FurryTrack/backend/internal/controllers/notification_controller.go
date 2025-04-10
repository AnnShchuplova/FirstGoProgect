package controllers

import (
    "net/http"
    
    "FurryTrack/internal/services"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type NotificationController struct {
    service *services.NotificationService
}

func NewNotificationController(service *services.NotificationService) *NotificationController {
    return &NotificationController{service: service}
}

func (c *NotificationController) GetNotifications(ctx *gin.Context) {
    userID, err := uuid.Parse(ctx.MustGet("userID").(string))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }
    
    notifications, err := c.service.GetUserNotifications(userID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{"data": notifications})
}

func (c *NotificationController) MarkAsRead(ctx *gin.Context) {
    userID, err := uuid.Parse(ctx.MustGet("userID").(string))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }
    
    notificationID, err := uuid.Parse(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
        return
    }
    
    if err := c.service.MarkAsRead(notificationID, userID); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (c *NotificationController) GetUnreadCount(ctx *gin.Context) {
    userID, err := uuid.Parse(ctx.MustGet("userID").(string))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }
    
    count, err := c.service.GetUnreadCount(userID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{"count": count})
}