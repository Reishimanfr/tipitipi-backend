package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) edit(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid post ID provided",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while getting multipart form data",
		})
		return
	}

	titleFields := form.Value["title"]
	contentFields := form.Value["content"]
	files := form.File["files[]"]

	if len(titleFields) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Expected a title, but got nothing",
		})
		return
	}

	title := titleFields[0]
	if title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Expected a title, but got nothing",
		})
		return
	}

	if len(contentFields) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Expected some content, but got nothing",
		})
		return
	}

	content := contentFields[0]
	if content == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Expected a title, but got nothing",
		})
		return
	}

	processedFiles := make([]core.ImageRecord, 0, len(files))
	seenNames := make(map[string]struct{}, len(files))

	curPath, err := os.Getwd()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while getting the current work directory",
		})
		return
	}

	for _, file := range files {
		if _, seen := seenNames[file.Filename]; seen {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Duplicate asset filename found in request body",
			})
			return
		}

		if file.Size > ATTACHMENT_SIZE_LIMIT {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": fmt.Sprintf("Attachment %v is too large (limit: %v bytes, got %v)", file.Filename, ATTACHMENT_SIZE_LIMIT, file.Size),
			})
			return
		}

		seenNames[file.Filename] = struct{}{}
		processedFiles = append(processedFiles, core.ImageRecord{
			Filename:   file.Filename,
			Path:       filepath.Join(curPath, "../backend/assets", file.Filename),
			BlogPostID: id,
		})
	}

	var existingPost core.BlogPost
	if err := h.Db.Where("id = ?", id).First(&existingPost).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while fetching existing post record",
		})
		return
	}

	for _, oldImage := range existingPost.Images {
		if err := os.Remove(oldImage.Path); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Error while removing old image",
			})
			return
		}

		if err := h.Db.Model(&core.BlogPost{}).Where("id = ?", id).UpdateColumn("Images", gorm.Expr("array_remove(Images, ?)", oldImage)).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Error while removing image entry from database",
			})
			return
		}
	}

	updatedPost := core.BlogPost{
		ID:        id,
		Edited_At: time.Now().Unix(),
		Title:     title,
		Content:   content,
		Images:    processedFiles,
	}

	tx := h.Db.Begin()

	if err := tx.Model(&core.BlogPost{}).Where("id = ?", id).Updates(updatedPost).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Transaction failed with error",
		})
		return
	}

	var oldImages []core.ImageRecord
	if err := tx.Where("blog_post_id = ?", id).Find(&oldImages).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Transaction failed with error",
		})
		return
	}

	fmt.Println(oldImages)

	for _, img := range oldImages {
		if err := os.Remove(img.Path); err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Error while deleting old asset file",
			})
			return
		}
	}

	if err := tx.Where("blog_post_id = ?", id).Delete(&core.ImageRecord{}).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Transaction failed with error",
		})
		return
	}

	if err := tx.Create(&updatedPost.Images).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Transaction failed with error",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Commiting transaction failed with error",
		})
		return
	}

	for _, file := range files {
		err := c.SaveUploadedFile(file, "assets/"+file.Filename)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Error while downloading attached file",
			})
		}
	}

	c.Status(http.StatusOK)
}
