package api

// IAffineTransform represents 2D transforms
type IAffineTransform interface {
	Matrix() *([16]float32)

	Components() (float32, float32, float32, float32, float32, float32)
	// ToIdentity sets the transform to an identity matrix
	ToIdentity()

	// --------------------------------------------
	// Setters
	// --------------------------------------------

	// SetByComp sets by component
	SetByComp(float32, float32, float32, float32, float32, float32)
	// SetByTransform sets point using another transform
	SetByTransform(IAffineTransform)

	// --------------------------------------------
	// Transforms
	// --------------------------------------------

	// TransformPoint applys affine transform to point
	TransformPoint(IPoint)
	// TransformToPoint applys affine transform to out point, "in" is not modified
	TransformToPoint(in IPoint, out IPoint)
	// TransformToComps applys transform and returns results, "in" is not modified
	TransformToComps(in IPoint) (x, y float32)
	TransformCompToPoint(x, y float32, out IPoint)
	// TransformPolygon

	// --------------------------------------------
	// Mutaters
	// --------------------------------------------

	// MakeTranslate sets the transform to a Translate matrix
	MakeTranslate(x, y float32)
	// MakeTranslateUsingPoint sets the transform to a Translate matrix
	MakeTranslateUsingPoint(p IPoint)

	// Translate mutates/concat "this" matrix using tx,ty
	Translate(tx, ty float32)

	// MakeScale sets the transform to a Scale matrix
	MakeScale(x, y float32)
	// Scale mutates "this" matrix using sx, sy
	Scale(sx, sy float32)

	// GetPsuedoScale returns the transform's "a" component, however,
	// this is only valid if the transform doesn't have a rotation or zoom applied.
	GetPsuedoScale() float32

	// MakeRotate sets the transform to a Rotate matrix
	MakeRotate(radians float64)
	// Rotate mutates "this" matrix using radian angle
	Rotate(radians float64)

	// --------------------------------------------
	// Inversions
	// --------------------------------------------

	// Invert (mutates) inverts "this" matrix
	Invert()
	// Invert (non-mutating) inverts "this" matrix and sends to "out"
	InvertTo(out IAffineTransform)

	// Transpose
	// Converts either from or to pre or post multiplication.
	//     a c
	//     b d
	// to
	//     a b
	//     c d
	Transpose()

	Populate(destination IMatrix4)

	String4x4() string
}
