package utils

import (
	"errors"
	"time"
	"fmt"
	"FurryTrack/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateToken(userID uuid.UUID, role models.Role, secret string) (string, error) {
    if !role.IsValid() {
        return "", fmt.Errorf("invalid role: %s", role)
    }

    claims := jwt.MapClaims{
        "userID": userID.String(),
        "role":   role,
        "exp":    time.Now().Add(time.Hour * 24).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}


func ValidateToken(tokenString string, secret string) (uuid.UUID, error) {
	if tokenString == "" {
		return uuid.Nil, errors.New("empty token string")
	}

	// Парсим токен с секретным ключом
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("token parsing failed: %w", err)
	}

	// Проверяем валидность токена
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims format")
	}

	// Проверяем срок действия токена
	if err := checkTokenExpiration(claims); err != nil {
		return uuid.Nil, err
	}

	// Извлекаем userID из claims
	userID, err := extractUserID(claims)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

// checkTokenExpiration - проверяет срок действия токена
func checkTokenExpiration(claims jwt.MapClaims) error {
	exp, ok := claims["exp"]
	if !ok {
		return errors.New("token missing expiration claim")
	}

	var expTime int64
	switch v := exp.(type) {
	case float64:
		expTime = int64(v)
	case int64:
		expTime = v
	default:
		return errors.New("invalid exp claim format")
	}

	if time.Now().Unix() > expTime {
		return errors.New("token expired")
	}

	return nil
}

// extractUserID извлекает и парсит userID из claims
func extractUserID(claims jwt.MapClaims) (uuid.UUID, error) {
	userIDValue, ok := claims["userID"]
	if !ok {
		userIDValue, ok = claims["user_id"]
		if !ok {
			return uuid.Nil, errors.New("token missing userID claim")
		}
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		return uuid.Nil, errors.New("userID must be a string")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid userID format: %w", err)
	}

	return userID, nil
}
