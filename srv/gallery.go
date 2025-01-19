package srv

import (
	"bash06/tipitipi-backend/core"
	"bash06/tipitipi-backend/flags"
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	maxFileSize = *flags.MaxFileSize << 20
)

type GroupCreateBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GroupDeleteBody struct {
	IDs []int `json:"ids"`
}

/*
Returns the info and the files associated with a group
TODO: Add thumbnails
*/
func (s *Server) GalleryGetOne(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "id is not a valid integer",
		})
		return
	}

	var group *core.GalleryGroup

	err = s.Db.
		Model(&core.GalleryGroup{}).
		Preload("Images").
		Where("id = ?", id).
		First(&group).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Group not found",
			})
			return
		}

		s.Log.Error("Failed to find group: %v", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	c.JSON(http.StatusOK, group)
}

/*
Returns multiple gallery groups based on parameters
TODO: Add thumbnails
*/
func (s *Server) GalleryGetBulk(c *gin.Context) {
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil || offset < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "offset is invalid",
		})
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit > 10 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "limit is invalid",
		})
		return
	}

	var groups []core.GalleryGroup

	err = s.Db.
		Model(&core.GalleryGroup{}).
		Preload("Images").
		Offset(offset).
		Limit(limit).
		Find(&groups).
		Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	c.JSON(http.StatusOK, groups)
}

/*
Creates a new gallery group
*/
func (s *Server) GalleryCreateOne(c *gin.Context) {
	var group GroupCreateBody

	if err := c.BindJSON(&group); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to parse JSON",
		})
		return
	}

	if group.Name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Name is required",
		})
		return
	}

	// Fallback description if none is provided
	if group.Description == "" {
		group.Description = "Brak opisu"
	}

	err := s.Db.Create(&core.GalleryGroup{
		Name:        group.Name,
		Description: group.Description,
	}).Error
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: gallery_groups.name" {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"error": "Group with this name already exists",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	c.Status(http.StatusCreated)
}

/*
Posts images to a gallery group
*/
func (s *Server) GalleryPostBulk(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "id is not a valid integer",
		})
		return
	}

	var group int

	err = s.Db.
		Model(&core.GalleryGroup{}).
		Where("id = ?", id).
		Select("id").
		First(&group).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Group not found",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	// Set file size limit
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxFileSize*(*flags.MaxFileUploads))

	form, _ := c.MultipartForm()
	files := form.File["images[]"]

	var wg sync.WaitGroup
	downloadErrs := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if file.Size > maxFileSize {
				downloadErrs <- errors.New("file is too large")
				return
			}

			result, err := s.DownloadFile(file)
			if err != nil {
				downloadErrs <- err
				return
			}

			err = s.Db.Create(&core.GalleryRecord{
				Filename: result.Filename,
				Mimetype: result.Mimetype,
				Size:     result.Size,
				GroupID:  group,
			}).Error
			if err != nil {
				downloadErrs <- err
				return
			}
		}()
	}

	wg.Wait()
	close(downloadErrs)

	for err := range downloadErrs {
		s.Log.Error("Failed to download file: %v", zap.Error(err))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Failed to upload file",
			})
			return
		}
	}

	c.Status(http.StatusOK)
}

/*
Delete multiple images from a gallery group
*/
func (s *Server) GalleryDeleteMany(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "id is not a valid integer",
		})
		return
	}

	var group core.GalleryGroup

	err = s.Db.
		Where("id = ?", id).
		Preload("Images").
		First(&group).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Group not found",
			})
			return
		}

		s.Log.Error("Failed to find group: %v", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	var body GroupDeleteBody

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to parse JSON",
		})
		return
	}

	// Delete multiple records at once using a transaction
	tx := s.Db.Begin()

	for _, id := range body.IDs {
		tx.Delete(&core.GalleryRecord{}, id)
	}

	if err := tx.Commit().Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete images",
		})
		return
	}

	for _, image := range group.Images {
		for _, id := range body.IDs {
			if image.ID == id {
				err = s.DeleteFile(image.Filename)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error":   err.Error(),
						"message": "Failed to delete images",
					})
					return
				}
			}
		}
	}

	c.Status(http.StatusOK)
}

/*
Deletes a gallery group
*/
func (s *Server) GalleryDeleteGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "id is not a valid integer",
		})
		return
	}

	var group core.GalleryGroup

	err = s.Db.
		Where("id = ?", id).
		Preload("Images").
		First(&group).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Group not found",
			})
			return
		}

		s.Log.Error("Failed to find group: %v", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	err = s.Db.Delete(&group).Error
	if err != nil {
		s.Log.Error("Failed to delete group: %v", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	deleteErrs := make(chan error, len(group.Images))

	for _, image := range group.Images {
		err = s.Db.Delete(&image).Error
		if err != nil {
			deleteErrs <- err
			return
		}

		err = s.DeleteFile(image.Filename)
		if err != nil {
			deleteErrs <- err
			return
		}
	}

	close(deleteErrs)
	for err := range deleteErrs {
		if err != nil {
			s.Log.Error("Failed to delete group: %v", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Failed to delete group",
			})
			return
		}
	}
}
