package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	sortByOptions = []string{
		"newest",
		"oldest",
		"likes",
	}
)

func (h *Handler) posts(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "5")
	sort := c.DefaultQuery("sort", "newest")
	images := c.DefaultQuery("images", "false")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid offset value provided (expected an int)",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid offset value provided (expected an int)",
		})
		return
	}

	if !slices.Contains(sortByOptions, sort) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid sort option. Expected one of " + strings.Join(sortByOptions, ", "),
		})
		return
	}

	postRecords := make([]*core.BlogPost, limit)

	orderClause := clause.OrderByColumn{
		Desc: true,
		Column: clause.Column{
			Name: "likes",
		},
	}

	if sort != "likes" {
		orderClause = clause.OrderByColumn{
			Desc: false,
			Column: clause.Column{
				Name: "created_at",
			},
		}
	}

	result := new(gorm.DB)

	if images == "true" {
		result = h.Db.Preload("Images").Order(orderClause).Offset(offset).Limit(limit).Find(&postRecords)
	} else {
		result = h.Db.Order(orderClause).Offset(offset).Limit(limit).Find(&postRecords)
	}

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		h.Log.Error("Error while getting post records from database", zap.Error(result.Error))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   result.Error.Error(),
			"message": "Error while getting records from database",
		})
		return
	}

	if len(postRecords) < 1 {
		c.AbortWithStatusJSON(http.StatusNoContent, gin.H{
			"error": "Post not found",
		})
		return
	}

	c.JSON(http.StatusOK, postRecords)
}
