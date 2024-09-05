package routes

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UNUSED
func (h *Handler) uploadPostAttachments(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while retrieving multipart form data",
		})
		return
	}

	stringId := c.Param("id")
	_, err = strconv.Atoi(stringId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid post ID provided",
		})
		return
	}

	files := form.File["upload[]"]

	if len(files) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No files provided",
		})
		return
	}

	for _, file := range files {
		filename := filepath.Base(file.Filename)

		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Error while saving the uploaded file",
			})
			return
		}
	}

	c.Status(http.StatusOK)
}
