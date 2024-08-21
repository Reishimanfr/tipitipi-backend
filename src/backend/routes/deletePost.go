package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteBlogPost(ctx *gin.Context) {
	stringId := ctx.Param("id")
	id, err := strconv.Atoi(stringId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid post id provided",
		})
		return
	}

	var postRecord *core.BlogPost

	h.Db.Where("id = ?", id).First(&postRecord)

	if postRecord == nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	h.Db.Delete(core.BlogPost{ID: id})

	ctx.JSON(http.StatusOK, nil)
}
