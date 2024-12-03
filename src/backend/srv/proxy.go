package srv

import (
	"bash06/strona-fundacja/src/backend/flags"
	"bufio"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	errOptKeyNil       = "No image key provided"
	errOptBucketInv    = "Invalid bucket type provided"
	errFileReadFailure = "Error while reading file"
)

type PartialFileInfo struct {
	Size     int64
	Mimetype string
}

func (s *Server) Proxy(c *gin.Context) {
	key := c.Query("key")

	if key == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errOptKeyNil,
			"message": nil,
		})
		return
	}

	var file *PartialFileInfo

	if err := s.Db.Model(&PartialFileInfo{}).Where("filename LIKE ?", "%"+key+"%").Select("size", "mimetype").First(&file).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to find file",
		})
		return
	}

	filePath := filepath.Join(flags.BasePath, key)

	srcFile, err := os.Open(filePath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to open file",
		})
		return
	}

	defer srcFile.Close()

	reader := bufio.NewReader(srcFile)
	buffer := make([]byte, 1024)

	reader.Read(buffer)

	c.Header("Content-Type", file.Mimetype)
	c.Header("Content-Length", strconv.Itoa(int(file.Size)))
	c.Header("Cache-Control", "max-age=3600")

	c.Data(http.StatusOK, file.Mimetype, buffer)
}
