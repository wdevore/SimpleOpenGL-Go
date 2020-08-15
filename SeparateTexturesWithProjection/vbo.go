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

	// Specify the vertex attribute layout. This specifies--how during transmission--opengl
	// will pass the data to the shader and how the shader will extract the data.
	vertexInputAttrb := uint32(0)
	sizeOfInputAttrb := int32(3)
	stride := int32(0)
	gl.VertexAttribPointer(vertexInputAttrb, sizeOfInputAttrb, gl.FLOAT, false, stride, nil)
	gl.EnableVertexAttribArray(vertexInputAttrb)
}

func bindTextureVbo(points *[]float32, vbo uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(*points), gl.Ptr(*points), gl.DYNAMIC_DRAW)

	sizeOfFloat := int32(4)

	// If the data per-vertex is (x,y,z,s,t = 5) then
	// Stride = 5 * size of float
	// OR
	// If the data per-vertex is (x,y,z,r,g,b,s,t = 8) then
	// Stride = 8 * size of float

	// Our data layout is x,y,z,s,t
	stride := 5 * sizeOfFloat

	// position attribute
	size := int32(3)   // x,y,z
	offset := int32(0) // position is first thus this attrib is offset by 0
	attribIndex := uint32(0)
	gl.VertexAttribPointer(attribIndex, size, gl.FLOAT, false, stride, gl.PtrOffset(int(offset)))
	gl.EnableVertexAttribArray(0)

	// texture coord attribute is offset by 3 (i.e. x,y,z)
	size = int32(2)   // s,t
	offset = int32(3) // the preceeding component size = 3, thus this attrib is offset by 3
	attribIndex = uint32(1)
	gl.VertexAttribPointer(attribIndex, size, gl.FLOAT, false, stride, gl.PtrOffset(int(offset*sizeOfFloat)))
	gl.EnableVertexAttribArray(1)
}

// Update moves any modified data to the buffer.
func updateTextureVbo(data *[]float32, vbo uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(*data)*4, gl.Ptr(*data))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
