package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (h *Handler) posts(ctx *gin.Context) {
	sortBy := ctx.PostForm("sortBy")
	amount := ctx.PostForm("amount")
	offsetStr := ctx.PostForm("offset")

	if !slices.Contains([]string{"likes", "newest", "oldest"}, sortBy) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid sortBy value provided (expected one of: likes, newest, oldest)",
		})
		return
	}

	limit, err := strconv.Atoi(amount)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post limit",
		})
		return
	}

	if offsetStr == "" {
		offsetStr = "0"
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid offset",
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

	if sortBy != "likes" {
		orderClause = clause.OrderByColumn{
			Desc: false,
			Column: clause.Column{
				Name: "created_at",
			},
		}
	}

	h.Db.Order(orderClause).Offset(offset).Limit(limit).Find(&postRecords)

	if len(*postRecords) < 1 {
		ctx.JSON(http.StatusNoContent, gin.H{
			"error": "No posts found",
		})
		return
	}

	ctx.JSON(http.StatusOK, postRecords)
}
