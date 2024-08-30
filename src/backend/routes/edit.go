package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EditRequestBody struct {
	Title   string                 `json:"title"`
	Content string                 `json:"content"`
	Images  []RequestBodyImageData `json:"images"`
}

func (h *Handler) edit(c *gin.Context) {
	var data EditRequestBody

	if err := c.BindJSON(&data); err != nil {
		h.Log.Error("Failed to bind json data", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Malformed or invalid JSON",
		})
		return
	}

	var postRecord *core.BlogPost
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post id",
		})
		return
	}

	h.Db.Where("id = ?", id).First(&postRecord)

	if postRecord == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
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
		c.JSON(http.StatusNotModified, postRecord)
		return
	}

	newBlog.Edited_At = time.Now().Unix()

	h.Db.Where("id = ?", id).UpdateColumns(newBlog)

	c.JSON(http.StatusOK, newBlog)
}
