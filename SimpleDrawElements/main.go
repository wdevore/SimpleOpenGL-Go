package main

import (
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 800
	height = 500
)

//  +Y
//   ^
//   |
//   |
//   |
//   .--------> +X
//  0,0

var (
	// square = []float32{
	// 	-0.5, -0.5, 0.0,
	// 	0.5, -0.5, 0.0,
	// 	0.5, 0.5, 0.0,
	// 	0.5, 0.5, 0.0,
	// 	-0.5, 0.5, 0.0,
	// 	-0.5, -0.5, 0.0,
	// }

	square = []float32{
		0.5, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,
	}
	// square = []float32{
	// 	-0.5, -0.5, 0.0,
	// 	-0.5, 0.5, 0.0,
	// 	0.5, 0.5, 0.0,
	// 	0.5, -0.5, 0.0,
	// }

	indices = []uint32{
		0, 3, 2, // first triangle
		2, 1, 0, // second triangle
	}
	// indices = []uint32{
	// 	0, 1, 3, // first triangle
	// 	1, 2, 3, // second triangle
	// }
	// indices = []uint32{
	// 	0, 1, 2,
	// 	2, 3, 0,
	// }
)

var quitTriggered bool
var polygonMode bool
var pointMode bool

func main() {
	runtime.LockOSThread()
	polygonMode = false

	window := initGlfw()
	defer glfw.Terminate()

	window.SetKeyCallback(keyCallback)

	program := initOpenGL()

	// Binding activates a buffer (i.e. a Scope)
	// VAO "captures" activity acted on buffers WHILE the VAO has
	// been activated via binding.
	// --------- Scope capturing STARTs here -------------------
	vao := makeVao()

	// Activate VBO buffer while in the VAOs scope
	makeVbo(square)

	// Activate EBO buffer while in the VAOs scope
	makeEbo(indices)

	// Specify the vertex attribute layout. This specifies--how during transmission--opengl
	// will pass the data to the shader and how the shader will extract the data.
	vertexInputAttrb := uint32(0)
	sizeOfInputAttrb := int32(3)
	stride := int32(0)
	gl.VertexAttribPointer(vertexInputAttrb, sizeOfInputAttrb, gl.FLOAT, false, stride, nil)
	gl.EnableVertexAttribArray(vertexInputAttrb)
	// --------- Scope capturing ENDs here -------------------

	// This is a viewport configuration insures that a square is "square"
	// if width > height {
	// 	gl.Viewport(0, -int32((width-height)/2.0), width, width)
	// } else {
	// 	gl.Viewport(int32((width-height)/2.0), 0, height, height)
	// }
	// But normally you adjust the projection instead and use window
	// dimensions instead.
	gl.Viewport(0, 0, width, height)

	// The View matrix will contain a translation if centering
	// and possibly a scale

	// aspect := float32(width) / float32(height)

	gl.ClearColor(0.25, 0.25, 0.25, 1.0)

	for !window.ShouldClose() && !quitTriggered {
		draw(vao, window, program)
		time.Sleep(time.Millisecond)
	}
}

func draw(vao uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)

	// if errNum := gl.GetError(); errNum != gl.NO_ERROR {
	// 	log.Fatal("(1)GL Error: ", errNum)
	// }

	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

	// if errNum := gl.GetError(); errNum != gl.NO_ERROR {
	// 	log.Fatal("(2)GL Error: ", errNum)
	// }

	gl.BindVertexArray(0)

	glfw.PollEvents()
	window.SwapBuffers()
}
