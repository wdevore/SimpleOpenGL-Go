package main

import (
	"log"
	"math"
	"runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 800
	height = 450

	degreeToRadians = math.Pi / 180.0

	centered = true
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

var projection *matrix4
var view *matrix4
var model *matrix4
var modelOrb *matrix4

var projLoc int32
var viewLoc int32
var modelLoc int32

var angle float64

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
	vertexInputAttrb := uint32(0) // aPos
	sizeOfInputAttrb := int32(3)  // 3 floats. The 4th is supplied by the shader
	stride := int32(0)            // All the vertices will tightly packed
	gl.VertexAttribPointer(vertexInputAttrb, sizeOfInputAttrb, gl.FLOAT, false, stride, nil)
	gl.EnableVertexAttribArray(vertexInputAttrb)
	// --------- Scope capturing ENDs here -------------------

	// This is a manual viewport configuration that insures a square is "square"
	// if width > height {
	// 	gl.Viewport(0, -int32((width-height)/2.0), width, width)
	// } else {
	// 	gl.Viewport(int32((width-height)/2.0), 0, height, height)
	// }
	// But normally you adjust the projection instead by using window
	// dimensions instead.
	gl.Viewport(0, 0, width, height)

	gl.UseProgram(program)

	projLoc = gl.GetUniformLocation(program, gl.Str("projection\x00"))
	if projLoc < 0 {
		log.Fatal("no projLoc")
	}
	viewLoc = gl.GetUniformLocation(program, gl.Str("view\x00"))
	modelLoc = gl.GetUniformLocation(program, gl.Str("model\x00"))

	projection = NewMatrix4()
	view = NewMatrix4()
	model = NewMatrix4()
	modelOrb = NewMatrix4()

	projection.SetToOrtho(0.0, width, 0.0, height, 0.1, 100.0)

	// The View matrix will contain a translation if centering

	offsetX := float32(0.0)
	offsetY := float32(0.0)
	if centered {
		offsetX = width / 2.0
		offsetY = height / 2.0
	}
	view.SetTranslate3Comp(offsetX, offsetY, 1.0)
	// view.ScaleByComp(2.0, 2.0, 1.0)

	pm := projection.Matrix()
	gl.UniformMatrix4fv(projLoc, 1, false, &pm[0])

	vm := view.Matrix()
	gl.UniformMatrix4fv(viewLoc, 1, false, &vm[0])

	gl.ClearColor(0.25, 0.25, 0.25, 1.0)

	for !window.ShouldClose() && !quitTriggered {
		draw(vao, window, program)
	}
}

func draw(vao uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)

	orbitApproach3()

	gl.UniformMatrix4fv(modelLoc, 1, false, &model.Matrix()[0])

	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

	angle += 2.0

	gl.BindVertexArray(0)

	glfw.PollEvents()
	window.SwapBuffers()
}

//------------------------------------------------------------------
// Each orbit approach produces the same results: the square orbits the
// center of the window CCW. The square does not rotate because its rotation
// is the opposite with respect to the orbit.
//------------------------------------------------------------------

func orbitApproach1() {
	model.SetScale3Comp(1.0, 1.0, 1.0)

	model.Rotate(degreeToRadians * angle)
	model.TranslateBy3Comps(50.0, 0.0, 0.0)

	model.TranslateBy3Comps(0.0, 0.0, 0.0)
	model.Rotate(degreeToRadians * -angle)
	model.ScaleByComp(30.0, 30.0, 1.0)
}

func orbitApproach2() {
	model.SetScale3Comp(1.0, 1.0, 1.0)

	model.TranslateBy3Comps(0.0, 0.0, 0.0)
	model.Rotate(degreeToRadians * -angle)
	model.ScaleByComp(30.0, 30.0, 1.0)

	// Method #2
	modelOrb.SetScale3Comp(1.0, 1.0, 1.0)
	modelOrb.Rotate(degreeToRadians * angle)
	modelOrb.TranslateBy3Comps(50.0, 0.0, 0.0)

	model.PostMultiply(modelOrb)
}
func orbitApproach3() {
	model.SetScale3Comp(1.0, 1.0, 1.0)
	model.TranslateBy3Comps(0.0, 0.0, 0.0)
	model.Rotate(degreeToRadians * -angle)
	model.ScaleByComp(30.0, 30.0, 1.0)

	// Or method #3 using Inverse
	modelOrb.SetTranslate3Comp(-50.0, 0.0, 0.0)
	modelOrb.Rotate(degreeToRadians * -angle)
	modelOrb.ScaleByComp(1.0, 1.0, 1.0)
	modelOrb.Inverse()

	model.PostMultiply(modelOrb)
}

// Used if you want to force a dimension different than the window
func configureProjections(deviceWidth, deviceHeight, virtualWidth, virtualHeight int) (ratioCorrection float64) {

	// Calc the aspect ratio between the physical (aka device) dimensions and the
	// the virtual (aka user's design choice) dimensions.

	deviceRatio := float64(deviceWidth) / float64(deviceHeight)
	virtualRatio := float64(virtualWidth) / float64(virtualHeight)

	xRatioCorrection := float64(deviceWidth) / float64(virtualWidth)
	yRatioCorrection := float64(deviceHeight) / float64(virtualHeight)

	if virtualRatio < deviceRatio {
		ratioCorrection = yRatioCorrection
	} else {
		ratioCorrection = xRatioCorrection
	}

	return
}

// if errNum := gl.GetError(); errNum != gl.NO_ERROR {
// 	log.Fatal("(1)GL Error: ", errNum)
// }
