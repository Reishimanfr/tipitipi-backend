package core

import (
	"net/http"
	"strings"

	"github.com/h2non/bimg"
	_ "golang.org/x/image/draw"
)

// Attempts to optimize the provided buffer based on it's detected MIME type
func OptimizeBuffer(b []byte, quality int) ([]byte, error) {
	mime := http.DetectContentType(b)

	if strings.HasPrefix(mime, "image/") {
		return bimg.NewImage(b).Process(bimg.Options{
			Quality:  quality,
			Lossless: false,
			Type:     bimg.JPEG,
		})
	}

	// Don't do anything if the buffer is not an image
	return b, nil
}
