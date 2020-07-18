package main

import (
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 800
	height = 800
)

var quitTriggered bool
var polygonMode bool
var pointMode bool
var shapeID = 3

func main() {
	runtime.LockOSThread()
	polygonMode = false

	window := initGlfw()
	defer glfw.Terminate()

	window.SetKeyCallback(keyCallback)

	defaultProgram := initDefaultProgram()
	textureProgram := initTextureProgram()

	triVao := buildTri()

	quadVao, quadTbo := buildQuadTexture()

	gl.Viewport(0, 0, width, height)

	gl.ClearColor(0.25, 0.25, 0.25, 1.0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	for !window.ShouldClose() && !quitTriggered {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(defaultProgram)

		bindVao(triVao)
		drawShape(triVao, indicesTri, window, defaultProgram)

		bindTexture(quadTbo)
		bindVao(quadVao)
		gl.UseProgram(textureProgram)

		drawTexture(quadVao, quadTbo, indicesQuad, window, defaultProgram)

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

func drawTexture(vao, tbo uint32, indices []uint32, window *glfw.Window, program uint32) {
	gl.BindVertexArray(vao)
	gl.BindTexture(gl.TEXTURE_2D, tbo)
	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

var (
	triangle = []float32{
		-0.2, -0.5, 0.0,
		0.8, -0.5, 0.0,
		0.3, 0.314, 0.0,
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

var (
	quad = []float32{
		// Use these for the ctype ship
		// -0.5, -0.5, 0.0, 0.0, 0.0, // bottom left
		// 0.5, -0.5, 0.0, 1.0, 0.0, // bottom right
		// 0.5, 0.5, 0.0, 1.0, 1.0, // top right
		// -0.5, 0.5, 0.0, 0.0, 1.0, // top left
	}

	// Indices defined in CCW order
	indicesQuad = []uint32{
		// 0, 1, 3, // first triangle
		// 1, 2, 3, // second triangle
		// OR
		0, 1, 2, // first triangle
		0, 2, 3, // second triangle
	}
)

func buildQuadTexture() (vao, tbo uint32) {
	// --------- Scope capturing STARTs here -------------------
	vao = makeVao()

	// Activate VBO buffer while in the VAOs scope
	vbo := makeVbo()

	bindVao(vao)

	// vector-space
	//      -0.5,0.5
	//  ^     *-----------*  0.5,0.5
	//  |     |           |
	//  |     |     _     |
	//  |+Y   |           |
	//  |     |           |
	//  |     *-----------*  1,1
	//      -0.5,-0.5        0.5,-0.5
	//        A           B

	//       texture-space
	//     0,1 D           C 1,1
	//  |      *-----------*
	//  |      |           |
	//  |+Y    |     _     |
	//  |      |           |
	//  |      |           |
	//  v      *-----------*
	//       0,0           1,0
	//         A           B

	switch shapeID {
	case 1:
		// Selects Mine shape
		quad = append(quad, -0.5, -0.5, 0.0) // xy
		quad = append(quad, 0.0, 0.25)       // uv
		quad = append(quad, 0.5, -0.5, 0.0)
		quad = append(quad, 0.25, 0.25) // uv
		quad = append(quad, 0.5, 0.5, 0.0)
		quad = append(quad, 0.25, 0.5) // uv
		quad = append(quad, -0.5, 0.5, 0.0)
		quad = append(quad, 0.0, 0.5) // uv
	case 2:
		// Selects green ship
		quad = append(quad, -0.5, -0.5, 0.0) // xy
		quad = append(quad, 0.25, 0.75)      // uv
		quad = append(quad, 0.5, -0.5, 0.0)
		quad = append(quad, 0.5, 0.75) // uv
		quad = append(quad, 0.5, 0.5, 0.0)
		quad = append(quad, 0.5, 1.0) // uv
		quad = append(quad, -0.5, 0.5, 0.0)
		quad = append(quad, 0.25, 1.0) // uv
	case 3:
		// Selects upper right orange ship
		quad = append(quad, -0.5, -0.5, 0.0) // xy
		quad = append(quad, 0.75, 0.75)      // uv
		quad = append(quad, 0.5, -0.5, 0.0)
		quad = append(quad, 1.0, 0.75) // uv
		quad = append(quad, 0.5, 0.5, 0.0)
		quad = append(quad, 1.0, 1.0) // uv
		quad = append(quad, -0.5, 0.5, 0.0)
		quad = append(quad, 0.75, 1.0) // uv
	case 4:
		// Selects upper left ctype ship
		quad = append(quad, -0.5, -0.5, 0.0) // xy
		quad = append(quad, 0.0, 0.75)       // uv
		quad = append(quad, 0.5, -0.5, 0.0)
		quad = append(quad, 0.25, 0.75) // uv
		quad = append(quad, 0.5, 0.5, 0.0)
		quad = append(quad, 0.25, 1.0) // uv
		quad = append(quad, -0.5, 0.5, 0.0)
		quad = append(quad, 0.0, 1.0) // uv
	case 5:
		// Selects smaller bomb
		quad = append(quad, -0.5, -0.5, 0.0) // xy
		quad = append(quad, 0.75, 0.625)     // uv
		quad = append(quad, 0.5, -0.5, 0.0)
		quad = append(quad, 0.8125, 0.625) // uv
		quad = append(quad, 0.5, 0.5, 0.0)
		quad = append(quad, 0.8125, 0.6875) // uv
		quad = append(quad, -0.5, 0.5, 0.0)
		quad = append(quad, 0.75, 0.6875) // uv
	}

	bindVbo(quad, vbo)

	bindTextureVbo(quad, vbo)

	// Activate EBO buffer while in the VAOs scope
	makeEbo(indicesQuad)

	tbo = makeTbo()

	image, err := loadImage("texture-atlas.png")
	if err != nil {
		panic(err)
	}

	bindTbo(image, tbo)

	gl.BindVertexArray(0) // close scope
	// --------- Scope capturing ENDs here -------------------

	return vao, tbo
}

// Lerp returns a the value between min and max given t = 0->1
func lerp(min, max, t float64) float64 {
	return min*(1.0-t) + max*t
}
