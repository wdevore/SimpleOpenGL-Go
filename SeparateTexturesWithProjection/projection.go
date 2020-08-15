package main

// Projection provides an orthographic projection
type Projection struct {
	near, far                float32
	left, right, bottom, top float32
	Width, Height            float32
	ratioCorrection          float32

	// Projection matrix (orthographic)
	matrix IMatrix4
}

// NewCamera construct a Camera
func NewCamera() *Projection {
	c := new(Projection)
	c.matrix = NewMatrix4()
	return c
}

// Matrix returns internal 4x4 matrix
func (c *Projection) Matrix() IMatrix4 {
	return c.matrix
}

// SetProjection sets orthographic frustum dimensions
func (c *Projection) SetProjection(bottom, left, top, right, near, far float32) {

	c.near = near
	c.far = far
	c.bottom = bottom
	c.left = left
	c.top = top
	c.right = right
	c.Width = right - left
	c.Height = top - bottom

	//                  left, right, bottom, top, near, far
	c.matrix.SetToOrtho(c.left, c.right, c.bottom, c.top, near, far)
}

// SetCenteredProjection centers the projection and adjusts for aspect ratio.
func (c *Projection) SetCenteredProjection(near, far float32) {
	// Adjust for aspect ratio
	left := -c.Width / 2.0 / c.ratioCorrection
	right := c.Width / 2.0 / c.ratioCorrection
	bottom := -c.Height / 2.0 / c.ratioCorrection
	top := c.Height / 2.0 / c.ratioCorrection

	c.matrix.SetToOrtho(left, right, bottom, top, near, far)
}
