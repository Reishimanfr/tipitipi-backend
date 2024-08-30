package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (h *Handler) changePassword(c *gin.Context) {
	password := c.PostForm("password")

	if len(password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 8 characters",
		})
		c.Abort()
		return
	}

	hashSalt, err := h.A.GenerateHash([]byte(password), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while hashing the new password",
		})
		c.Abort()
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	var adminUser core.AdminUser

	result := h.Db.WithContext(ctx).Take(&adminUser, 1)

	// This shouldn't really happen since that would mean the user
	// managed to login without having an account created and changed
	// their password with some voodoo fucking magic
	if result.Error == gorm.ErrRecordNotFound {
		h.Log.Info("Admin user managed to change their password while not existing")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "This user doesn't exist",
		})
		return
	}

	update := h.Db.Model(&core.AdminUser{}).Where("id = 1").Updates(core.AdminUser{
		Hash: string(hashSalt.Hash),
		Salt: string(hashSalt.Salt),
	})

	if update.Error != nil {
		h.Log.Error("Error while updating admin password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while updating password",
		})
		c.Abort()
		return
	}

	c.Status(http.StatusOK)
}
