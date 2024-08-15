package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateRequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Images  []any  `json:"images"`
}

type EditRequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Images  []any  `json:"images"`
}

func (h *Handler) CreateBlogPost(c *gin.Context) {
	var data CreateRequestBody
	db, err := core.GetDb()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong while contacting the database. Please report this to the website administrators",
		})
		return
	}

	// Malformed JSON data
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error() + ". Make sure the JSON data is valid and try again.",
		})
		return
	}

	if data.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "data field is empty",
		})
		return
	}

	if data.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "content field is empty",
		})
		return
	}

	var foundPost core.BlogPost

	db.Where("title = ?", data.Title).First(&foundPost)

	if foundPost != (core.BlogPost{}) {
		c.JSON(http.StatusConflict, gin.H{
			"error": "post with this title already exists",
		})
		return
	}

	post := core.BlogPost{
		Title:      data.Title,
		Content:    data.Content,
		Created_At: time.Now(),
		Edited_At:  time.Now(),
		Images:     "", //TODO
	}

	db.Create(&post)

	c.JSON(http.StatusOK, post)
}

func (h *Handler) EditBlogPost(c *gin.Context) {
	var data EditRequestBody
	db, err := core.GetDb()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong while contacting the database. Please report this to the website administrators",
		})
		return
	}

	// Malformed JSON data
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var record core.BlogPost
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to convert id string to int. Make sure it's correct and try again",
		})
		return
	}

	db.Where("id = ?", id).First(&record)

	// Record doesn't exist
	if record == (core.BlogPost{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "record not found in database",
		})
		return
	}

	var newBlogData core.BlogPost

	if data.Title != "" {
		newBlogData.Title = data.Title
	}

	if data.Content != "" {
		newBlogData.Content = data.Content
	}

	if len(data.Images) > 0 {
		newBlogData.Images = "" //TODO
	}

	db.Where("id = ?", id).UpdateColumns(newBlogData)

	c.JSON(http.StatusOK, record)
}

func (h *Handler) DeleteBlogPost(c *gin.Context) {
	db, err := core.GetDb()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong while contacting the database. Please report this to the website administrators",
		})
		return
	}

	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to convert id string to int. Make sure it's correct and try again",
		})
		return
	}

	var post core.BlogPost

	db.Where("id = ?", id).First(&post)

	if post == (core.BlogPost{}) {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	db.Delete(core.BlogPost{ID: id})

	c.JSON(http.StatusOK, nil)
}
