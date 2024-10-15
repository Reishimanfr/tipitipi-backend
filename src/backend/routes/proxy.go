package routes

import (
    "io"
    "net/http"
    "strconv"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/gin-gonic/gin"
)

func (h *Handler) proxy(c *gin.Context) {
    key := c.Query("key")

    if key == "" {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "error": "key is nil",
        })
        return
    }

    obj, err := h.Ovh.S3.GetObject(&s3.GetObjectInput{
        Bucket: aws.String("tipi-tipi-test-container"),
        Key:    aws.String(key),
    })
    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to proxy file from S3",
        })
        return
    }

    defer obj.Body.Close()

    c.Header("Content-Type", *obj.ContentType)
    c.Header("Content-Length", strconv.Itoa(int(*obj.ContentLength)))
    c.Header("Cache-Control", "max-age=3600")

    body, err := io.ReadAll(obj.Body)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to read file",
        })
        return
    }

    c.Data(http.StatusOK, *obj.ContentType, body)
}