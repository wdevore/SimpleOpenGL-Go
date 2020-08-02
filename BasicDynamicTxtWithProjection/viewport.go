// Package graphics provides visual
package main

import "github.com/go-gl/gl/v4.5-core/gl"

// Viewport is a basic wrapper of an OpenGL viewport
type Viewport struct {
	x, y, width, height int32
}

// NewViewport construct a viewport
func NewViewport() *Viewport {
	v := new(Viewport)
	return v
}

// SetDimensions set viewport dimensions
func (v *Viewport) SetDimensions(x, y, width, height int) {
	v.x = int32(x)
	v.y = int32(y)
	v.width = int32(width)
	v.height = int32(height)
}

// Apply set the actual OpenGL viewport
func (v *Viewport) Apply() {
	gl.Viewport(v.x, v.y, v.width, v.height)
}
