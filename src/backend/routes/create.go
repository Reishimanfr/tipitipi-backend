package routes

import (
	"bash06/strona-fundacja/src/backend/core"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
			"error": "Missing title field",
		})
		return false
	}

	if len(data.Title) > 255 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Field title is too long (max: 255 characters)",
		})
		return false
	}

	if data.Content == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing content field",
		})
		return false
	}

	// var seenIds []int

	// for _, v := range data.Images {
	// 	fmt.Println(seenIds)
	// 	if slices.Contains(seenIds, v.Id) {
	// 		ctx.JSON(http.StatusConflict, gin.H{
	// 			"error": "duplicate image ids found",
	// 		})
	// 		return false
	// 	}

	// 	seenIds = append(seenIds, v.Id)
	// }

	return true
}

func (h *Handler) create(ctx *gin.Context) {
	var data CreateRequestBody

	if err := ctx.BindJSON(&data); err != nil {
		h.Log.Error("Failed to bind json data", zap.Error(err))
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Failed to bind json data",
		})
		return
	}

	sanitizeData(&data, ctx)

	var count int64
	result := h.Db.Model(&core.BlogPost{}).Where("title = ?", data.Title).Count(&count)

	if result.Error != nil {
		h.Log.Error("Error while checking for post in database", zap.Error(result.Error))
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Error while checking if post exists",
		})
		return
	}

	if count > 0 {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Post with this title already exists",
		})
		return
	}

	newPost := &core.BlogPost{
		Title:      data.Title,
		Content:    data.Content,
		Created_At: time.Now().Unix(),
		Edited_At:  time.Now().Unix(),
		Images:     "", // TODO
	}

	h.Db.Create(&newPost)
	ctx.JSON(http.StatusOK, newPost)
}
