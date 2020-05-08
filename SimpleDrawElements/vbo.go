package main

import "github.com/go-gl/gl/v4.5-core/gl"

// makeVbo initializes and returns a vertex buffer object from the points provided.
func makeVbo(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	return vbo
}
