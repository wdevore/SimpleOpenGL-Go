package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

const (
	vertexShaderSource = `
    #version 450
    in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 450
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(1, 0.5, 0.0, 1);
    }
` + "\x00"

	vertexTextureShaderSource = `
    #version 450
    layout (location = 0) in vec3 aPos;
    layout (location = 1) in vec3 aColor;
    layout (location = 2) in vec2 aTexCoord;
    
    //out vec3 ourColor;
    out vec2 TexCoord;
    
    void main() {
        gl_Position = vec4(aPos, 1.0);
        //ourColor = aColor;
        TexCoord = vec2(aTexCoord.xy);
    }
` + "\x00"

	fragmentTextureShaderSource = `
    #version 450
    out vec4 FragColor;
    
    //in vec3 ourColor;
    in vec2 TexCoord;
    
    // texture sampler
    uniform sampler2D texture1;
    
    void main()
    {
        // Basic binary transparency with blending disabled
        // vec4 texColor = texture(texture1, TexCoord);
        // if (texColor.a < 0.75)
        //     discard;
        // FragColor = texColor;

        FragColor = texture(texture1, TexCoord);
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
