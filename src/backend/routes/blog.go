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
	sortOptions = []string{"newest", "oldest", "likes"}
)

const (
	err_post_not_found         = "Post not found"
	err_post_id_invalid        = "Invalid post ID provided"
	err_post_exists            = "Post with this title already exists"
	err_opt_offset_invalid     = "Offset value is not a valid integer"
	err_opt_offset_too_small   = "Offset value must be at least 0"
	err_opt_limit_invalid      = "Limit is not a valid integer"
	err_opt_limit_too_small    = "Limit must be at least 1"
	err_opt_sort_invalid       = "Invalid sort option provided"
	err_sql_query              = "Error while executing SQL query"
	err_sql_transaction_failed = "Failed to commit SQL transaction"
	err_getwd_failed           = "Failed to get the current working directory"
	err_multipart_parse        = "Failed to parse multipart form"
	err_multipart_no_title     = "No post title provided"
	err_multipart_no_content   = "No post content provided"
	err_attach_delete          = "Failed to delete files"
	err_aws_upload_failed      = "Some files failed to upload to S3"
	ok_post_and_attach_delete  = "Post and it's files deleted successfully"
	ok_post_create_success     = "Post created successfully"
)

func getMimeType(fileName string) string {
	ext := filepath.Ext(fileName)

	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	case ".mp4":
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}

// blog/post/:id
// Returns a single post under some ID
func (h *Handler) getOne(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)
	atts := c.Query("attachments") == "true"

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err_post_id_invalid,
			"message": nil,
		})
		return
	}

	postRecord := new(core.BlogPost)
	result := new(gorm.DB)

	if atts {
		result = h.Db.Preload("Attachments").Where("id = ?", id).First(&postRecord)
	} else {
		result = h.Db.Where("id = ?", id).First(&postRecord)
	}

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		h.Log.Error("Error while searching for a post record", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": err_sql_query,
		})
		return
	}

	if postRecord == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   err_post_not_found,
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
	atts := c.DefaultQuery("attachments", "false") == "true"
	partial := c.DefaultQuery("partial", "false") == "true"

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err_opt_offset_invalid,
			"message": nil,
		})
		return
	}

	if offset < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err_opt_offset_too_small,
			"message": nil,
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err_opt_limit_invalid,
			"message": nil,
		})
		return
	}

	if limit < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err_opt_limit_too_small,
			"message": nil,
		})
		return
	}

	if !slices.Contains(sortOptions, sort) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err_opt_sort_invalid,
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
	} else if atts {
		result = h.Db.Preload("Attachments").Order(orderClause).Offset(offset).Limit(limit).Find(&postRecords)
	} else {
		result = h.Db.Order(orderClause).Offset(offset).Limit(limit).Find(&postRecords)
	}

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		h.Log.Error("Error while getting post records from database", zap.Error(result.Error))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": err_sql_query,
		})
		return
	}

	if len(postRecords) < 1 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   err_post_not_found,
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
			"message": err_multipart_parse,
		})
		return
	}

	title := strings.Trim(c.PostForm("title"), "")
	content := strings.Trim(c.PostForm("content"), "")

	if title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err_multipart_no_title,
			"message": nil,
		})
		return
	}

	if content == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err_multipart_no_content,
			"message": nil,
		})
		return
	}

	var exists bool
	h.Db.Model(&core.BlogPost{}).Select("count(*) > 0").Where("title = ?", title).Find(&exists)

	if exists {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error":   err_post_exists,
			"message": nil,
		})
		return
	}

	var nextID int
	h.Db.Model(&core.BlogPost{}).Select("COALESCE(MAX(id), 0) + 1").Scan(&nextID)

	form := c.Request.MultipartForm
	files := form.File["files[]"]

	// todo: rewrite this to use multipart uploads within AWS' library
	if len(files) > 0 {
		var wg sync.WaitGroup
		errors := make(chan error, len(files))

		for i, fileHeader := range files {
			wg.Add(1)

			go func(fileHeader *multipart.FileHeader, nextId int) {
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

				ext := filepath.Ext(fileHeader.Filename)
				altFilename := fmt.Sprintf("%v-%v%v", nextID, i, ext)
				mime := getMimeType(fileHeader.Filename)

				url, err := h.Ovh.AddObject(buffer, os.Getenv("AWS_BLOG_BUCKET_NAME"), altFilename, mime)
				if err != nil {
					errors <- fmt.Errorf("failed to upload file %s: %v", fileHeader.Filename, err)
				}

				h.Db.Create(&core.AttachmentRecord{
					BlogPostID: nextID,
					URL:        url,
					Filename:   altFilename,
				})
			}(fileHeader, nextID)
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
	}

	h.Db.Create(&core.BlogPost{
		Created_At: time.Now().Unix(),
		Edited_At:  time.Now().Unix(),
		Title:      title,
		Content:    content,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": ok_post_create_success,
	})
}

// blog/post/:id
// Deletes a post under some ID (and it's related attachments if any)
func (h *Handler) deleteOne(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err_post_id_invalid,
			"message": nil,
		})
		return
	}

	post := new(core.BlogPost)

	if err := h.Db.Preload("Attachments").Where("id = ?", id).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   err_post_not_found,
				"message": nil,
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err_sql_query,
			"error":   err.Error(),
		})
		return
	}

	tx := h.Db.Begin()

	if err := tx.Model(&core.BlogPost{}).Delete("id = ?", id).Error; err != nil {
		tx.Rollback()

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": err_sql_query,
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": err_sql_transaction_failed,
		})
		return
	}

	if len(post.Attachments) > 0 {
		bucketKeys := make([]string, len(post.Attachments))

		for _, at := range post.Attachments {
			bucketKeys = append(bucketKeys, at.Filename)
		}

		fmt.Println(bucketKeys)

		err := h.Ovh.DeleteObjectsBulk(os.Getenv("AWS_BLOG_BUCKET_NAME"), bucketKeys)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": err_attach_delete,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": ok_post_and_attach_delete,
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