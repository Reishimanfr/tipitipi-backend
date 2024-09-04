package middleware

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SizeLimiter struct {
	// Max body size in bytes
	maxSize int64
}

func (s *SizeLimiter) Allow(c *gin.Context) bool {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, s.maxSize)

	_, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if err == http.ErrBodyReadAfterClose || err == io.ErrUnexpectedEOF {
			return false
		}
		fmt.Println("Failed to read response body ", err)
		return false
	}

	return true
}

func NewSizeLimiter(maxSize int64) *SizeLimiter {
	limiter := &SizeLimiter{
		maxSize: maxSize,
	}

	return limiter
}

func BodySizeLimiterMiddleware(limiter *SizeLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow(c) {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "Request body too large",
			})
			return
		}
		c.Next()
	}
}
