package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (h *Handler) post(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)
	images := c.Query("images") == "true"

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post id",
		})
		return
	}

	postRecord := new(core.BlogPost)
	result := new(gorm.DB)

	if images {
		result = h.Db.Preload("Attachments").Where("id = ?", id).First(&postRecord)
	} else {
		result = h.Db.Where("id = ?", id).First(&postRecord)
	}

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		h.Log.Error("Error while searching for a post record", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": "Error while checking if post exists in database",
		})
		return
	}

	if postRecord == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Post with this ID doesn't exist",
		})
		return
	}

	c.JSON(http.StatusOK, postRecord)
}
