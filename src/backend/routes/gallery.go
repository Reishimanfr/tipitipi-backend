package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// WORKS
func (h *Handler) getInfoAllGroups(c *gin.Context) {
	var allGroups []*core.GalleryGroup

	if err := h.Db.Model(&core.GalleryGroup{}).Select("id", "name").Where("1 = 1").Find(&allGroups).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": errSqlQuery,
		})
		return
	}

	if len(allGroups) < 1 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   "No gallery groups exist",
			"message": "",
		})
		return
	}

	c.JSON(http.StatusOK, allGroups)
}

// WORKS
func (h *Handler) getInfoOnGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Group ID is not a valid integer",
			"message": nil,
		})
		return
	}

	var groupInfo *core.GalleryGroup

	if err := h.Db.Where("id = ?", groupId).Select("id", "name").First(&groupInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Group under this ID doesn't exist",
				"message": nil,
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing the request",
		})
		return
	}

	c.JSON(http.StatusOK, groupInfo)
}

// TODO
func (h *Handler) getImagesAllGroups(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "3"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Limit option is not a valid integer",
			"message": nil,
		})
		return
	}

	if limit < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Limit option must be at least 0",
			"message": nil,
		})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Offset option is not a valid integer",
			"message": nil,
		})
		return
	}

	if offset < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Offset option must be at least 0",
			"message": nil,
		})
		return
	}

	var galleryRecords []*core.GalleryRecord

	if err := h.Db.Model(&core.GalleryRecord{}).Offset(offset).Limit(limit).Scan(&galleryRecords).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Records not found",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while querying the database",
		})
		return
	}

	formatRecords := make(map[int][]*core.GalleryRecord)

	for _, record := range galleryRecords {
		formatRecords[record.GroupID] = append(formatRecords[record.GroupID], record)
	}

	c.JSON(http.StatusOK, formatRecords)
}

// WORKS
func (h *Handler) getImagesFromGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"error":   "Group ID is not a valid integer",
			"message": nil,
		})
		return
	}

	var groupImages []*core.GalleryRecord

	if err := h.Db.Where("group_id = ?", groupId).Find(&groupImages).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Something went wrong while querying the SQL database",
			})
			return
		}

		h.Log.Error("SQL query failed", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while querying the SQL database",
		})
		return
	}

	c.JSON(http.StatusOK, groupImages)
}

// WORKS
func (h *Handler) createGalleryGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Group ID is not a valid integer",
			"message": nil,
		})
		return
	}

	if groupId < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Group ID must not be smaller than 0",
			"message": nil,
		})
		return
	}

	name := c.Param("name")

	if name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "No group name provided",
			"message": nil,
		})
		return
	}

	var existsGroup bool
	h.Db.Model(&core.BlogPost{}).Select("count(*) > 0").Where("group_id = ? OR name = ?", groupId, name).Find(&existsGroup)

	if existsGroup {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error":   "Group with this ID or name already exist",
			"message": nil,
		})
		return
	}

	if err := h.Db.Create(&core.GalleryGroup{
		ID:   groupId,
		Name: name,
	}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": errSqlQuery,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error":   nil,
		"message": "Group created successfully",
	})
}

