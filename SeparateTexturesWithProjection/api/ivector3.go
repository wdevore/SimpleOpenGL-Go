package api

// IVector3 is a 3x1 vector
type IVector3 interface {
	Clone() IVector3
	Set3Components(x, y, z float32)
	Set2Components(x, y float32)
	Components2D() (x, y float32)
	Components3D() (x, y, z float32)
	X() float32
	Y() float32
	Z() float32
	Set(source IVector3)
	Add(src IVector3)
	Add2Components(x, y float32)
	Sub(src IVector3)
	Sub2Components(x, y float32)
	ScaleBy(s float32)
	ScaleBy2Components(sx, sy float32)
	MulAdd(src IVector3, scalar float32)
	Length() float32
	LengthSquared() float32
	Equal(other IVector3) bool
	EqEpsilon(other IVector3) bool
	Distance(src IVector3) float32
	DistanceSquared(src IVector3) float32
	DotByComponent(x, y, z float32) float32
	Dot(o IVector3) float32
	Cross(o IVector3)
	Mul(m IMatrix4)
}
