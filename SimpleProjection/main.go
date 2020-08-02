package main

import (
	"math"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width           = 1200
	height          = 800
	degreeToRadians = math.Pi / 180.0
)

var (
	quitTriggered    bool
	polygonMode      bool
	pointMode        bool
	quadVbo, quadVao uint32
)

func main() {
	runtime.LockOSThread()
	polygonMode = false

	window := initGlfw()
	defer glfw.Terminate()

	window.SetKeyCallback(keyCallback)

	defaultProgram := initDefaultProgram()

	triVao := buildTri()

	// -----------------------------------------------------------
	viewport := NewViewport()
	viewport.SetDimensions(0, 0, width, height)
	viewport.Apply()

	projection := buildProjection()

	view := buildView()

	// -----------------------------------------------------------
	model := NewMatrix4()
	// model.Rotate(45.0 * degreeToRadians)
	// model.TranslateBy2Comps(100.0, 0.0)
	model.ScaleByComp(100.0, 100.0, 1.0)

	mo := model.Matrix()

	// -----------------------------------------------------------
	gl.UseProgram(defaultProgram)

	projLoc := gl.GetUniformLocation(defaultProgram, gl.Str("projection\x00"))
	if projLoc < 0 {
		panic("NodeManager: couldn't find 'projection' uniform variable")
	}
	pm := projection.Matrix().Matrix()
	gl.UniformMatrix4fv(projLoc, 1, false, &pm[0])

	viewLoc := gl.GetUniformLocation(defaultProgram, gl.Str("view\x00"))
	if viewLoc < 0 {
		panic("NodeManager: couldn't find 'view' uniform variable")
	}
	gl.UniformMatrix4fv(viewLoc, 1, false, &view.Matrix()[0])

	modelLoc := gl.GetUniformLocation(defaultProgram, gl.Str("model\x00"))
	if modelLoc < 0 {
		panic("World: couldn't find 'model' uniform variable")
	}

	// -----------------------------------------------------------
	gl.ClearColor(0.25, 0.25, 0.25, 1.0)

	for !window.ShouldClose() && !quitTriggered {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UniformMatrix4fv(modelLoc, 1, false, &mo[0])

		bindVao(triVao)
		drawShape(triVao, indicesTri, window, defaultProgram)

		glfw.PollEvents()
		window.SwapBuffers()

		time.Sleep(time.Millisecond)
	}
}

func drawShape(vao uint32, indices []uint32, window *glfw.Window, program uint32) {
	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

const (
	trX = 0.0
	trY = 0.0
)

var (
	triangle = []float32{
		-0.5 + trX, -0.5 + trY, 0.0,
		0.5 + trX, -0.5 + trY, 0.0,
		0.0 + trX, 0.314 + trY, 0.0,
	}

	indicesTri = []uint32{
		0, 1, 2,
	}
)

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

func buildProjection() *Projection {
	projection := NewCamera()

	// This projection contains centering "built in". You don't need
	// a view matrix as long as your not interested in moving the view.
	// Otherwise you would need to rebuilt the projection. Typically
	// a second view-matrix is used for controlling a "camera".
	// sH := float32(height) / 2
	// sW := float32(width) / 2
	// projection.SetProjection(
	// 	-sH, -sW, // bottom,left
	// 	sH, sW, //top,right
	// 	-1.0, 1.0)

	// This projection is "just" a projection without any centering.
	projection.SetProjection(
		0.0, 0.0, // bottom,left
		float32(height), float32(width), //top,right
		-1.0, 1.0) //0.1, 100.0

	return projection
}

func buildView() IMatrix4 {
	centered := true
	offsetX := float32(0.0)
	offsetY := float32(0.0)

	if centered {
		offsetX = float32(width) / 2.0
		offsetY = float32(height) / 2.0
	}

	view := NewMatrix4()
	view.SetTranslate3Comp(offsetX, offsetY, 0.5)

	return view
}
