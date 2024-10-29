package srv

import (
	"bash06/strona-fundacja/src/backend/core"
	"fmt"
	"io"
	"net/http"
	"os"
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
			"message": errSqlQuery,
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
			"error":   "Group ID is not a valid integer",
			"message": nil,
		})
		return
	}

	var groupInfo *core.GalleryGroup

	if err := s.Db.Where("id = ?", groupId).Select("id", "name").First(&groupInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Group with this ID doesn't exist",
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

/*
Returns info about all images in a specified group by it's ID. If the group contains no images this will return an empty array. If the group doesn't exist this will return 404

Available query params:
NAME       DEFAULT_VAL       DESCRIPTION
limit      (3)               Sets the amount of returned images in one request
offset     (0)               Offsets the returned images by X

How offset works:
* offset = 0
[* * * * * *]
^^^^^^^^^^^^^ (All images will be returned (until we hit the limit). Nothing will be omitted)

* offset = 3
[* * * * * *]
-------^^^^^^ (Images from this point on will be returned. The first 3 will be ignored)

How limit works:
* limit = 2 (and offset 0)
[* * * * * *]
-^^^^^------- (Only these images will be returned. Everything else is omitted)

* limit = 3 and offset = 5
[* * * * * *]
-----------^ (This image will be returned BUT because there are not enough images the array will be of len(1))

Example request urls:
GET /gallery/groups/1/images
GET /gallery/groups/123/images?limit=3&offset=3
*/
func (s *Server) GalleryGetImagesOne(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"error":   "Group ID is not a valid integer",
			"message": nil,
		})
		return
	}

	var groupInfo *core.GalleryGroup

	if err := s.Db.Preload("Images").Where("id = ?", groupId).First(&groupInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
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

	c.JSON(http.StatusOK, groupInfo.Images)
}

/*
Creates a new gallery group. If a group with the specified name already exists this will return 409 Conflict.
*/
func (s *Server) GalleryCreateOne(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "No group name provided",
			"message": nil,
		})
		return
	}

	if err := s.Db.Create(&core.GalleryGroup{Name: name}).Error; err != nil {
		if err.Error() == "UNIQUE constraint failed: gallery_groups.name" {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"error":   "Group with this name exists already",
				"message": nil,
			})
			return
		}

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

	err = c.Request.ParseMultipartForm(8 << 20) // 8 MiB memory limit
	if err != nil {
		if err.Error() == "http: request body too large" {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":   "Request body too large",
				"message": nil,
			})
			return
		}

		s.Log.Error("Error while parsing multipart upload", zap.Error(err))

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

	tx := s.Db.Begin()

	var wg sync.WaitGroup
	errors := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)

		go func() {
			defer wg.Done()

			filename := core.RandomFilename(10)

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

			url, err := s.Ovh.AddObject(os.Getenv("AWS_GALLERY_BUCKET_NAME"), filename, buffer)
			if err != nil {
				errors <- fmt.Errorf("error while sending file %v to S3: %v", file.Filename, err)
				return
			}

			if err := tx.Create(&core.GalleryRecord{
				GroupID: groupId,
				AltText: "", // TODO
				Key:     filename,
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

	var resolveImage *core.GalleryRecord

	err = s.Db.Where("group_id = ? AND id = ?", groupId, imageId).First(&resolveImage).Error
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

	if err := tx.Delete(&resolveImage).Error; err != nil {
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

	err = s.Ovh.DeleteObjectsBulk(os.Getenv("AWS_GALLERY_BUCKET_NAME"), []string{resolveImage.Key})
	if err != nil {
		s.Log.Error("Deleting files from S3 failed", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while deleting images from S3",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   nil,
		"message": "Image deleted successfully",
	})
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

	// TODO: this should also delete images from S3
	if err := tx.Where("group_id = ?", groupId).Delete(&core.GalleryRecord{}).Error; err != nil {
		tx.Rollback()
		s.Log.Error("Failed to remove children of gallery group", zap.Int("Gallery Group ID", groupId), zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while deleting all related records",
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
		"message": "All images from group deleted successfully",
	})
}

func (s *Server) GalleryDelete(c *gin.Context) {
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

	var resolveGroup *core.GalleryGroup
	imageKeys := make([]string, len(resolveGroup.Images))

	tx := s.Db.Begin()

	if err := tx.Where("id = ?", groupId).First(&resolveGroup).Error; err != nil {
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

	for _, image := range resolveGroup.Images {
		imageKeys = append(imageKeys, image.Key)
	}

	err = s.Ovh.DeleteObjectsBulk(os.Getenv("AWS_GALLERY_BUCKET_NAME"), imageKeys)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Some files failed to be deleted from S3",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Group deleted successfully",
		"error":   nil,
	})
}
