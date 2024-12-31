package srv

import (
	"bash06/tipitipi-backend/core"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
Returns the ID and name of all existing groups. If no groups exists this will return an empty array

Example request urls:
GET /gallery/groups/all/info
*/
func (s *Server) GalleryGetGroupsBulk(c *gin.Context) {
	var allGroups []*core.GalleryGroup

	if err := s.Db.Model(&core.GalleryGroup{}).Select("id", "name").Where("1 = 1").Find(&allGroups).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	c.JSON(http.StatusOK, allGroups)
}

/*
Returns info on a specified group based on it's ID (as a url param). If the group doesn't exist this will return 404

Example request urls:
GET /gallery/groups/1/info
GET /gallery/groups/123/info
*/
func (s *Server) GalleryGetGroupOne(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Group ID is not a valid integer",
		})
		return
	}

	var group *core.GalleryGroup

	if err := s.Db.Where("id = ?", groupId).Select("id", "name").First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Group with this ID doesn't exist",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing the request",
		})
		return
	}

	c.JSON(http.StatusOK, group)
}

/*
Returns all images in a specified group by it's ID. If the group contains no images this will return an empty array. If the group doesn't exist this will return 404
Example request urls:
GET /gallery/groups/1/images
GET /gallery/groups/123/images?limit=3&offset=3
*/
func (s *Server) GalleryGetImagesOne(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"error": "Group ID is not a valid integer",
		})
		return
	}

	var group *core.GalleryGroup

	if err := s.Db.Preload("Images").Where("id = ?", groupId).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Group with this ID doesn't exist",
			})
			return
		}

		s.Log.Error("SQL query failed", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while querying the SQL database",
		})
		return
	}

	c.JSON(http.StatusOK, group.Images)
}

func (s *Server) GetEverything(c *gin.Context) {
	var everything []*core.GalleryGroup

	if err := s.Db.Preload("Images").Where("1 = 1").Find(&everything).Error; err != nil {
		s.Log.Error("Failed to get every image entry from the database", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong while processing your request",
		})
		return
	}

	c.JSON(http.StatusOK, everything)
}

func (s *Server) GalleryGetImagesBulk(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Offset is not a valid integer",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "1"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Limit is not a valid integer",
		})
		return
	}

	if offset < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Offset must be higher than 0",
		})
		return
	}

	if limit < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Limit must be higher than 0",
		})
		return
	}

	var groups []*core.GalleryGroup

	if err := s.Db.Preload("Images").Where("1 = 1").Offset(offset).Limit(limit).Find(&groups).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Group with this ID doesn't exist",
				"message": nil,
			})
			return
		}

		s.Log.Error("SQL query failed", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while querying the SQL database",
		})
		return
	}

	c.JSON(http.StatusOK, groups)
}

