package srv

import (
	"bash06/tipitipi-backend/core"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PageContentUpdateBody struct {
	Content string `json:"content"`
}

/*
Returns the contents of a page
*/
func (s *Server) PageGetOne(c *gin.Context) {
	name := c.Param("name")

	var page core.PageContent

	err := s.Db.
		Where("name = ?", name).
		Find(&page).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}

		s.Log.Error("Failed to get page content", zap.Error(err))

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get page content"})
		return
	}

	c.JSON(http.StatusOK, page)
}

/*
Updates the contents of a page
*/
func (s *Server) PageUpdateOne(c *gin.Context) {
	name := c.Param("name")

	var newContent PageContentUpdateBody

	err := c.BindJSON(&newContent)
	if err != nil {
		s.Log.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON"})
		return
	}

	var pageContent core.PageContent

	err = s.Db.
		Where("name = ?", name).
		Find(&pageContent).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}

		s.Log.Error("Failed to get page content", zap.Error(err))

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get page content"})
		return
	}

	pageContent.Content = newContent.Content
	// Set the name just in case
	pageContent.Name = name

	err = s.Db.
		Save(&pageContent).
		Error

	if err != nil {
		s.Log.Error("Failed to update page content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update page content"})
		return
	}

	c.JSON(http.StatusOK, pageContent)
}
