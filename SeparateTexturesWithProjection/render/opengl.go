package render

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

// initDefaultProgram initializes OpenGL and returns an intiialized program.
func InitDefaultProgram() uint32 {
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

	return prog
}
