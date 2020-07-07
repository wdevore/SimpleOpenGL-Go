package main

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Simple", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func keyCallback(glfwW *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// fmt.Println("key pressed ", key)

	if action == glfw.Press {
		switch key {
		case glfw.KeyEscape:
			quitTriggered = true
		case glfw.KeyM:
			if !polygonMode {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
			} else {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
			}
			polygonMode = !polygonMode
		case glfw.KeyP:
			if !pointMode {
				gl.PointSize(5)
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.POINT)
			} else {
				gl.PointSize(1)
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
			}
			pointMode = !pointMode
		}
	}
}
