package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) post(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post id",
		})
		return
	}

	var postRecord *core.BlogPost

	h.Db.Where("id = ?", id).First(&postRecord)

	if postRecord == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, postRecord)
}
