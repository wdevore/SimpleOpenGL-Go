package main

import (
	"log"

	"github.com/go-gl/gl/v4.5-core/gl"
)

func makeEbo(indices []uint32) uint32 {
	var ebo uint32
	gl.GenBuffers(1, &ebo)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	sizeOfUInt32 := 4
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, sizeOfUInt32*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	if errNum := gl.GetError(); errNum != gl.NO_ERROR {
		log.Fatal("(ebo)GL Error: ", errNum)
	}

	return ebo
}
