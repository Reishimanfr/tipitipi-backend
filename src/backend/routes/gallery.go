package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	err_multipart_no_images = "No images provided"
	ok_image_upload_success = "Gallery images uploaded successfully"
)

func (h *Handler) uploadToGallery(c *gin.Context) {
	err := c.Request.ParseMultipartForm(50 << 20)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": err_multipart_parse,
		})
		return
	}

	form := c.Request.MultipartForm
	images := form.File["images[]"]

	if len(images) < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err_multipart_no_images,
			"message": nil,
		})
		return
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(images))

	for _, img := range images {
		wg.Add(1)

		go func(fileHeader *multipart.FileHeader) {
			defer wg.Done()

			file, err := fileHeader.Open()
			if err != nil {
				errors <- fmt.Errorf("failed to open file %s: %v", fileHeader.Filename, err)
				return
			}

			defer file.Close()

			buffer := new(bytes.Buffer)

			if _, err := io.Copy(buffer, file); err != nil {
				errors <- fmt.Errorf("failed to read file %s: %v", fileHeader.Filename, err)
				return
			}

			filename := core.RandomFilename(10)
			mime := core.GetMIMEType(fileHeader.Filename)

			if !strings.HasPrefix(mime, "image/") {
				errors <- fmt.Errorf("file %s is not an image", fileHeader.Filename)
				return
			}

			optimizeBuf, err := core.OptimizeImage(buffer.Bytes(), 70)
			if err != nil {
				errors <- fmt.Errorf("failed to optimize file %s: %v", fileHeader.Filename, err)
				return
			}

			url, err := h.Ovh.AddObject(optimizeBuf, os.Getenv("AWS_GALLERY_BUCKET_NAME"), filename, mime)
			if err != nil {
				errors <- fmt.Errorf("failed to upload file %s: %v", fileHeader.Filename, err)
				return
			}

			// TODO: Implement alt text
			h.Db.Create(&core.GalleryRecord{
				URL:     url,
				AltText: "",
			})
		}(img)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": err_aws_upload_failed,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": ok_image_upload_success,
	})
}

func (h *Handler) deleteFromGallery(c *gin.Context) {

}

func (h *Handler) getGallery(c *gin.Context) {
	galleryRecords := []*core.GalleryRecord{}

	if err := h.Db.Model(&core.GalleryRecord{}).Where("1 = 1").Scan(&galleryRecords).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": err_sql_query,
		})
		return
	}

	c.JSON(http.StatusOK, galleryRecords)
}
