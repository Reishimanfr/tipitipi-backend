package routes

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func (h *Handler) validateJWT(c *gin.Context) {
    // The reason we can do this is that the JWT middleware
    // already handles everything for us and writing the same
    // code again would just be a waste of time
    c.JSON(http.StatusOK, gin.H{
        "message": "JWT token is valid",
    })
}