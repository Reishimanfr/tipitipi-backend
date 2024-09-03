package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Heartbeat(c *gin.Context) {
	c.Status(http.StatusOK)
}
