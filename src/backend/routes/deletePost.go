package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (h *Handler) delete(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post id",
		})
		return
	}

	postRecord := new(core.BlogPost)

	result := h.Db.Where("id = ?", id).First(&postRecord)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		h.Log.Error("Error while searching for a post record", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": "Error while looking for record in database",
		})
		return
	}

	if postRecord == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	result = h.Db.Delete(core.BlogPost{ID: id})
	if result.Error != nil {
		h.Log.Error("Error while deleting record from database", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": "Error while deleting record from database",
		})
		return
	}

	c.Status(http.StatusOK)
}
