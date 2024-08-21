package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateToken(username string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	jwtSecret := os.Getenv("JWT_SECRET")

	fmt.Println(jwtSecret)

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "strona-fundacja-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return os.Getenv("JWT_SECRET"), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("token")

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			ctx.Abort()
			return
		}

		claims, err := ValidateToken(tokenString)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}

		ctx.Set("username", claims.Username)
		ctx.Next()
	}
}
