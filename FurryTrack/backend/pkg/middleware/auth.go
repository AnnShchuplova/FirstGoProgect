package middleware

import (
	//"FurryTrack/pkg/utils"
	"net/http"
	"strings"
	"log"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"fmt"
	"FurryTrack/internal/models"
)

const (
	RoleKey   = "role"
	UserIDKey = "userID"
)

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
	RoleVet   Role = "VET"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser, RoleVet:
		return true
	}
	return false
}

func GinAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		log.Printf("Authorization header: %s", authHeader)
		
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token not found"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Получаем userID
		userIDStr, ok := claims["userID"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid userID in token"})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid userID format"})
			return
		}

		// Получаем роль и преобразуем в models.Role
		roleStr, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid role in token"})
			return
		}

		role := models.Role(roleStr)
		if !role.IsValid() {  // Убедитесь, что models.Role имеет метод IsValid()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid role value: " + roleStr})
			return
		}

		// Устанавливаем в контекст
		c.Set("userID", userID)
		c.Set("role", role)  // Теперь сохраняем models.Role
		c.Next()
	}
}

