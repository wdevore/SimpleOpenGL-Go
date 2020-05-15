package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"
)

// makeVao initializes and returns a vertex array from the points provided.
func makeVao() uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	fmt.Println("Vao ID: ", vao)
	return vao
}

func bindVao(vao uint32) {
	gl.BindVertexArray(vao)
}
