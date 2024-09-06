package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type CreatePostBody struct {
	Title   string              `form:"title"`
	Content string              `form:"content"`
	Files   []*core.ImageRecord `form:"files[]"`
}

const (
	ATTACHMENT_SIZE_LIMIT = 10 << 20 // 10 MB
)

func (h *Handler) create(c *gin.Context) {
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

	var exists bool
	h.Db.Model(&core.BlogPost{}).Select("count(*) > 0").Where("title = ?", title).Find(&exists)

	if exists {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error": "Blog post with this title already exists",
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
			"error": "Expected some content, but got nothing",
		})
		return
	}

	files := form.File["files[]"]

	processedFiles := make([]core.ImageRecord, 0, len(files))

	var nextID int
	h.Db.Model(&core.BlogPost{}).Select("COALESCE(MAX(id), 0) + 1").Scan(&nextID)

	curPath, err := os.Getwd()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while getting the current work directory",
		})
		return
	}

	for idx, file := range files {
		if file.Size > ATTACHMENT_SIZE_LIMIT {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": fmt.Sprintf("Attachment %v is too large (limit: %v bytes, got %v)", file.Filename, ATTACHMENT_SIZE_LIMIT, file.Size),
			})
			return
		}

		extension := filepath.Ext(file.Filename)
		internalFilename := fmt.Sprintf("%v-%v%v", nextID, idx+1, extension)
		files[idx].Filename = internalFilename

		processedFiles = append(processedFiles, core.ImageRecord{
			Filename:   internalFilename,
			Path:       filepath.Join(curPath, "../backend/assets", internalFilename),
			BlogPostID: nextID,
		})
	}

	post := core.BlogPost{
		ID:         nextID,
		Created_At: time.Now().Unix(),
		Edited_At:  time.Now().Unix(),
		Title:      title,
		Content:    content,
		Images:     processedFiles,
	}

	result := h.Db.Create(&post)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": "Error while creating post record",
		})
		return
	}

	for _, v := range files {
		err := c.SaveUploadedFile(v, "assets/"+v.Filename)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Error while downloading attached file",
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post added successfully",
	})
}
