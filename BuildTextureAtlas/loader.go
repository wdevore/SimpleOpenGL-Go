package main

import (
	"bufio"
	"image"
	"image/draw"
	_ "image/png" // For 'png' images
	"os"
)

func loadManifest(maniFile string) []string {
	file, err := os.Open(maniFile)

	defer file.Close()

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	return txtlines
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func imageToNRGBA(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	nrgba := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))

	r := image.Rect(0, 0, bounds.Dx(), bounds.Dy())
	nrg := image.NewNRGBA(r)

	// Transfer data to image
	draw.Draw(nrg, nrgba.Bounds(), img, bounds.Min, draw.Src)

	return nrgba
}

func flipImage(img *image.NRGBA) *image.NRGBA {
	r := image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy())
	flippedImg := image.NewNRGBA(r)

	// Flip vertically or around the X-axis
	height := img.Bounds().Dy()
	for j := 0; j < height; j++ {
		for i := 0; i < img.Bounds().Dx(); i++ {
			flippedImg.Set(i, height-j, img.At(i, j))
		}
	}

	return flippedImg
}
