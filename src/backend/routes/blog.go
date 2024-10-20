package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	sortOptions       = []string{"newest", "oldest", "likes"}
	awsBlogBucketName = os.Getenv("AWS_BLOG_BUCKET_NAME")
)

const (
	errPostNotFound       = "Post not found"
	errPostIdInvalid      = "Invalid post ID provided"
	errPostDupName        = "Post with this title already exists"
	errOptOffsetInv       = "Offset value is not a valid integer"
	errOptOffsetSmall     = "Offset value must be at least 0"
	errOptLimitInv        = "Limit is not a valid integer"
	errOptLimitSmall      = "Limit must be at least 1"
	errOptSortInv         = "Invalid sort option provided"
	errSqlQuery           = "Error while executing SQL query"
	errSqlTransaction     = "Failed to commit SQL transaction"
	errGetwd              = "Failed to get the current working directory"
	errMultipartParse     = "Failed to parse multipart form"
	errMultipartNoTitle   = "No post title provided"
	errMultipartNoContent = "No post content provided"
	errAttachmentDelete   = "Failed to delete files"
	errAwsUpload          = "Some files failed to upload to S3"
	postDeletedOk         = "Post and it's files deleted successfully"
	postCreatedOk         = "Post created successfully"
)

// blog/post/:id
// Returns a single post under some ID
func (h *Handler) getOne(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)
	attach := c.Query("attachments") == "true"

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errPostIdInvalid,
			"message": nil,
		})
		return
	}

	postRecord := new(core.BlogPost)
	result := new(gorm.DB)

	if attach {
		result = h.Db.Preload("Attachments").Where("id = ?", id).First(&postRecord)
	} else {
		result = h.Db.Where("id = ?", id).First(&postRecord)
	}

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		h.Log.Error("Error while searching for a post record", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": errSqlQuery,
		})
		return
	}

	if postRecord == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   errPostNotFound,
			"message": nil,
		})
		return
	}

	c.JSON(http.StatusOK, postRecord)
}

// /blog/posts
// Returns multiple posts depending on the provided settings
func (h *Handler) getMany(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "5")
	sort := strings.ToLower(c.DefaultQuery("sort", "newest"))
	attach := c.DefaultQuery("attachments", "false") == "true"
	partial := c.DefaultQuery("partial", "false") == "true"

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errOptOffsetInv,
			"message": nil,
		})
		return
	}

	if offset < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errOptOffsetSmall,
			"message": nil,
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errOptLimitInv,
			"message": nil,
		})
		return
	}

	if limit < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errOptLimitSmall,
			"message": nil,
		})
		return
	}

	if !slices.Contains(sortOptions, sort) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errOptSortInv,
			"message": nil,
		})
		return
	}

	postRecords := make([]*core.BlogPost, limit)

	orderClause := clause.OrderByColumn{
		Desc: true,
		Column: clause.Column{
			Name: "created_at",
		},
	}

	if sort == "oldest" {
		orderClause = clause.OrderByColumn{
			Desc: false,
			Column: clause.Column{
				Name: "created_at",
			},
		}
	}

	result := new(gorm.DB)

	if partial {
		result = h.Db.Select([]string{"created_at", "title", "id"}).Order(orderClause).Offset(offset).Limit(limit).Find(&postRecords)
	} else if attach {
		result = h.Db.Preload("Attachments").Order(orderClause).Offset(offset).Limit(limit).Find(&postRecords)
	} else {
		result = h.Db.Order(orderClause).Offset(offset).Limit(limit).Find(&postRecords)
	}

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		h.Log.Error("Error while getting post records from database", zap.Error(result.Error))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": errSqlQuery,
		})
		return
	}

	if len(postRecords) < 1 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   errPostNotFound,
			"message": nil,
		})
		return
	}

	c.JSON(http.StatusOK, postRecords)
}

// /blog/post
// Creates a single blog post
func (h *Handler) createOne(c *gin.Context) {
	err := c.Request.ParseMultipartForm(50 << 20)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": errMultipartParse,
		})
		return
	}

	title := strings.Trim(c.PostForm("title"), "")
	content := strings.Trim(c.PostForm("content"), "")

	if title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errMultipartNoTitle,
			"message": nil,
		})
		return
	}

	if content == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errMultipartNoContent,
			"message": nil,
		})
		return
	}

	var exists bool
	h.Db.Model(&core.BlogPost{}).Select("count(*) > 0").Where("title = ?", title).Find(&exists)

	if exists {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error":   errPostDupName,
			"message": nil,
		})
		return
	}

	form := c.Request.MultipartForm
	files := form.File["files[]"]

	h.Db.Create(&core.BlogPost{
		Created_At: time.Now().Unix(),
		Edited_At:  time.Now().Unix(),
		Title:      title,
		Content:    content,
	})

	// Exit early if there are no files to be uploaded
	if len(files) < 1 {
		c.JSON(http.StatusOK, gin.H{
			"message": postCreatedOk,
		})
		return
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(files))

	var nextPostId int
	h.Db.Model(&core.BlogPost{}).Select("COALESCE(MAX(id), 0)").Scan(&nextPostId)

	for i, fHeader := range files {
		wg.Add(1)

		go func(fHeader *multipart.FileHeader, i int) {
			defer wg.Done()

			f, err := fHeader.Open()
			if err != nil {
				errors <- err
				return
			}

			defer f.Close()

			buffer := new(bytes.Buffer)

			if _, err := io.Copy(buffer, f); err != nil {
				errors <- err
				return
			}

			ext := filepath.Ext(fHeader.Filename)
			key := fmt.Sprintf("%v-%v%v", nextPostId, i, ext)

			// TODO: implement optimizing attachments based on the mimetype
			url, err := h.Ovh.AddObject(awsBlogBucketName, key, buffer.Bytes())
			if err != nil {
				errors <- err
			}

			h.Db.Create(&core.AttachmentRecord{
				BlogPostID: nextPostId,
				URL:        *url,
				Filename:   key,
			})

		}(fHeader, i)
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
}

