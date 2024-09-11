package middleware

import (
	"github.com/gin-gonic/gin"
)

func FileSizeLimiterMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// if c.Request.Method == http.MethodPost {
		// 	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		// 	if err := c.Request.ParseMultipartForm(maxSize); err != nil {
		// 		c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Request body too large"})
		// 		return
		// 	}

		// 	for _, files := range c.Request.MultipartForm.File {
		// 		for _, file := range files {
		// 			if file.Size > maxSize {
		// 				c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File size exceeds limit"})
		// 				return
		// 			}
		// 		}
		// 	}
		// }

		c.Next()
	}
}
