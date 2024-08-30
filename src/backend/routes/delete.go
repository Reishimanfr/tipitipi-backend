package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) delete(ctx *gin.Context) {
	stringId := ctx.Param("id")
	id, err := strconv.Atoi(stringId)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post id",
		})
		return
	}

	var postRecord *core.BlogPost

	h.Db.Where("id = ?", id).First(&postRecord)

	if postRecord == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	h.Db.Delete(core.BlogPost{ID: id})
	ctx.Status(http.StatusOK)
}
