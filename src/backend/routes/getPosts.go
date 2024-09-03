package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostsInBatchBody struct {
	SortBy string `json:"sortBy" binding:"required"`
	Limit  uint8  `json:"limit" binding:"required,max=10,min=1"`
	Offset uint16 `json:"offset"`
}

func (h *Handler) posts(c *gin.Context) {
	var body PostsInBatchBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "Malformed or invalid JSON",
		})
		return
	}

	if !slices.Contains([]string{"likes", "newest", "oldest"}, body.SortBy) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid sortBy value. Expected one of likes, newest, oldest, got " + body.SortBy + " instead",
		})
		return
	}

	var postRecords *[]core.BlogPost
	orderClause := clause.OrderByColumn{
		Desc: true,
		Column: clause.Column{
			Name: "likes",
		},
	}

	if body.SortBy != "likes" {
		orderClause = clause.OrderByColumn{
			Desc: false,
			Column: clause.Column{
				Name: "created_at",
			},
		}
	}

	result := h.Db.Order(orderClause).Offset(int(body.Offset)).Limit(int(body.Limit)).Find(&postRecords)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": "Error while getting records from database",
		})
		return
	}

	if len(*postRecords) < 1 {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{
			"error": "Post not found",
		})
		return
	}

	c.JSON(http.StatusOK, postRecords)
}
