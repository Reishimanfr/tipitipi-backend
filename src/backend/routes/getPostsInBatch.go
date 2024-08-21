package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (h *Handler) GetPostsInBatch(ctx *gin.Context) {
	sortBy := ctx.Param("sort_by")
	limitString := ctx.Param("limit")

	if !slices.Contains([]string{"likes", "newest", "oldest"}, sortBy) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid sort_by value provided (expected one of: [likes, newest, oldest])",
		})
		return
	}

	limit, err := strconv.Atoi(limitString)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid limit provided",
		})
		return
	}

	var postRecords *[]core.BlogPost
	orderClause := clause.OrderByColumn{
		Desc: true,
		Column: clause.Column{
			Table: "likes",
		},
	}

	if sortBy != "likes" {
		orderClause = clause.OrderByColumn{
			Desc: sortBy == "newest",
			Column: clause.Column{
				Table: "created_at",
			},
		}
	}

	h.Db.Order(orderClause).Limit(limit).Find(&postRecords)

	if len(*postRecords) < 1 {
		ctx.JSON(http.StatusNoContent, gin.H{
			"error": "not posts found in the database",
		})
		return
	}

	ctx.JSON(http.StatusOK, postRecords)
}
