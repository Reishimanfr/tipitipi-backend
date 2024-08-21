package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestBodyImageData struct {
	Id       int    `json:"id"`
	Path     string `json:"path"`
	Position string `json:"position"`
}

type CreateRequestBody struct {
	Title   string                 `json:"title"`
	Content string                 `json:"content"`
	Images  []RequestBodyImageData `json:"images"`
}

func sanitizeData(data *CreateRequestBody, ctx *gin.Context) bool {
	if strings.Trim(data.Title, "") == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "(required) title field is empty",
		})
		return false
	}

	if len(data.Title) > 255 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "title field is too long (max: 255 characters)",
		})
		return false
	}

	if data.Content == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "(required) field content is empty",
		})
		return false
	}

	var seenIds []int

	for _, v := range data.Images {
		if slices.Contains(seenIds, v.Id) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "duplicate image ids found",
			})
			return false
		}

		seenIds = append(seenIds, v.Id)
	}

	return true
}

func (h *Handler) CreateBlogPost(ctx *gin.Context) {
	var data CreateRequestBody

	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "failed to bind json data: " + err.Error(),
		})
		return
	}

	sanitizeData(&data, ctx)

	var postRecord *core.BlogPost

	h.Db.Where("title = ?", data.Title).First(&postRecord)

	if postRecord != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "post with this name already exists",
		})
		return
	}

	newPost := &core.BlogPost{
		Title:      data.Title,
		Content:    data.Content,
		Created_At: time.Now(),
		Edited_At:  time.Now(),
		Images:     "", // TODO
	}

	h.Db.Create(&newPost)
	ctx.JSON(http.StatusOK, newPost)
}
