package main

// IPoint represents 2D points
type IPoint interface {
	// Components returns the x,y parts
	Components() (float32, float32)

	// ComponentsAsInt32 return x,y parts for render context
	ComponentsAsInt32() (int32, int32)

	// X sets the x component
	X() float32
	// Y sets the y component
	Y() float32
	// SetByComp sets by component
	SetByComp(x, y float32)
	// SetByPoint sets point using another point
	SetByPoint(IPoint)

	MulPoint(IMatrix4)
}
