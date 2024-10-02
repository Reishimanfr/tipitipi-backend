package core

import (
	"image"
	"os"
	"path"
	"path/filepath"

	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/draw"
)

// TODO: finish thumbnails
func MakeThumbnail(fpath string, width int, height int) error {
	inputFile, err := os.Open(fpath)
	if err != nil {
		return err
	}

	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return err
	}

	thumbnail := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.ApproxBiLinear.Scale(thumbnail, thumbnail.Bounds(), img, img.Bounds(), draw.Over, nil)

	thumbFilename := "thumb-" + filepath.Base(fpath)
	thumbDirname := path.Join(filepath.Dir(fpath), "thumbnails", thumbFilename)

	outputFile, err := os.Create(thumbDirname)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	// err = webp.Decode(outputFile, thumbnail, nil)

	return nil
}
