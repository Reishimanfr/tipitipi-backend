package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

type RequestBodyImageData struct {
	Id       uint8  `json:"id" binding:"required"`
	Path     string `json:"path" binding:"required"`
	Position string `json:"position" binding:"required"`
}

type CreateRequestBody struct {
	Title   string                  `json:"title" binding:"required,max=255"`
	Content string                  `json:"content" binding:"required"`
	Images  []*RequestBodyImageData `json:"images" binding:"max=20" `
}

func validateImages(data []*RequestBodyImageData, c *gin.Context) bool {
	var seenIds []uint8

	for _, v := range data {
		if slices.Contains(seenIds, v.Id) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Duplicate image ids",
			})
			return false
		}

		seenIds = append(seenIds, v.Id)
	}

	return true
}

func (h *Handler) create(c *gin.Context) {
	var data CreateRequestBody

	if err := c.ShouldBindWith(&data, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "Malformed or invalid JSON",
		})
		return
	}

	validateImages(data.Images, c)

	var count int64
	result := h.Db.Model(&core.BlogPost{}).Where("title = ?", data.Title).Count(&count)

	if result.Error != nil {
		h.Log.Error("Error while checking for post in database", zap.Error(result.Error))
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error":   result.Error.Error(),
			"message": "Error while cheching if post exists",
		})
		return
	}

	if count > 0 {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "Post with this title already exists",
		})
		return
	}

	newPost := &core.BlogPost{
		Title:      data.Title,
		Content:    data.Content,
		Created_At: time.Now().Unix(),
		Edited_At:  time.Now().Unix(),
		Images:     "", // TODO
	}

	result = h.Db.Create(&newPost)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": "Error while creating database record",
		})
		return
	}

	c.JSON(http.StatusOK, newPost)
}
