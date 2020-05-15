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
	square = []float32{
		0.5, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,
	}
	triangle = []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.314, 0.0,
	}

	indicesSqr = []uint32{
		0, 3, 2, // first triangle
		2, 1, 0, // second triangle
	}
	indicesTri = []uint32{
		0, 1, 2,
	}
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

	sqrVao := buildSquare()

	triVao := buildTri()

	gl.Viewport(0, 0, width, height)

	gl.ClearColor(0.25, 0.25, 0.25, 1.0)

	gl.UseProgram(program)

	for !window.ShouldClose() && !quitTriggered {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		bindVao(sqrVao)
		draw(sqrVao, indicesSqr, window, program)

		bindVao(triVao)
		draw(triVao, indicesTri, window, program)

		glfw.PollEvents()
		window.SwapBuffers()
		time.Sleep(time.Millisecond)
	}
}

func draw(vao uint32, indices []uint32, window *glfw.Window, program uint32) {

	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)

	// if errNum := gl.GetError(); errNum != gl.NO_ERROR {
	// 	log.Fatal("(2)GL Error: ", errNum)
	// }

}

func buildSquare() uint32 {
	// Binding activates a buffer (i.e. a Scope)
	// VAO "captures" activity acted on buffers WHILE the VAO has
	// been activated via binding.

	// --------- Scope capturing STARTs here -------------------
	vao := makeVao()

	// Activate VBO buffer while in the VAOs scope
	vbo := makeVbo()

	bindVao(vao)

	bindVbo(square, vbo)

	// Activate EBO buffer while in the VAOs scope
	makeEbo(indicesSqr)

	gl.BindVertexArray(0) // close scope
	// --------- Scope capturing ENDs here -------------------

	return vao
}

func buildTri() uint32 {
	// --------- Scope capturing STARTs here -------------------
	vao := makeVao()

	// Activate VBO buffer while in the VAOs scope
	vbo := makeVbo()

	bindVao(vao)

	bindVbo(triangle, vbo)

	// Activate EBO buffer while in the VAOs scope
	makeEbo(indicesTri)

	gl.BindVertexArray(0) // close scope
	// --------- Scope capturing ENDs here -------------------

	return vao
}
