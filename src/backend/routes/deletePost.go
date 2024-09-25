package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Record not found",
			})
			return
		}

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

	curDir, err := os.Getwd()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while getting the current work directory",
		})
		return
	}

	attDirectory := path.Join(curDir, "../assets", fmt.Sprintf("%v", post.ID))

	err = os.RemoveAll(attDirectory)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error while deleting attachments related to the deleted post",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post and its attachments deleted successfully",
	})
}
