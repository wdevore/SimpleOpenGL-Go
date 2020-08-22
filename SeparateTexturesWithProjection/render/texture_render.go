package render

import (
	"SimpleOpenGL-Go/SeparateTexturesWithProjection/api"
	"SimpleOpenGL-Go/SeparateTexturesWithProjection/display"
	"SimpleOpenGL-Go/SeparateTexturesWithProjection/maths"
	"SimpleOpenGL-Go/SeparateTexturesWithProjection/textures"
	"image"
	"log"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type TextureRender struct {
	vao, tbo, vbo, ebo uint32

	shaderProgram uint32
	textureAtlas  *textures.TextureAtlas

	projLoc, viewLoc, modelLoc int32

	modelM api.IMatrix4

	quad    []float32
	indices []uint32
}

func NewTextureRender(textureAtlas *textures.TextureAtlas) *TextureRender {
	o := new(TextureRender)
	o.modelM = maths.NewMatrix4()
	o.modelM.ScaleByComp(64.0, 64.0, 1.0)

	o.textureAtlas = textureAtlas
	return o
}

func (t *TextureRender) Build(name string) {
	gl.GenVertexArrays(1, &t.vao)

	gl.GenBuffers(1, &t.vbo)

	// Activate VBO buffer while in the VAOs scope
	gl.BindVertexArray(t.vao)

	t.shaderProgram = t.initShaderProgram()

	// Indices defined in CCW order
	t.indices = []uint32{
		// 0, 1, 3, // first triangle
		// 1, 2, 3, // second triangle
		// OR
		0, 1, 2, // first triangle
		0, 2, 3, // second triangle
	}

	coords := t.textureAtlas.TextureCoords(name)
	if coords == nil {
		panic("Sub texture not found")
	}

	t.quad = append(t.quad, -0.5, -0.5, 0.0)          // xy = aPos
	t.quad = append(t.quad, coords[0].S, coords[0].T) // uv = aTexCoord
	t.quad = append(t.quad, 0.5, -0.5, 0.0)
	t.quad = append(t.quad, coords[1].S, coords[1].T) // uv
	t.quad = append(t.quad, 0.5, 0.5, 0.0)
	t.quad = append(t.quad, coords[2].S, coords[2].T) // uv
	t.quad = append(t.quad, -0.5, 0.5, 0.0)
	t.quad = append(t.quad, coords[3].S, coords[3].T) // uv

	t.bindTextureVbo()

	// Activate EBO buffer while in the VAOs scope
	gl.GenBuffers(1, &t.ebo)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, t.ebo)
	sizeOfUInt32 := 4
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, sizeOfUInt32*len(t.indices), gl.Ptr(t.indices), gl.STATIC_DRAW)

	if errNum := gl.GetError(); errNum != gl.NO_ERROR {
		log.Fatal("(ebo)GL Error: ", errNum)
	}

	gl.GenTextures(1, &t.tbo)

	t.bindTbo(t.textureAtlas.Atlas())

	gl.BindVertexArray(0) // close scope
	// --------- Scope capturing ENDs here -------------------
}

func (t *TextureRender) SetPosition(x, y float32) {
	t.modelM.SetTranslate3Comp(x, y, 0.0)
	t.modelM.ScaleByComp(64.0, 64.0, 1.0)
}

func (t *TextureRender) Draw() {
	gl.UseProgram(t.shaderProgram)

	gl.UniformMatrix4fv(t.modelLoc, 1, false, &t.modelM.Matrix()[0])

	gl.BindVertexArray(t.vao)

	gl.ActiveTexture(gl.TEXTURE0)

	gl.BindTexture(gl.TEXTURE_2D, t.tbo)
	gl.DrawElements(gl.TRIANGLES, int32(len(t.indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

	gl.BindVertexArray(0)
}

func (t *TextureRender) ChangeShape(name string) {
	coords := t.textureAtlas.TextureCoords(name)
	if coords == nil {
		panic("Sub texture not found")
	}

	i := 3
	t.quad[i] = coords[0].S
	t.quad[i+1] = coords[0].T
	i += 3 + 2
	t.quad[i] = coords[1].S
	t.quad[i+1] = coords[1].T
	i += 3 + 2
	t.quad[i] = coords[2].S
	t.quad[i+1] = coords[2].T
	i += 3 + 2
	t.quad[i] = coords[3].S
	t.quad[i+1] = coords[3].T

	t.updateTextureVbo()
}

func (t *TextureRender) SetUniforms(proj *display.Projection, view api.IMatrix4) {
	pm := proj.Matrix().Matrix()
	gl.UniformMatrix4fv(t.projLoc, 1, false, &pm[0])

	gl.UniformMatrix4fv(t.viewLoc, 1, false, &view.Matrix()[0])
}

// Update moves any modified data to the buffer.
func (t *TextureRender) updateTextureVbo() {
	gl.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(t.quad)*4, gl.Ptr(t.quad))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (t *TextureRender) initShaderProgram() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	vertexShader, err := compileShader(vertexTextureShaderSourcePrj, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentTextureShaderSource, gl.FRAGMENT_SHADER)
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
		panic("TextureRender: couldn't find 'projection' uniform variable")
	}

	t.viewLoc = gl.GetUniformLocation(prog, gl.Str("view\x00"))
	if t.viewLoc < 0 {
		panic("TextureRender: couldn't find 'view' uniform variable")
	}

	t.modelLoc = gl.GetUniformLocation(prog, gl.Str("model\x00"))
	if t.modelLoc < 0 {
		panic("TextureRender: couldn't find 'model' uniform variable")
	}

	return prog
}

func (t *TextureRender) bindTextureVbo() {
	gl.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(t.quad), gl.Ptr(t.quad), gl.DYNAMIC_DRAW)

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

func (t *TextureRender) bindTbo(texture *image.NRGBA) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.tbo)

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	width := int32(texture.Bounds().Dx())
	height := int32(texture.Bounds().Dy())

	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	// Give the image to OpenGL
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(texture.Pix))
	// gl.GenerateMipmap(gl.TEXTURE_2D)
}
