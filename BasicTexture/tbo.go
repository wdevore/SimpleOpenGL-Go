package main

import (
	"image"
	"image/draw"
	_ "image/png" // For 'png' images
	"os"

	"github.com/go-gl/gl/v4.5-core/gl"
)

// makeTbo initializes
func makeTbo() uint32 {
	var tbo uint32
	gl.GenTextures(1, &tbo)
	return tbo
}

func bindTexture(tbo uint32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tbo)
}

func bindTbo(texture *image.NRGBA, tbo uint32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tbo)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	width := int32(texture.Bounds().Dx())
	height := int32(texture.Bounds().Dy())

	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	// Give the image to OpenGL
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(texture.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
}

func loadImage(path string) (*image.NRGBA, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	nrgba := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))

	r := image.Rect(0, 0, bounds.Dx(), bounds.Dy())
	flippedImg := image.NewNRGBA(r)

	// Transfer data to image
	draw.Draw(nrgba, nrgba.Bounds(), img, bounds.Min, draw.Src)

	// Flip horizontally or along Y-axis
	// for j := 0; j < nrgba.Bounds().Dy(); j++ {
	// 	for i := 0; i < nrgba.Bounds().Dx(); i++ {
	// 		flippedImg.Set(bounds.Dx()-i, j, nrgba.At(i, j))
	// 	}
	// }

	for j := 0; j < nrgba.Bounds().Dy(); j++ {
		for i := 0; i < nrgba.Bounds().Dx(); i++ {
			flippedImg.Set(i, bounds.Dy()-j, nrgba.At(i, j))
		}
	}

	return flippedImg, nil
}
