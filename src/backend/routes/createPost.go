package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) create(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while parsing multipart form",
		})
		return
	}

	title := c.PostForm("title")

	if title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No title provided",
		})
		return
	}

	content := c.PostForm("content")

	if content == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No content provided",
		})
		return
	}

	var exists bool
	h.Db.Model(&core.BlogPost{}).Select("count(*) > 0").Where("title = ?", title).Find(&exists)

	if exists {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error": "Post with this title already exists",
		})
		return
	}

	files := form.File["files[]"]
	processedFiles := make([]core.AttachmentRecord, 0, len(files))

	var nextID int
	h.Db.Model(&core.BlogPost{}).Select("COALESCE(MAX(id), 0) + 1").Scan(&nextID)
	idString := strconv.Itoa(nextID)

	curPath, err := os.Getwd()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while getting the current work directory",
		})
		return
	}

	for idx, file := range files {
		extension := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%v%v", idx+1, extension)
		files[idx].Filename = filename

		processedFiles = append(processedFiles, core.AttachmentRecord{
			Filename:   filename,
			Path:       filepath.Join(curPath, "../backend/assets", idString, filename),
			BlogPostID: nextID,
		})
	}

	post := core.BlogPost{
		ID:          nextID,
		Created_At:  time.Now().Unix(),
		Edited_At:   time.Now().Unix(),
		Title:       title,
		Content:     content,
		Attachments: processedFiles,
	}

	if err := h.Db.Create(&post).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while creating blog post. Transaction cancelled",
		})
		return
	}

	for _, file := range files {
		filePath := path.Join(curPath, "assets", idString, file.Filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Error while downloading attached file",
			})
			return
		}

		// TODO: finish thumbnails
		// go func() {
		// 	err := core.MakeThumbnail(filePath, 100, 100)
		// 	if err != nil {
		// 		fmt.Printf("\nError while doing shit: %v", err)
		// 	}
		// }()

	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post added successfully",
	})

}
