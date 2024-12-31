package middleware

import (
	"bash06/tipitipi-backend/core"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB, l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing authorization header",
			})

			return
		}

		if authHeader[:6] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unexpected token type",
			})
			return
		}

		token := authHeader[7:]

		var found bool

		if err := db.Model(&core.Token{}).Where("token = ?", token).Select("count(*) > 0").Find(&found).Error; err != nil {
			l.Error("Failed to check if opaque token exists in db", zap.Error(err))

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong while processing your request",
			})
			return
		}

		if !found {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}
	}
}
