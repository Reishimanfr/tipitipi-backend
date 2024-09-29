package routes

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

const (
	IMAGE_UPLOAD_AMOUNT_LIMIT = 5
)

var (
	acceptedMimeTypes = []string{
		"image/png",
		"image/jpeg",
		"image/webp",
	}
)

func (h *Handler) createGallery(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while parsing form",
		})
		return
	}

	files := form.File["files[]"]

	if len(files) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No images provided",
		})
		return
	}

	if len(files) > IMAGE_UPLOAD_AMOUNT_LIMIT {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Too many images provided at once",
		})
		return
	}

	declinedFiles := make(map[string]string, 0)

	for _, file := range files {
		header := file.Header["Content-Type"]
		sizeAsMiB := (file.Size / 1000) / 1000

		if sizeAsMiB > 10 {
			declinedFiles[file.Filename] = "File is too large (max 10MiB per file)"
		}

		if !slices.Contains(acceptedMimeTypes, header[0]) {
			declinedFiles[file.Filename] = fmt.Sprintf("Invalid MIME type. Expected one of %v, but got %v instead", acceptedMimeTypes, header[0])
		}
	}

	if len(declinedFiles) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Some files don't meet the constrains and were rejected",
			"files": declinedFiles,
		})
		return
	}

	// We know all files are valid images so we can download everything without worrying
	for _, file := range files {
		if err := c.SaveUploadedFile(file, "gallery/"+file.Filename); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Error while downloading the uploaded file",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All images uploaded successfully",
	})

}
