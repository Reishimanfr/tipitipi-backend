package srv

import (
	"bash06/tipitipi-backend/core"
	"bash06/tipitipi-backend/flags"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	errOptKeyNil       = "No image key provided"
	errOptTypeNil      = "No type provided"
	errOptBucketInv    = "Invalid bucket type provided"
	errFileReadFailure = "Error while reading file"
)

type PartialFileInfo struct {
	Size     int64
	Mimetype string
	Filename string
}

func (s *Server) Proxy(c *gin.Context) {
	key := c.Query("key")
	_type := c.Query("type")

	if key == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errOptKeyNil,
			"message": nil,
		})
		return
	}

	if _type == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errOptTypeNil,
			"message": nil,
		})
		return
	}

	var file *PartialFileInfo

	switch _type {
	case "blog":
		{
			if err := s.Db.Model(&core.File{}).Where("filename LIKE ?", key+"%").Select("size", "mimetype", "filename").First(&file).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":   err.Error(),
					"message": "Failed to find file",
				})
				return
			}
		}

	case "gallery":
		{
			if err := s.Db.Model(&core.GalleryRecord{}).Where("filename LIKE ?", key+"%").Select("size", "mimetype", "filename").First(&file).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":   err.Error(),
					"message": "Failed to find file",
				})
				return
			}
		}
	}

	if file == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Image not found",
		})
		return
	}

	filePath := filepath.Join(flags.BasePath, "files", file.Filename)

	srcFile, err := os.Open(filePath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to open file",
		})
		return
	}

	defer srcFile.Close()

	// buffer := make([]byte, 0)

	// reader.ReadBytes(buffer)

	c.Header("Content-Type", file.Mimetype)
	c.Header("Content-Length", strconv.Itoa(int(file.Size)))
	c.Header("Cache-Control", "max-age=3600")

	http.ServeContent(c.Writer, c.Request, file.Filename, time.Now(), srcFile)

	// c.Data(http.StatusOK, file.Mimetype, buffer)
}
