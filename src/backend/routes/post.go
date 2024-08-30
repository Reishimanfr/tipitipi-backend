package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) post(ctx *gin.Context) {
	stringId := ctx.Param("id")
	id, err := strconv.Atoi(stringId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post id",
		})
		return
	}

	var postRecord *core.BlogPost

	h.Db.Where("id = ?", id).First(&postRecord)

	if postRecord == nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	ctx.JSON(http.StatusOK, postRecord)
}