// WORKS
func (h *Handler) postImageToGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"error":   "Group ID is not a valid integer",
			"message": nil,
		})
		return
	}

	if err := h.Db.Model(&core.GalleryGroup{}).Where("id = ?", c.Param("groupId")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Group with this ID not found",
				"message": nil,
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing the request",
		})
		return
	}

	// Max 100MiB and 5 images per request
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 100<<20)

	err = c.Request.ParseMultipartForm(8 << 20) // 8 MiB memory limit
	if err != nil {
		if err.Error() == "http: request body too large" {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":   "Request body too large",
				"message": nil,
			})
			return
		}

		h.Log.Error("Error while parsing multipart upload", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": errMultipartParse,
		})
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["files[]"]

	if len(files) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "No files provided",
			"message": nil,
		})
		return
	}

	if len(files) > 5 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Too many files provided (max 5)",
			"message": nil,
		})
		return
	}

	tx := h.Db.Begin()

	var wg sync.WaitGroup
	errors := make(chan error, len(files))

	for idx, file := range files {
		wg.Add(1)

		go func() {
			defer wg.Done()

			ext := filepath.Ext(file.Filename)
			formatName := fmt.Sprintf("%v-%v%v", groupId, idx, ext)

			f, err := file.Open()
			if err != nil {
				errors <- fmt.Errorf("error while opening multipart file %v: %v", file.Filename, err)
				return
			}

			buffer, err := io.ReadAll(f)
			if err != nil {
				errors <- fmt.Errorf("error while reading file %v into memory: %v", file.Filename, err)
				return
			}

			url, err := h.Ovh.AddObject(os.Getenv("AWS_GALLERY_BUCKET_NAME"), formatName, buffer)
			if err != nil {
				errors <- fmt.Errorf("error while sending file %v to S3: %v", file.Filename, err)
				return
			}

			if err := tx.Create(&core.GalleryRecord{
				GroupID: groupId,
				AltText: "", // TODO
				URL:     *url,
			}).Error; err != nil {
				tx.Rollback()

				errors <- fmt.Errorf("error while adding database record for file %v: %v", file.Filename, err)
				return
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": errAwsUpload,
			})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while finishing database transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   nil,
		"message": "Files uploaded to gallery group successfully!",
	})
}

func (h *Handler) deleteImageFromGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"error":   "Group ID is not a valid integer",
			"message": nil,
		})
		return
	}

	if groupId < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Group ID must be at least 0",
			"message": nil,
		})
		return
	}

	imageId, err := strconv.Atoi(c.Param("imageId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"error":   "Group ID is not a valid integer",
			"message": nil,
		})
		return
	}

	if imageId < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Image ID must be at least 0",
			"message": nil,
		})
		return
	}

	var existsGroup, existsImage bool
	if err := h.Db.Model(&core.GalleryGroup{}).Select("count(*) > 0").Where("id = ?", groupId).Find(&existsGroup).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while querying the SQL database",
		})
		return
	}

	if !existsGroup {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   "Group with this ID doesn't exist",
			"message": nil,
		})
		return
	}

	if err := h.Db.Model(&core.GalleryRecord{}).Select("count(*) > 0").Where("id = ?", imageId).Find(&existsImage).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while querying the SQL database",
		})
		return
	}

	if !existsImage {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   "Image with this ID doesn't exist",
			"message": nil,
		})
		return
	}

	tx := h.Db.Begin()

	if err := tx.Model(&core.GalleryRecord{}).Delete("group_id = ? AND id = ?", groupId, imageId).Error; err != nil {
		tx.Rollback()

		h.Log.Error("SQL query failed", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while querying the SQL database",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": errSqlTransaction,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   nil,
		"message": "Image deleted successfully",
	})
}

func (h *Handler) deleteAllFromGroup(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Group ID is not a valid integer",
			"message": nil,
		})
		return
	}

	if groupId < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Group ID must be at least 0",
			"message": nil,
		})
		return
	}

	var existsGroup bool
	err = h.Db.Model(&core.GalleryGroup{}).Select("count(*) > 0").Where("id = ?", groupId).Find(&existsGroup).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while querying the SQL database",
		})
		return
	}

	// TODO: this should also delete images from S3
	if err := h.Db.Where("group_id = ?", groupId).Delete(&core.GalleryRecord{}).Error; err != nil {
		h.Log.Error("Failed to remove children of gallery group", zap.Int("Gallery Group ID", groupId), zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while deleting all related records",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   nil,
		"message": "All images from group deleted successfully",
	})
}

func (h *Handler) deleteGroup(c *gin.Context) {
}
