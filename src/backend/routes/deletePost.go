package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) delete(c *gin.Context) {
	stringId := c.Param("id")
	id, err := strconv.Atoi(stringId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post id",
		})
		return
	}

	post := new(core.BlogPost)

	if err := h.Db.Preload("Attachments").Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while fetching post record",
		})
		return
	}

	tx := h.Db.Begin()

	if err := tx.Delete(core.BlogPost{ID: id}).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while deleting post record from database. Transaction rolled back",
		})
		return
	}

	if err := tx.Where("blog_post_id = ?", post.ID).Delete(&core.AttachmentRecord{}).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while deleting image records for post. Transaction rolled back",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Committing transaction failed",
		})
		return
	}

	for _, old := range post.Attachments {
		if err := os.Remove(old.Path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Error while deleting one of the attachments from assets",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post and its attachments deleted successfully",
	})
}
