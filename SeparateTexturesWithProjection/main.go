package main

import (
	"SimpleOpenGL-Go/SeparateTexturesWithProjection/api"
	"SimpleOpenGL-Go/SeparateTexturesWithProjection/display"
	"SimpleOpenGL-Go/SeparateTexturesWithProjection/maths"
	"SimpleOpenGL-Go/SeparateTexturesWithProjection/render"
	"SimpleOpenGL-Go/SeparateTexturesWithProjection/textures"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	textureAtlas        = textures.NewTextureAtlas("assets/texture_manifest.txt")
	texture2Atlas       = textures.NewTextureAtlas("assets/texture2_manifest.txt")
	textureRender       *render.TextureRender
	texture2Render      *render.TextureRender
	activeTextureRender *render.TextureRender
)

func main() {
	runtime.LockOSThread()
	display.PolygonMode = false

	window := display.InitGlfw()
	defer glfw.Terminate()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	window.SetKeyCallback(KeyCallback)

	// -----------------------------------------------------------
	viewport := display.NewViewport()
	viewport.SetDimensions(0, 0, display.Width, display.Height)
	viewport.Apply()

	projection := buildProjection()

	view := buildView()

	textureAtlas.Build()
	texture2Atlas.Build()

	textureRender = render.NewTextureRender(textureAtlas)
	textureRender.Build("orange ship")
	textureRender.SetUniforms(projection, view)
	activeTextureRender = textureRender
	textureRender.SetPosition(-200.0, 0.0)

	texture2Render = render.NewTextureRender(texture2Atlas)
	texture2Render.Build("green ship")
	texture2Render.SetUniforms(projection, view)
	texture2Render.SetPosition(200.0, 0.0)

	triangleRender := render.NewTriangleRender()
	triangleRender.Build("Triangle")
	triangleRender.SetUniforms(projection, view)
	triangleRender.SetAngle(0.0)

	// -----------------------------------------------------------
	angle := 0.0

	gl.ClearColor(0.25, 0.25, 0.25, 1.0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	for !window.ShouldClose() && !display.QuitTriggered {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		triangleRender.Draw()
		angle++
		triangleRender.SetAngle(angle)

		textureRender.Draw()
		texture2Render.Draw()

		glfw.PollEvents()
		window.SwapBuffers()

		time.Sleep(time.Millisecond)
	}
}

func buildProjection() *display.Projection {
	projection := display.NewCamera()

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
		float32(display.Height), float32(display.Width), //top,right
		-1.0, 1.0) //0.1, 100.0

	return projection
}

func buildView() api.IMatrix4 {
	centered := true
	offsetX := float32(0.0)
	offsetY := float32(0.0)

	if centered {
		offsetX = float32(display.Width) / 2.0
		offsetY = float32(display.Height) / 2.0
	}

	view := maths.NewMatrix4()
	view.SetTranslate3Comp(offsetX, offsetY, 0.5)

	return view
}

func KeyCallback(glfwW *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// fmt.Println("key pressed ", key)

	if action == glfw.Press {
		switch key {
		case glfw.KeyEscape:
			display.QuitTriggered = true
		case glfw.KeyM:
			if !display.PolygonMode {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
			} else {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
			}
			display.PolygonMode = !display.PolygonMode
		case glfw.KeyP:
			if !display.PointMode {
				gl.PointSize(5)
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.POINT)
			} else {
				gl.PointSize(1)
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
			}
			display.PointMode = !display.PointMode
		case glfw.Key0:
			if activeTextureRender == textureRender {
				activeTextureRender = texture2Render
			} else {
				activeTextureRender = textureRender
			}

		case glfw.Key1:
			fmt.Println("mine")
			activeTextureRender.ChangeShape("mine")
		case glfw.Key2:
			fmt.Println("green ship")
			activeTextureRender.ChangeShape("green ship")
		case glfw.Key3:
			fmt.Println("orange ship")
			activeTextureRender.ChangeShape("orange ship")
		case glfw.Key4:
			fmt.Println("ctype ship")
			activeTextureRender.ChangeShape("ctype ship")
		case glfw.Key5:
			fmt.Println("bomb")
			activeTextureRender.ChangeShape("bomb")
		}
	}
}
