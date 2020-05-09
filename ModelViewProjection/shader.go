package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

// const (
// 	vertexShaderSource = `
//     #version 450
//     in vec3 vp;
//     void main() {
//         gl_Position = vec4(vp, 1.0);
//     }
// ` + "\x00"

// 	fragmentShaderSource = `
//     #version 450
//     out vec4 frag_colour;
//     void main() {
//         frag_colour = vec4(1, 0.5, 0.0, 1);
//     }
// ` + "\x00"
// )

const (
	vertexShaderSource = `
	#version 450 core

	// This input attribute handles the stream of vertices sent to this shader
	layout (location = 0) in vec3 aPos;

	uniform mat4 model;
	// This uniforms are set once at the start of the client App
	uniform mat4 view;
	uniform mat4 projection;
	
	void main()
	{
		// note that we read the multiplication from right to left
		gl_Position = projection * view * model * vec4(aPos, 1.0);
	}
` + "\x00"

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
