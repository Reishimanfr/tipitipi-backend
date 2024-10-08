package core

import (
	"github.com/h2non/bimg"
	_ "golang.org/x/image/draw"
)

func OptimizeAttachment(b []byte, quality int) ([]byte, error) {
	options := bimg.Options{
		Quality:  quality,
		Lossless: false,
		Type:     bimg.WEBP,
	}

	return bimg.NewImage(b).Process(options)
}
