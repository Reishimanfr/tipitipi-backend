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
	Title   string                  `json:"title" binding:"required"`
	Content string                  `json:"content" binding:"required"`
	Images  []*RequestBodyImageData `json:"images"`
}

func (h *Handler) edit(c *gin.Context) {
	var data EditRequestBody

	if err := c.BindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "Malformed or invalid JSON",
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

	result := h.Db.Where("id = ?", id).UpdateColumns(newBlog)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": "Error while updating post",
		})
	}

	c.JSON(http.StatusOK, newBlog)
}
