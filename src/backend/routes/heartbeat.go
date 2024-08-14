package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Heartbeat(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}
