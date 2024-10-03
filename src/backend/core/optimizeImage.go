package core

import (
	"bytes"

	"github.com/h2non/bimg"
	_ "golang.org/x/image/draw"
)

// TODO: make this work lol
func OptimizeImage(buffer []byte, quality int) (*bytes.Buffer, error) {
	converted, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		return nil, err
	}

	processed, err := bimg.NewImage(converted).Process(bimg.Options{
		Quality: quality,
	})

	return bytes.NewBuffer(processed), err
}
