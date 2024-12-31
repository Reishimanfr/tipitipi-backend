package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"admin":   true,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Missing authorization header",
			})
			return
		}

		if !strings.HasPrefix(auth, "Bearer") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Expected a bearer token, but got something else",
			})
			return
		}

		authSplit := strings.Split(auth, " ")

		if len(authSplit) < 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Missing JWT token in authorization header",
			})
			return
		}

		tokenString := authSplit[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !claims["admin"].(bool) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
