package srv

import (
	"bash06/strona-fundacja/src/backend/core"
	"fmt"
	"net/http"
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

var sortOptions = []string{"newest", "oldest"}

// blog/post/:id
// Returns a single post under some ID
func (s *Server) BlogGetOne(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID provided",
		})
		return
	}

	files := c.DefaultQuery("files", "false") == "true"

	var postRecord *core.BlogPost
	var dbErr error

	if files {
		dbErr = s.Db.Preload("Files").Where("id = ?", id).First(&postRecord).Error
	} else {
		dbErr = s.Db.Where("id = ?", id).First(&postRecord).Error
	}

	if dbErr != nil {
		if dbErr == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Record not found",
			})
			return
		}

		s.Log.Error("Error while searching for a post record", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   dbErr.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	c.JSON(http.StatusOK, postRecord)
}

// /blog/posts
// Returns multiple posts depending on the provided settings
func (s *Server) BlogGetBulk(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "5")
	sort := strings.ToLower(c.DefaultQuery("sort", "newest"))
	files := c.DefaultQuery("files", "false") == "true"
	partial := c.DefaultQuery("partial", "false") == "true"

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Offset is not a valid integer",
		})
		return
	}

	if offset < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Offset must be at least 0",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Limit is not a valid integer",
		})
		return
	}

	if limit < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Limit must be at least 1",
		})
		return
	}

	if !slices.Contains(sortOptions, sort) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid sort option provided",
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

	if partial {
		err = s.Db.Select([]string{"created_at", "title", "id"}).Order(orderClause).Offset(offset).Limit(limit).Find(&postRecords).Error
	} else if files {
		err = s.Db.Preload("Files").Order(orderClause).Offset(offset).Limit(limit).Find(&postRecords).Error
	} else {
		err = s.Db.Order(orderClause).Offset(offset).Limit(limit).Find(&postRecords).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Post not found",
			})
			return
		}

		s.Log.Error("Error while getting post records from database", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Something went wrong while processing your request",
		})
		return
	}

	c.JSON(http.StatusOK, postRecords)
}

// /blog/post
// Creates a single blog post
func (s *Server) BlogCreateOne(c *gin.Context) {
	err := c.Request.ParseMultipartForm(4 << 20)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	title := strings.Trim(c.PostForm("title"), "")
	content := strings.Trim(c.PostForm("content"), "")

	if title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No post title provided",
		})
		return
	}

	if content == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No post content provided",
		})
		return
	}

	var exists bool
	s.Db.Model(&core.BlogPost{}).Select("count(*) > 0").Where("title = ?", title).Find(&exists)

	if exists {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error": "Post with this title already exists",
		})
		return
	}

	form := c.Request.MultipartForm
	files := form.File["files[]"]

	s.Db.Create(&core.BlogPost{
		Created_At: time.Now().Unix(),
		Edited_At:  time.Now().Unix(),
		Title:      title,
		Content:    content,
	})

	// Exit early if there are no files to be uploaded
	if len(files) < 1 {
		c.Status(http.StatusOK)
		return
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(files))

	var nextPostId int
	s.Db.Model(&core.BlogPost{}).Select("COALESCE(MAX(id), 0)").Scan(&nextPostId)

	for _, file := range files {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if result, err := s.DownloadFile(file); err != nil {
				s.Log.Error("Failed to download file", zap.Error(err))
				errors <- err
			} else {
				if err := s.Db.Create(&core.File{
					Filename:   result.Filename,
					BlogPostID: nextPostId,
					Size:       result.Size,
					Mimetype:   result.Mimetype,
				}).Error; err != nil {
					errors <- fmt.Errorf("failed to create database record for %s: %v", result.Filename, err)
				}
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Some files failed ot be downloaded",
			})
			return
		}
	}

	c.Status(http.StatusOK)
}

// blog/post/:id
// Deletes a post under some ID (and it's related attachments if any)
func (s *Server) BlogDeleteOne(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID provided",
		})
		return
	}

	post := new(core.BlogPost)

	if err := s.Db.Preload("Files").Where("id = ?", id).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Post not found",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong while processing your request",
			"error":   err.Error(),
		})
		return
	}

	tx := s.Db.Begin()

	if err := tx.Model(&core.BlogPost{}).Delete("id = ?", id).Error; err != nil {
		tx.Rollback()

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	if err := tx.Model(&core.File{}).Delete("blog_post_id = ?", post.ID).Error; err != nil {
		tx.Rollback()

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

	if len(post.Files) < 1 {
		c.Status(http.StatusOK)
		return
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(post.Files))

	for _, file := range post.Files {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := s.DeleteFile(file.Filename); err != nil {
				s.Log.Error("Failed to delete file from disk", zap.String("Filename", file.Filename), zap.Error(err))
				errors <- err
				return
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"errors":  errors,
				"message": "Some files failed to be deleted from disk",
			})
			return
		}
	}

	c.Status(http.StatusOK)
}

// blog/post/:id
// Edits a post under the specified ID
func (s *Server) BlogEditOne(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID provided",
		})
		return
	}

	var post *core.BlogPost

	if err := s.Db.Where("id = ?", id).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Post not found",
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

	err = c.Request.ParseMultipartForm(4 << 20)
	if err != nil {
		s.Log.Error("Failed to parse multipart form", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	title := strings.Trim(c.PostForm("title"), "")
	content := strings.Trim(c.PostForm("content"), "")

	if title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No post title provided",
		})
		return
	}

	if content == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No post content provided",
		})
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["files[]"]

	var wg sync.WaitGroup
	errors := make(chan error, len(post.Files))

	for _, file := range post.Files {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := s.DeleteFile(file.Filename); err != nil {
				s.Log.Error("Failed to delete file from disk", zap.String("Filename", file.Filename), zap.Error(err))
				errors <- err
				return
			}

			if err := s.Db.Model(&core.File{}).Delete("blog_post_id = ?", id).Error; err != nil {
				s.Log.Error("Failed to delete file record from database", zap.String("Filename", file.Filename), zap.Error(err))
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"errors":  errors,
				"message": "Something went wrong while processing your request",
			})
			return
		}
	}

	if err := s.Db.UpdateColumns(&core.BlogPost{
		ID:        id,
		Title:     title,
		Content:   content,
		Edited_At: time.Now().Unix(),
	}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong while processing your request",
		})
		return
	}

	if len(files) < 1 {
		c.Status(http.StatusOK)
		return
	}

	var nextPostId int
	s.Db.Model(&core.BlogPost{}).Select("COALESCE(MAX(id), 0)").Scan(&nextPostId)

	errors = make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)

		go func() {
			defer wg.Done()

			result, err := s.DownloadFile(file)
			if err != nil {
				errors <- err
				return
			}

			if err := s.Db.Create(&core.File{
				Filename:   result.Filename,
				Size:       result.Size,
				Mimetype:   result.Mimetype,
				BlogPostID: nextPostId,
			}).Error; err != nil {
				errors <- err
				return
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"errors":  errors,
				"message": "Something went wrong while processing your request",
			})
			return
		}
	}
}
