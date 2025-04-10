package utils

import "github.com/gin-gonic/gin"

// RespondWithError - обработчик ошибок
func RespondWithError(ctx *gin.Context, statusCode int, errorMessage string) {
	ctx.AbortWithStatusJSON(statusCode, gin.H{
		"success": false,
		"error":   errorMessage,
	})
}

// RespondWithJSON - обработчик успешных ответов
func RespondWithJSON(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}
