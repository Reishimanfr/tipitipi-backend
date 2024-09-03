package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CredentialsRequestBody struct {
	Password string `json:"password" binding:"required,min=8"`
	Username string `json:"username" binding:"required"`
}

func (h *Handler) changePassword(c *gin.Context) {
	var newCredentials CredentialsRequestBody

	if err := c.ShouldBindWith(&newCredentials, binding.JSON); err != nil {
		h.Log.Error("Failed to bind json data", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "Malformed or invalid JSON",
		})
		return
	}

	hashSalt, err := h.A.GenerateHash([]byte(newCredentials.Password), nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error while hashing the new password",
			"error":   err.Error(),
		})
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
			"error":   result.Error.Error(),
			"message": "Admin user managed to try to change their password while not existing",
		})
		return
	}

	update := h.Db.Model(&core.AdminUser{}).Where("id = 1").Updates(core.AdminUser{
		Hash:     string(hashSalt.Hash),
		Salt:     string(hashSalt.Salt),
		Username: newCredentials.Username,
	})

	if update.Error != nil {
		h.Log.Error("Error while updating admin password", zap.Error(update.Error))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error while updating password",
			"error":   update.Error.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