// blog/post/:id
// Deletes a post under some ID (and it's related attachments if any)
func (h *Handler) deleteOne(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   errPostIdInvalid,
			"message": nil,
		})
		return
	}

	post := new(core.BlogPost)

	if err := h.Db.Preload("Attachments").Where("id = ?", id).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   errPostNotFound,
				"message": nil,
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": errSqlQuery,
			"error":   err.Error(),
		})
		return
	}

	tx := h.Db.Begin()

	if err := tx.Model(&core.BlogPost{}).Delete("id = ?", id).Error; err != nil {
		tx.Rollback()

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": errSqlQuery,
		})
		return
	}

	if err := tx.Model(&core.AttachmentRecord{}).Delete("blog_post_id = ?", post.ID).Error; err != nil {
		tx.Rollback()

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": errSqlQuery,
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

	if len(post.Attachments) > 0 {
		bucketKeys := []string{}

		for _, at := range post.Attachments {
			bucketKeys = append(bucketKeys, at.Filename)
		}

		err := h.Ovh.DeleteObjectsBulk(os.Getenv("AWS_BLOG_BUCKET_NAME"), bucketKeys)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": errAttachmentDelete,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": postDeletedOk,
		"error":   nil,
	})
}

// blog/post/:id
// Edits a post under the specified ID
func (h *Handler) editOne(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotImplemented, gin.H{
		"error": "Endpoint not implemented yet!",
	})
	// stringId := c.Param("id")
	// id, err := strconv.Atoi(stringId)

	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	// 		"error":   err_post_id_invalid,
	// 		"message": nil,
	// 	})
	// 	return
	// }

	// post := new(core.BlogPost)

	// if err := h.Db.Where("id = ?", id).First(&post).Error; err != nil {
	// 	if err == gorm.ErrRecordNotFound {
	// 		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
	// 			"error":   err_post_not_found,
	// 			"message": nil,
	// 		})
	// 		return
	// 	}

	// 	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
	// 		"error":   err_post_not_found,
	// 		"message": nil,
	// 	})
	// 	return
	// }

	// err = c.Request.ParseMultipartForm(50 << 20)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	// 		"error":   err.Error(),
	// 		"message": err_multipart_parse,
	// 	})
	// 	return
	// }

	// title := strings.Trim(c.PostForm("title"), "")
	// content := strings.Trim(c.PostForm("content"), "")

	// if title == "" {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	// 		"error":   err_multipart_no_title,
	// 		"message": nil,
	// 	})
	// 	return
	// }

	// if content == "" {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	// 		"error":   err_multipart_no_content,
	// 		"message": nil,
	// 	})
	// 	return
	// }

	// form := c.Request.MultipartForm
	// files := form.File["files[]"]

	// if len(files) > 0 {

	// }

	// if len(rawFiles) > 0 {
	// 	files := make([]*core.AttachmentRecord, 0, len(rawFiles))

	// 	currentDir, err := os.Getwd()
	// 	if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	// 			"error":   err.Error(),
	// 			"message": err_getwd_failed,
	// 		})
	// 		return
	// 	}

	// 	for range rawFiles {
	// 		name := core.RandStringBytesMaskImprSrcUnsafe(15)

	// 		files = append(files, &core.AttachmentRecord{
	// 			Filename:   name,
	// 			Path:       filepath.Join(currentDir, "../assets", stringId, name), //TODO: validate this code
	// 			BlogPostID: id,
	// 		})
	// 	}
	// }

	// someImagePath := post.Attachments[0].Path
	// dirPath := filepath.Join(someImagePath, "..")

	// fmt.Println(dirPath)

	// for _, image := range post.Attachments {
	// 	if err := os(image.Path); err != nil {
	// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	// 			"error":   err.Error(),
	// 			"message": err_attach_delete,
	// 		})
	// 		return
	// 	}
	// }
}
