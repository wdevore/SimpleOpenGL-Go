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
	quitTriggered                   bool
	polygonMode                     bool
	pointMode                       bool
	textureAtlas                    = NewTextureAtlas("texture_manifest.txt")
	texture2Atlas                   = NewTextureAtlas("texture2_manifest.txt")
	quadVbo, quadVao, quadTbo       uint32
	quadVbo2, quadVao2, quadTbo2    uint32
	actQuadVbo                      uint32
	textureProgram, texture2Program uint32
)

func main() {
	runtime.LockOSThread()
	polygonMode = false

	window := initGlfw()
	defer glfw.Terminate()

	window.SetKeyCallback(keyCallback)

	defaultProgram := initDefaultProgram()
	textureProgram = initTextureProgram()
	texture2Program = initTextureProgram()

	triVao := buildTri()

	quadVao, quadTbo, quadVbo = buildQuadTexture(textureAtlas, "green ship", &indicesQuad, &quad)
	changeShape("ctype ship", quadVbo)

	quadVao2, quadTbo2, quadVbo2 = buildQuadTexture(texture2Atlas, "orange ship", &indicesQuad2, &quad2)
	changeShape("orange ship", quadVbo2)

	actQuadVbo = quadVbo

	// -----------------------------------------------------------
	viewport := NewViewport()
	viewport.SetDimensions(0, 0, width, height)
	viewport.Apply()

	projection := buildProjection()

	view := buildView()

	// -----------------------------------------------------------
	angle := 0.0
	modelTri := NewMatrix4()
	modelTri.SetRotation(angle * degreeToRadians)
	modelTri.TranslateBy2Comps(100.0, 0.0)
	modelTri.ScaleByComp(100.0, 100.0, 1.0)
	mo := modelTri.Matrix()

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

	modelTriLoc := gl.GetUniformLocation(defaultProgram, gl.Str("model\x00"))
	if modelTriLoc < 0 {
		panic("World: couldn't find 'model' uniform variable")
	}

	// -----------------------------------------------------------
	modelQuad := NewMatrix4()
	modelQuad.ScaleByComp(64.0, 64.0, 1.0)
	modelQuad.TranslateBy2Comps(-5.0, 0.0)
	moq := modelQuad.Matrix()

	modelQuad2 := NewMatrix4()
	modelQuad2.ScaleByComp(64.0, 64.0, 1.0)
	modelQuad2.TranslateBy2Comps(5.0, 0.0)
	moq2 := modelQuad2.Matrix()

	gl.UseProgram(textureProgram)

	projLoc = gl.GetUniformLocation(textureProgram, gl.Str("projection\x00"))
	if projLoc < 0 {
		panic("NodeManager: couldn't find 'projection' uniform variable")
	}
	pm = projection.Matrix().Matrix()
	gl.UniformMatrix4fv(projLoc, 1, false, &pm[0])

	viewLoc = gl.GetUniformLocation(textureProgram, gl.Str("view\x00"))
	if viewLoc < 0 {
		panic("NodeManager: couldn't find 'view' uniform variable")
	}
	gl.UniformMatrix4fv(viewLoc, 1, false, &view.Matrix()[0])

	modelQuadLoc := gl.GetUniformLocation(textureProgram, gl.Str("model\x00"))
	if modelQuadLoc < 0 {
		panic("World: couldn't find 'model' uniform variable")
	}

	// -------------------------------------------------------------
	gl.UseProgram(texture2Program)

	projLoc2 := gl.GetUniformLocation(texture2Program, gl.Str("projection\x00"))
	if projLoc2 < 0 {
		panic("NodeManager: couldn't find 'projection' uniform variable")
	}
	pm = projection.Matrix().Matrix()
	gl.UniformMatrix4fv(projLoc2, 1, false, &pm[0])

	viewLoc2 := gl.GetUniformLocation(texture2Program, gl.Str("view\x00"))
	if viewLoc2 < 0 {
		panic("NodeManager: couldn't find 'view' uniform variable")
	}
	gl.UniformMatrix4fv(viewLoc2, 1, false, &view.Matrix()[0])

	modelQuadLoc2 := gl.GetUniformLocation(texture2Program, gl.Str("model\x00"))
	if modelQuadLoc2 < 0 {
		panic("World: couldn't find 'model' uniform variable")
	}

	// -----------------------------------------------------------

	gl.ClearColor(0.25, 0.25, 0.25, 1.0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	for !window.ShouldClose() && !quitTriggered {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(defaultProgram)
		gl.UniformMatrix4fv(modelTriLoc, 1, false, &mo[0])
		bindVao(triVao)
		drawShape(triVao, indicesTri, window, defaultProgram)

		modelTri.SetRotation(angle * degreeToRadians)
		modelTri.TranslateBy2Comps(100.0, 0.0)
		modelTri.ScaleByComp(25.0, 25.0, 1.0)
		angle++

		gl.UseProgram(textureProgram)
		gl.UniformMatrix4fv(modelQuadLoc, 1, false, &moq[0])
		bindTexture(quadTbo)
		bindVao(quadVao)
		drawTexture(quadVao, quadTbo, &indicesQuad, window, textureProgram)

		gl.UseProgram(texture2Program)
		gl.UniformMatrix4fv(modelQuadLoc2, 1, false, &moq2[0])
		bindTexture(quadTbo2)
		bindVao(quadVao2)
		drawTexture(quadVao2, quadTbo2, &indicesQuad2, window, texture2Program)

		glfw.PollEvents()
		window.SwapBuffers()

		time.Sleep(time.Millisecond)
	}
}

func changeShape(name string, selectedVbo uint32) {
	coords := textureAtlas.TextureCoords(name)
	if coords == nil {
		panic("Sub texture not found")
	}

	i := 3
	quad[i] = coords[0].s
	quad[i+1] = coords[0].t
	i += 3 + 2
	quad[i] = coords[1].s
	quad[i+1] = coords[1].t
	i += 3 + 2
	quad[i] = coords[2].s
	quad[i+1] = coords[2].t
	i += 3 + 2
	quad[i] = coords[3].s
	quad[i+1] = coords[3].t

	updateTextureVbo(quad, selectedVbo)
}

func drawShape(vao uint32, indices []uint32, window *glfw.Window, program uint32) {
	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

func drawTexture(vao, tbo uint32, indices *[]uint32, window *glfw.Window, program uint32) {
	gl.BindTexture(gl.TEXTURE_2D, tbo)
	gl.DrawElements(gl.TRIANGLES, int32(len(*indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
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
	makeEbo(&indicesTri)

	gl.BindVertexArray(0) // close scope
	// --------- Scope capturing ENDs here -------------------

	return vao
}

var (
	quad  = []float32{}
	quad2 = []float32{}

	// Indices defined in CCW order
	indicesQuad = []uint32{
		// 0, 1, 3, // first triangle
		// 1, 2, 3, // second triangle
		// OR
		0, 1, 2, // first triangle
		0, 2, 3, // second triangle
	}

	// Indices defined in CCW order
	indicesQuad2 = []uint32{
		0, 1, 2, // first triangle
		0, 2, 3, // second triangle
	}
)

func buildQuadTexture(textureAtlas *TextureAtlas, subTexture string, indicesQ *[]uint32, quadFlt *[]float32) (vao, tbo, vbo uint32) {
	vao = makeVao()

	// Activate VBO buffer while in the VAOs scope
	vbo = makeVbo()

	// --------- Scope capturing STARTs here -------------------
	bindVao(vao)

	//       vector-space
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
	//  ^      *-----------*
	//  |      |           |
	//  |+Y    |     _     |
	//  |      |           |
	//  |      |           |
	//        *-----------*
	//       0,0           1,0
	//         A           B

	textureAtlas.Build()

	coords := textureAtlas.TextureCoords(subTexture)
	if coords == nil {
		panic("Sub texture not found")
	}
	*quadFlt = append(*quadFlt, -0.5, -0.5, 0.0)          // xy = aPos
	*quadFlt = append(*quadFlt, coords[0].s, coords[0].t) // uv = aTexCoord
	*quadFlt = append(*quadFlt, 0.5, -0.5, 0.0)
	*quadFlt = append(*quadFlt, coords[1].s, coords[1].t) // uv
	*quadFlt = append(*quadFlt, 0.5, 0.5, 0.0)
	*quadFlt = append(*quadFlt, coords[2].s, coords[2].t) // uv
	*quadFlt = append(*quadFlt, -0.5, 0.5, 0.0)
	*quadFlt = append(*quadFlt, coords[3].s, coords[3].t) // uv

	bindTextureVbo(quad, vbo)

	// Activate EBO buffer while in the VAOs scope
	makeEbo(indicesQ)

	tbo = makeTbo()

	atlas := textureAtlas.Atlas()

	bindTbo(atlas, tbo)

	gl.BindVertexArray(0) // close scope
	// --------- Scope capturing ENDs here -------------------

	return vao, tbo, vbo
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
