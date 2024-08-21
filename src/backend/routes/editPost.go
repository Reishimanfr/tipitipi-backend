package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type EditRequestBody struct {
	Title   string                 `json:"title"`
	Content string                 `json:"content"`
	Images  []RequestBodyImageData `json:"images"`
}

func (h *Handler) EditBlogPost(ctx *gin.Context) {
	var data EditRequestBody

	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "failed to bind json data: " + err.Error(),
		})
		return
	}

	var postRecord *core.BlogPost
	stringId := ctx.Param("id")
	id, err := strconv.Atoi(stringId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "an invalid post id was provided",
		})
		return
	}

	h.Db.Where("id = ?", id).First(&postRecord)

	if postRecord == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "post not found",
		})
		return
	}

	var newBlog core.BlogPost

	if strings.Trim(data.Title, "") != "" && data.Title != postRecord.Title {
		newBlog.Title = data.Title
	}

	if data.Content != "" && data.Content != postRecord.Content {
		newBlog.Content = data.Content
	}

	if len(data.Images) > 0 {
		// TODO
	}

	if newBlog == (core.BlogPost{}) {
		ctx.JSON(http.StatusNotModified, postRecord)
		return
	}

	newBlog.Edited_At = time.Now()

	h.Db.Where("id = ?", id).UpdateColumns(newBlog)

	ctx.JSON(http.StatusOK, newBlog)
}
