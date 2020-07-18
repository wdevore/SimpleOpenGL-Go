package main

import (
	"image"
	_ "image/png" // For 'png' images

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
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

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
