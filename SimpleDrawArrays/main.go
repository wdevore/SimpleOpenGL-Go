package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 600
	height = 500
)

var (
	square = []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.5, 0.5, 0.0,
		0.5, 0.5, 0.0,
		-0.5, 0.5, 0.0,
		-0.5, -0.5, 0.0,
	}
)

var quitTriggered bool
var polygonMode bool

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	window.SetKeyCallback(keyCallback)

	program := initOpenGL()
	vao := makeVao(square)

	for !window.ShouldClose() && !quitTriggered {
		draw(vao, window, program)
	}
}

func draw(vao uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))

	glfw.PollEvents()
	window.SwapBuffers()
}
