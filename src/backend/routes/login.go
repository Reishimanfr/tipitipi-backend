package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	// password := ctx.PostForm("password")

	// TODO: implement password hashing
	var adminUser core.AdminUser

	if err := h.Db.Where("username = ?", username).First(&adminUser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve entry from database",
		})
		ctx.Abort()
		return
	}

}
