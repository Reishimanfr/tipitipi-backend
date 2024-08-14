package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"

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

type DeleteRequestBody struct {
	Id int `json:"id"`
}

func (h *Handler) CreateBlogPost(c *gin.Context) {
	var data CreateRequestBody
	db, err := core.GetDb()

	if err != nil || db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong while contacting the database. Please report this to the website administrators",
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

	if data.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data field is empty",
		})
		return
	}

	if data.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Content field is empty",
		})
		return
	}
}

func (h *Handler) EditBlogPost(c *gin.Context) {

}

func (h *Handler) DeleteBlogPost(c *gin.Context) {

}
