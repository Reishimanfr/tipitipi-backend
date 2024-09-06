package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"fmt"
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

	if err := h.Db.Preload("Images").Where("id = ?", id).First(post).Error; err != nil {
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
			"message": "Error while deleting post record from database",
		})
		return
	}

	if err := tx.Where("blog_post_id = ?", post.ID).Delete(&core.ImageRecord{}).Error; err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while deleting image records for post",
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

	fmt.Println(post.Images)

	for _, oldImage := range post.Images {
		fmt.Println(oldImage.Path)
		if err := os.Remove(oldImage.Path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Error while removing post image",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post and its images deleted successfully",
	})
}