/*
Creates a new gallery group. If a group with the specified name already exists this will return 409 Conflict.
*/
func (s *Server) GalleryCreateOne(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No group name provided",
		})
		return
	}

	if err := s.Db.Create(&core.GalleryGroup{Name: name}).Error; err != nil {
		if err.Error() == "UNIQUE constraint failed: gallery_groups.name" {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"error": "Group with this name exists already",
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
Posts an image to an existing gallery group by it's ID. If the specified group doesn't exist this will return 404.
Requests should include an image (or multiple images) to be uploaded specified with the "files[]" key in a multipart form.

Example request:
POST /gallery/groups/1/images (multipart: files[] -> file1.png)
*/
func (s *Server) GalleryPostBulk(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"error":   "Group ID is not a valid integer",
			"message": nil,
		})
		return
	}

	var resolveGroup *core.GalleryGroup

	if err := s.Db.Where("id = ?", groupId).First(&resolveGroup).Error; err != nil {
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

	err = c.Request.ParseMultipartForm(8 << 20)
	if err != nil {
		if err.Error() == "http: request body too large" {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Request body too large",
			})
			return
		}

		s.Log.Error("Error while parsing multipart upload", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["files[]"]

	if len(files) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No files provided",
		})
		return
	}

	if len(files) > 5 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Too many files provided (max 5)",
		})
		return
	}

	tx := s.Db.Begin()

	var wg sync.WaitGroup
	errors := make(chan error, len(files))

	for _, f := range files {
		wg.Add(1)

		go func() {
			defer wg.Done()

			r, err := s.DownloadFile(f)
			if err != nil {
				errors <- err
				return
			}

			if err := s.Db.Create(&core.GalleryRecord{
				GroupID:  groupId,
				Filename: r.Filename,
				Mimetype: r.Mimetype,
				Size:     r.Size,
			}).Error; err != nil {
				errors <- fmt.Errorf("failed to create database record for %s: %v", f.Filename, err)
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Some files failed to be downloaded",
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

	c.Status(http.StatusOK)
}

/*
Deletes a specified image from a specified group. If either the group or image doesn't exist this will return 404.

Example request url:
DELETE /gallery/groups/1/images/1
*/
func (s *Server) GalleryDeleteOne(c *gin.Context) {
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

	var image *core.GalleryRecord

	err = s.Db.Where("group_id = ? AND id = ?", groupId, imageId).First(&image).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Record not found",
				"message": nil,
			})
			return
		}

		s.Log.Error("SQL query failed", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while querying the SQL database",
		})
		return
	}

	tx := s.Db.Begin()

	if err := tx.Delete(&image).Error; err != nil {
		tx.Rollback()

		s.Log.Error("Failed to delete gallery image record", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while deleting records from SQL database",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Error("SQL transaction commit failed", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while committing SQL transaction",
		})
		return
	}

	if err := s.DeleteFile(image.Filename); err != nil {
		s.Log.Error("Failed to delete file from disk", zap.String("Filename", image.Filename), zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete file from disk",
		})
		return
	}

	c.Status(http.StatusOK)
}

/*
Deletes all images from a group without actually deleting the group (basically a group purge). If the specified group doesn't exist this will return 404. If a group is empty this will return 200 Ok as if everything was deleted successfully.

Example request url: DELETE /gallery/groups/1
*/
func (s *Server) GalleryDeleteAll(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Group ID is not a valid integer",
			"message": nil,
		})
		return
	}

	var resolveGroup *core.GalleryGroup

	tx := s.Db.Begin()

	if err := tx.Where("id = ?", groupId).First(&resolveGroup).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Group with this ID doesn't exist",
				"message": nil,
			})
			return
		}

		s.Log.Error("Failed to look for gallery group", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while querying the SQL database",
		})
		return
	}

	if err := tx.Where("group_id = ?", groupId).Delete(&core.GalleryRecord{}).Error; err != nil {
		tx.Rollback()
		s.Log.Error("Failed to remove children of gallery group", zap.Int("GalleryGroupID", groupId), zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) GalleryDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Group ID is not a valid integer",
			"message": nil,
		})
		return
	}

	if id < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Group ID must be at least 0",
			"message": nil,
		})
		return
	}

	var resolveGroup *core.GalleryGroup

	tx := s.Db.Begin()

	if err := tx.Where("id = ?", id).First(&resolveGroup).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Group not found",
				"message": nil,
			})
			return
		}

		s.Log.Error("Failed to find gallery group", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while querying the database",
		})
		return
	}

	if err := tx.Select(clause.Associations).Delete(&resolveGroup).Error; err != nil {
		tx.Rollback()
		s.Log.Error("Failed to cascade delete group", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while deleting the group",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Error("SQL transaction failed", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while committing SQL transaction",
		})
		return
	}

	for _, img := range resolveGroup.Images {
		if err := s.DeleteFile(img.Filename); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Failed to delete file from disk ",
			})
			return
		}
	}

	c.Status(http.StatusOK)
}
