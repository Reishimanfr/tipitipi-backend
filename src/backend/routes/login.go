package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"bash06/strona-fundacja/src/backend/middleware"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (h *Handler) AdminLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var adminData core.AdminUser

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		var err error
		result := h.Db.WithContext(ctx).First(&adminData, "id = 1")

		if result.Error == gorm.ErrRecordNotFound {
			hashSalt, genErr := h.A.GenerateHash([]byte(password), nil)

			if genErr == nil {
				adminData = core.AdminUser{
					ID:       1,
					Username: "admin",
					Hash:     string(hashSalt.Hash),
					Salt:     string(hashSalt.Salt),
				}
				h.Db.Create(&adminData)
			}
			err = genErr
		} else {
			err = result.Error
		}

		errChan <- err
	}()

	err := <-errChan

	if err != nil {
		h.Log.Error("Error during password operation", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	err = h.A.Compare([]byte(adminData.Hash), []byte(adminData.Salt), []byte(password))

	if username != adminData.Username || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		c.Abort()
		return
	}

	token, err := middleware.GenerateJWT(username, true)
	if err != nil {
		h.Log.Error("Error while generating JWT token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not generate token",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
