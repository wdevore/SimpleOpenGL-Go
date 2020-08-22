package api

// IMatrix4 is a 4x4 matrix
type IMatrix4 interface {
	Translate(v IVector3)
	TranslateBy3Comps(x, y, z float32)
	TranslateBy2Comps(x, y float32)

	SetTranslateUsingVector(v IVector3)
	SetTranslate3Comp(x, y, z float32)
	PostTranslate(tx, ty, tz float32)
	GetTranslation(out IVector3)

	SetRotation(angle float64)
	Rotate(angle float64)

	Scale(v IVector3)
	SetScale(v IVector3)
	ScaleByComp(x, y, z float32)
	SetScale3Comp(sx, sy, sz float32)
	SetScale2Comp(sx, sy float32)
	GetPsuedoScale() float32

	Set(src IMatrix4)
	SetFromAffine(src IAffineTransform)

	Multiply(a, b IMatrix4)
	PreMultiply(b IMatrix4)
	PostMultiply(b IMatrix4)

	C(i int) float32
	Matrix() *([16]float32)

	Clone() IMatrix4
	Eq(IMatrix4) bool
	ToIdentity()
	Invert() bool

	// Graphics
	SetToOrtho(left, right, bottom, top, near, far float32)
}
