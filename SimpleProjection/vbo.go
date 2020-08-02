package main

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

// makeVbo initializes and returns a vertex buffer object from the points provided.
func makeVbo() uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	return vbo
}

func bindVbo(points []float32, vbo uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	// Specify the vertex attribute layout. This specifies--Where and How during transmission--opengl
	// will pass the data to the shader and how the shader will extract the data.
	vertexInputAttrb := uint32(0)
	sizeOfInputAttrb := int32(3)
	stride := int32(0)
	gl.VertexAttribPointer(vertexInputAttrb, sizeOfInputAttrb, gl.FLOAT, false, stride, nil)
	gl.EnableVertexAttribArray(vertexInputAttrb)
}
