package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	var postRecord *core.BlogPost

	h.Db.Where("id = ?", id).First(&postRecord)

	if postRecord == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	result := h.Db.Delete(core.BlogPost{ID: id})
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": "Error while deleting record from database",
		})
		return
	}

	c.Status(http.StatusOK)
}
