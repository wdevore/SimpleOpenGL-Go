package render

import (
	"SimpleOpenGL-Go/SeparateTexturesWithProjection/api"
	"SimpleOpenGL-Go/SeparateTexturesWithProjection/display"
	"SimpleOpenGL-Go/SeparateTexturesWithProjection/maths"
	"log"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type TriangleRender struct {
	vao, tbo, vbo, ebo uint32

	shaderProgram uint32

	projLoc, viewLoc, modelLoc int32

	modelM api.IMatrix4

	triangle []float32
	indices  []uint32
}

func NewTriangleRender() *TriangleRender {
	o := new(TriangleRender)
	o.modelM = maths.NewMatrix4()
	o.modelM.ScaleByComp(25.0, 25.0, 1.0)

	return o
}

func (t *TriangleRender) Build(name string) {
	gl.GenVertexArrays(1, &t.vao)

	gl.GenBuffers(1, &t.vbo)

	// Activate VBO buffer while in the VAOs scope
	gl.BindVertexArray(t.vao)

	t.shaderProgram = t.initShaderProgram()

	// Indices defined in CCW order
	t.indices = []uint32{
		0, 1, 2, // triangle
	}

	t.triangle = []float32{
		-0.2, -0.5, 0.0,
		0.8, -0.5, 0.0,
		0.3, 0.314, 0.0,
	}

	t.bindVbo()

	// Activate EBO buffer while in the VAOs scope
	gl.GenBuffers(1, &t.ebo)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, t.ebo)
	sizeOfUInt32 := 4
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, sizeOfUInt32*len(t.indices), gl.Ptr(t.indices), gl.STATIC_DRAW)

	if errNum := gl.GetError(); errNum != gl.NO_ERROR {
		log.Fatal("(ebo)GL Error: ", errNum)
	}

	gl.BindVertexArray(0) // close scope
	// --------- Scope capturing ENDs here -------------------
}

func (t *TriangleRender) SetAngle(radians float64) {
	t.modelM.SetRotation(radians * display.DegreeToRadians)
	t.modelM.TranslateBy2Comps(100.0, 0.0)
	t.modelM.ScaleByComp(25.0, 25.0, 1.0)
}

func (t *TriangleRender) Draw() {
	gl.UseProgram(t.shaderProgram)

	gl.UniformMatrix4fv(t.modelLoc, 1, false, &t.modelM.Matrix()[0])

	gl.BindVertexArray(t.vao)

	gl.DrawElements(gl.TRIANGLES, int32(len(t.indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

	gl.BindVertexArray(0)
}

func (t *TriangleRender) SetUniforms(proj *display.Projection, view api.IMatrix4) {
	pm := proj.Matrix().Matrix()
	gl.UniformMatrix4fv(t.projLoc, 1, false, &pm[0])

	gl.UniformMatrix4fv(t.viewLoc, 1, false, &view.Matrix()[0])
}

func (t *TriangleRender) bindVbo() {
	gl.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(t.triangle), gl.Ptr(t.triangle), gl.STATIC_DRAW)

	// Specify the vertex attribute layout. This specifies--how during transmission--opengl
	// will pass the data to the shader and how the shader will extract the data.
	vertexInputAttrb := uint32(0)
	sizeOfInputAttrb := int32(3)
	stride := int32(0)
	gl.VertexAttribPointer(vertexInputAttrb, sizeOfInputAttrb, gl.FLOAT, false, stride, nil)
	gl.EnableVertexAttribArray(vertexInputAttrb)
}

func (t *TriangleRender) initShaderProgram() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	vertexShader, err := compileShader(vertexShaderSourcePrj, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	gl.UseProgram(prog)

	t.projLoc = gl.GetUniformLocation(prog, gl.Str("projection\x00"))
	if t.projLoc < 0 {
		panic("TriangleRender: couldn't find 'projection' uniform variable")
	}

	t.viewLoc = gl.GetUniformLocation(prog, gl.Str("view\x00"))
	if t.viewLoc < 0 {
		panic("TriangleRender: couldn't find 'view' uniform variable")
	}

	t.modelLoc = gl.GetUniformLocation(prog, gl.Str("model\x00"))
	if t.modelLoc < 0 {
		panic("TriangleRender: couldn't find 'model' uniform variable")
	}

	return prog
}
