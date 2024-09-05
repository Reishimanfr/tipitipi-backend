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
	Title   string                  `json:"title" binding:"required"`
	Content string                  `json:"content" binding:"required"`
	Images  []*RequestBodyImageData `json:"images"`
}

func (h *Handler) edit(c *gin.Context) {
	body := new(EditRequestBody)

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "Malformed or invalid JSON",
		})
		return
	}

	postRecord := new(core.BlogPost)

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

	newBlog := new(core.BlogPost)

	if strings.Trim(body.Title, "") != "" && body.Title != postRecord.Title {
		newBlog.Title = body.Title
	}

	if body.Content != "" && body.Content != postRecord.Content {
		newBlog.Content = body.Content
	}

	if len(body.Images) > 0 {
		// TODO
	}

	if newBlog.Title == "" && newBlog.Content == "" && len(newBlog.Images) == 0 {
		c.JSON(http.StatusNotModified, postRecord)
		return
	}

	newBlog.Edited_At = time.Now().Unix()

	result := h.Db.Where("id = ?", id).UpdateColumns(newBlog)

	if result.Error != nil {
		h.Log.Error("Error while updating post record", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": "Error while updating post",
		})
	}

	c.JSON(http.StatusOK, newBlog)
}
