package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

const (
	vertexShaderSourcePrj = `
    #version 450
    layout (location = 0) in vec3 vp;

    uniform mat4 model;
    // These uniforms don't change and are set once at the start of the client App
    uniform mat4 view;
    uniform mat4 projection;

    void main() {
        gl_Position = projection * view * model * vec4(vp, 1.0);
    }
` + "\x00"

	// -------------------------------------------------------------------
	// Fragments
	// -------------------------------------------------------------------
	fragmentShaderSource = `
    #version 450
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(1, 0.5, 0.0, 1);
    }
` + "\x00"
)

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
