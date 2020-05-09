package main

import (
	"fmt"
	"math"
)

//    2x3          3x3            4x4         OpenGL Array Index      M(Cell/Row)
// | a c e |    | a c e |      | a c 0 e |      |00 04 08 12|     |M00 M01 M02 M03|
// | b d f |    | b d f |  =>  | b d 0 f | ==>  |01 05 09 13| ==> |M10 M11 M12 M13|
// 	            | 0 0 1 |      | 0 0 1 0 |      |02 06 10 14|     |M20 M21 M22 M23|
// 	                           | 0 0 0 1 |      |03 07 11 15|     |M30 M31 M32 M33|

// Array indices
const (
	// M00 XX: Typically the unrotated X component for scaling, also the cosine of the
	// angle when rotated on the Y and/or Z axis. On
	// Vector3 multiplication this value is multiplied with the source X component
	// and added to the target X component.
	M00 = 0
	// M01 XY: Typically the negative sine of the angle when rotated on the Z axis.
	// On Vector3 multiplication this value is multiplied
	// with the source Y component and added to the target X component.
	M01 = 4
	// M02 XZ: Typically the sine of the angle when rotated on the Y axis.
	// On Vector3 multiplication this value is multiplied with the
	// source Z component and added to the target X component.
	M02 = 8
	// M03 XW: Typically the translation of the X component.
	// On Vector3 multiplication this value is added to the target X component.
	M03 = 12

	// M10 YX: Typically the sine of the angle when rotated on the Z axis.
	// On Vector3 multiplication this value is multiplied with the
	// source X component and added to the target Y component.
	M10 = 1
	// M11 YY: Typically the unrotated Y component for scaling, also the cosine
	// of the angle when rotated on the X and/or Z axis. On
	// Vector3 multiplication this value is multiplied with the source Y
	// component and added to the target Y component.
	M11 = 5
	// M12 YZ: Typically the negative sine of the angle when rotated on the X axis.
	// On Vector3 multiplication this value is multiplied
	// with the source Z component and added to the target Y component.
	M12 = 9
	// M13 YW: Typically the translation of the Y component.
	// On Vector3 multiplication this value is added to the target Y component.
	M13 = 13

	// M20 ZX: Typically the negative sine of the angle when rotated on the Y axis.
	// On Vector3 multiplication this value is multiplied
	// with the source X component and added to the target Z component.
	M20 = 2
	// M21 ZY: Typical the sine of the angle when rotated on the X axis.
	// On Vector3 multiplication this value is multiplied with the
	// source Y component and added to the target Z component.
	M21 = 6
	// M22 ZZ: Typically the unrotated Z component for scaling, also the cosine of the
	// angle when rotated on the X and/or Y axis.
	// On Vector3 multiplication this value is multiplied with the source Z component
	// and added to the target Z component.
	M22 = 10
	// M23 ZW: Typically the translation of the Z component.
	// On Vector3 multiplication this value is added to the target Z component.
	M23 = 14

	// M30 WX: Typically the value zero. On Vector3 multiplication this value is ignored.
	M30 = 3
	// M31 WY: Typically the value zero. On Vector3 multiplication this value is ignored.
	M31 = 7
	// M32 WZ: Typically the value zero. On Vector3 multiplication this value is ignored.
	M32 = 11
	// M33 WW: Typically the value one. On Vector3 multiplication this value is ignored.
	M33 = 15
)

// matrix4 represents a column major opengl array.
type matrix4 struct {
	e [16]float32
}

// A temporary matrix for multiplication
var tempM0 = NewMatrix4()
var mulM = NewMatrix4()

// NewMatrix4 creates a Matrix4 initialized to an identity matrix
func NewMatrix4() *matrix4 {
	m := new(matrix4)
	m.ToIdentity()
	return m
}

// E returns the internal 4x4 matrix
func (m *matrix4) Matrix() *([16]float32) {
	return &m.e
}

// --------------------------------------------------------------------------
// Translation
// --------------------------------------------------------------------------

// TranslateBy2Comps adds a translational component to the matrix in the 4th column.
// Z is unmodified. The other columns are unmodified.
func (m *matrix4) TranslateBy2Comps(x, y float32) {
	m.TranslateBy3Comps(x, y, 0.0)
}

// TranslateBy3Comps adds a translational component to the matrix in the 4th column.
// The other columns are unmodified.
func (m *matrix4) TranslateBy3Comps(x, y, z float32) {
	e := tempM0.Matrix()
	e[M00] = 1.0
	e[M01] = 0.0
	e[M02] = 0.0
	e[M03] = x
	e[M10] = 0.0
	e[M11] = 1.0
	e[M12] = 0.0
	e[M13] = y
	e[M20] = 0.0
	e[M21] = 0.0
	e[M22] = 1.0
	e[M23] = z
	e[M30] = 0.0
	e[M31] = 0.0
	e[M32] = 0.0
	e[M33] = 1.0

	m.PreMultiply(tempM0)
}

// SetTranslate3Comp sets the translational component to the matrix in the 4th column.
// The other columns are unmodified.
func (m *matrix4) SetTranslate3Comp(x, y, z float32) {
	m.ToIdentity()

	m.e[M03] = x
	m.e[M13] = y
	m.e[M23] = z
}

// --------------------------------------------------------------------------
// Rotation
// --------------------------------------------------------------------------

// SetRotation set a rotation matrix about Z axis. 'angle' is specified
// in radians.
//
//      [  M00  M01   _    _   ]
//      [  M10  M11   _    _   ]
//      [   _    _    _    _   ]
//      [   _    _    _    _   ]
func (m *matrix4) SetRotation(angle float64) {
	if angle == 0 {
		return
	}

	m.ToIdentity()

	// Column major
	c := float32(math.Cos(angle))
	s := float32(math.Sin(angle))

	m.e[M00] = c
	m.e[M01] = -s
	m.e[M10] = s
	m.e[M11] = c
}

// Rotate postmultiplies this matrix with a (counter-clockwise) rotation matrix whose
// angle is specified in radians.
// |M00 M01 M02 M03|
// |M10 M11 M12 M13|
// |M20 M21 M22 M23|
// |M30 M31 M32 M33|
func (m *matrix4) Rotate(angle float64) {
	if angle == 0.0 {
		return
	}

	// Column major
	c := float32(math.Cos(angle))
	s := float32(math.Sin(angle))

	e := tempM0.Matrix()
	e[M00] = c
	e[M01] = -s
	e[M02] = 0.0
	e[M03] = 0.0
	e[M10] = s
	e[M11] = c
	e[M12] = 0.0
	e[M13] = 0.0
	e[M20] = 0.0
	e[M21] = 0.0
	e[M22] = 1.0
	e[M23] = 0.0
	e[M30] = 0.0
	e[M31] = 0.0
	e[M32] = 0.0
	e[M33] = 1.0

	m.PreMultiply(tempM0)
}

// --------------------------------------------------------------------------
// Scale
// --------------------------------------------------------------------------

// SetScale3Comp sets the scale components of an identity matrix and captures
// scale values into Scale property.
func (m *matrix4) SetScale3Comp(sx, sy, sz float32) {
	m.ToIdentity()

	m.e[M00] = sx
	m.e[M11] = sy
	m.e[M22] = sz
}

// SetScale2Comp sets the scale components of an identity matrix and captures
// scale values into Scale property where Z component = 1.0.
func (m *matrix4) SetScale2Comp(sx, sy float32) {
	m.ToIdentity()

	m.e[M00] = sx
	m.e[M11] = sy
	m.e[M22] = 1.0
}

// ScaleByComp scales the scale components.
func (m *matrix4) ScaleByComp(sx, sy, sz float32) {
	e := tempM0.Matrix()

	e[M00] = sx
	e[M01] = 0
	e[M02] = 0
	e[M03] = 0
	e[M10] = 0
	e[M11] = sy
	e[M12] = 0
	e[M13] = 0
	e[M20] = 0
	e[M21] = 0
	e[M22] = sz
	e[M23] = 0
	e[M30] = 0
	e[M31] = 0
	e[M32] = 0
	e[M33] = 1

	m.PreMultiply(tempM0)
}

func (m *matrix4) GetPsuedoScale() float32 {
	return m.e[M00]
}

// --------------------------------------------------------------------------
// Transforms
// --------------------------------------------------------------------------

// --------------------------------------------------------------------------
// Matrix methods
// --------------------------------------------------------------------------

func multiply4(a, b, out *([16]float32)) {
	out[M00] = a[M00]*b[M00] + a[M01]*b[M10] + a[M02]*b[M20] + a[M03]*b[M30]
	out[M01] = a[M00]*b[M01] + a[M01]*b[M11] + a[M02]*b[M21] + a[M03]*b[M31]
	out[M02] = a[M00]*b[M02] + a[M01]*b[M12] + a[M02]*b[M22] + a[M03]*b[M32]
	out[M03] = a[M00]*b[M03] + a[M01]*b[M13] + a[M02]*b[M23] + a[M03]*b[M33]
	out[M10] = a[M10]*b[M00] + a[M11]*b[M10] + a[M12]*b[M20] + a[M13]*b[M30]
	out[M11] = a[M10]*b[M01] + a[M11]*b[M11] + a[M12]*b[M21] + a[M13]*b[M31]
	out[M12] = a[M10]*b[M02] + a[M11]*b[M12] + a[M12]*b[M22] + a[M13]*b[M32]
	out[M13] = a[M10]*b[M03] + a[M11]*b[M13] + a[M12]*b[M23] + a[M13]*b[M33]
	out[M20] = a[M20]*b[M00] + a[M21]*b[M10] + a[M22]*b[M20] + a[M23]*b[M30]
	out[M21] = a[M20]*b[M01] + a[M21]*b[M11] + a[M22]*b[M21] + a[M23]*b[M31]
	out[M22] = a[M20]*b[M02] + a[M21]*b[M12] + a[M22]*b[M22] + a[M23]*b[M32]
	out[M23] = a[M20]*b[M03] + a[M21]*b[M13] + a[M22]*b[M23] + a[M23]*b[M33]
	out[M30] = a[M30]*b[M00] + a[M31]*b[M10] + a[M32]*b[M20] + a[M33]*b[M30]
	out[M31] = a[M30]*b[M01] + a[M31]*b[M11] + a[M32]*b[M21] + a[M33]*b[M31]
	out[M32] = a[M30]*b[M02] + a[M31]*b[M12] + a[M32]*b[M22] + a[M33]*b[M32]
	out[M33] = a[M30]*b[M03] + a[M31]*b[M13] + a[M32]*b[M23] + a[M33]*b[M33]
}

// Multiply4 multiplies a * b and places result into 'out', (i.e. out = a * b)
func Multiply4(a, b, out *matrix4) {
	oe := out.Matrix()
	ae := a.Matrix()
	be := b.Matrix()
	multiply4(ae, be, oe)
}

// Multiply multiplies a * b and places result into this matrix, (i.e. m = a * b)
func (m *matrix4) Multiply(a, b *matrix4) {
	Multiply4(a, b, m)
}

// PreMultiply pre-multiplies 'b' matrix with 'm' and places the result into 'm' matrix.
// (i.e. m = m * b)
func (m *matrix4) PreMultiply(b *matrix4) {
	Multiply4(m, b, mulM)
	m.Set(mulM)
}

// PostMultiply post-multiplies 'm' matrix with 'b' and places the result into 'm' matrix.
// (i.e. m = b * m)
func (m *matrix4) PostMultiply(b *matrix4) {
	Multiply4(b, m, mulM)
	m.Set(mulM)
}

// MultiplyIntoA multiplies a * b and places result into 'a', (i.e. a = a * b)
func MultiplyIntoA(a, b *matrix4) {
	Multiply4(a, b, mulM)
	a.Set(mulM)
}

// PostTranslate postmultiplies this matrix by a translation matrix.
// Postmultiplication is also used by OpenGL ES.
func (m *matrix4) PostTranslate(tx, ty, tz float32) {
	te := tempM0.Matrix()

	te[M00] = 1.0
	te[M01] = 0.0
	te[M02] = 0.0
	te[M03] = tx
	te[M10] = 0.0
	te[M11] = 1.0
	te[M12] = 0.0
	te[M13] = ty
	te[M20] = 0.0
	te[M21] = 0.0
	te[M22] = 1.0
	te[M23] = tz
	te[M30] = 0.0
	te[M31] = 0.0
	te[M32] = 0.0
	te[M33] = 1.0

	m.PostMultiply(tempM0)
}

// --------------------------------------------------------------------------
// Projections
// --------------------------------------------------------------------------

// SetToOrtho sets the matrix for a 2d ortho graphic projection
//  |M00 M01 M02 M03|
//  |M10 M11 M12 M13|
//  |M20 M21 M22 M23|
//  |M30 M31 M32 M33|
//
//  [M00,M10,M20,M30,M01,M11,M21,M31,M02,M12,M22,M32,M03,M13,M23,M33]
func (m *matrix4) SetToOrtho(left, right, bottom, top, near, far float32) {
	xorth := 2.0 / (right - left)
	yorth := 2.0 / (top - bottom)
	zorth := 2.0 / (near - far)

	tx := (right + left) / (left - right)
	ty := (top + bottom) / (bottom - top)
	tz := (far + near) / (far - near)

	m.e[M00] = xorth
	m.e[M10] = 0.0
	m.e[M20] = 0.0
	m.e[M30] = 0.0
	m.e[M01] = 0.0
	m.e[M11] = yorth
	m.e[M21] = 0.0
	m.e[M31] = 0.0
	m.e[M02] = 0.0
	m.e[M12] = 0.0
	m.e[M22] = zorth
	m.e[M32] = 0.0
	m.e[M03] = tx
	m.e[M13] = ty
	m.e[M23] = tz
	m.e[M33] = 1.0
}

// --------------------------------------------------------------------------
// Misc
// --------------------------------------------------------------------------

// Eq does an epsilon compare
func (m *matrix4) Eq(other *matrix4) bool {
	eq := true
	o := other.Matrix()

	eq = eq && (m.e[M00] == o[M00])
	eq = eq && (m.e[M01] == o[M01])
	eq = eq && (m.e[M02] == o[M02])
	eq = eq && (m.e[M03] == o[M03])

	eq = eq && (m.e[M10] == o[M10])
	eq = eq && (m.e[M11] == o[M11])
	eq = eq && (m.e[M12] == o[M12])
	eq = eq && (m.e[M13] == o[M13])

	eq = eq && (m.e[M20] == o[M20])
	eq = eq && (m.e[M21] == o[M21])
	eq = eq && (m.e[M22] == o[M22])
	eq = eq && (m.e[M23] == o[M23])

	eq = eq && (m.e[M30] == o[M30])
	eq = eq && (m.e[M31] == o[M31])
	eq = eq && (m.e[M32] == o[M32])
	eq = eq && (m.e[M33] == o[M33])

	return eq
}

// C returns a cell value based on Mxx index
func (m *matrix4) C(i int) float32 {
	return m.e[i]
}

// Clone returns a clone of this matrix
func (m *matrix4) Clone() *matrix4 {
	c := new(matrix4)
	c.Set(m)
	return c
}

// Set copies src into this matrix
func (m *matrix4) Set(src *matrix4) {
	se := src.Matrix()

	m.e[M00] = se[M00]
	m.e[M01] = se[M01]
	m.e[M02] = se[M02]
	m.e[M03] = se[M03]

	m.e[M10] = se[M10]
	m.e[M11] = se[M11]
	m.e[M12] = se[M12]
	m.e[M13] = se[M13]

	m.e[M20] = se[M20]
	m.e[M21] = se[M21]
	m.e[M22] = se[M22]
	m.e[M23] = se[M23]

	m.e[M30] = se[M30]
	m.e[M31] = se[M31]
	m.e[M32] = se[M32]
	m.e[M33] = se[M33]
}

// ToIdentity set this matrix to the identity matrix
func (m *matrix4) ToIdentity() {
	m.e[M00] = 1.0
	m.e[M01] = 0.0
	m.e[M02] = 0.0
	m.e[M03] = 0.0

	m.e[M10] = 0.0
	m.e[M11] = 1.0
	m.e[M12] = 0.0
	m.e[M13] = 0.0

	m.e[M20] = 0.0
	m.e[M21] = 0.0
	m.e[M22] = 1.0
	m.e[M23] = 0.0

	m.e[M30] = 0.0
	m.e[M31] = 0.0
	m.e[M32] = 0.0
	m.e[M33] = 1.0
}

func (m *matrix4) Inverse() bool {
	inv := tempM0.Matrix()

	inv[0] = m.e[5]*m.e[10]*m.e[15] -
		m.e[5]*m.e[11]*m.e[14] -
		m.e[9]*m.e[6]*m.e[15] +
		m.e[9]*m.e[7]*m.e[14] +
		m.e[13]*m.e[6]*m.e[11] -
		m.e[13]*m.e[7]*m.e[10]

	inv[4] = -m.e[4]*m.e[10]*m.e[15] +
		m.e[4]*m.e[11]*m.e[14] +
		m.e[8]*m.e[6]*m.e[15] -
		m.e[8]*m.e[7]*m.e[14] -
		m.e[12]*m.e[6]*m.e[11] +
		m.e[12]*m.e[7]*m.e[10]

	inv[8] = m.e[4]*m.e[9]*m.e[15] -
		m.e[4]*m.e[11]*m.e[13] -
		m.e[8]*m.e[5]*m.e[15] +
		m.e[8]*m.e[7]*m.e[13] +
		m.e[12]*m.e[5]*m.e[11] -
		m.e[12]*m.e[7]*m.e[9]

	inv[12] = -m.e[4]*m.e[9]*m.e[14] +
		m.e[4]*m.e[10]*m.e[13] +
		m.e[8]*m.e[5]*m.e[14] -
		m.e[8]*m.e[6]*m.e[13] -
		m.e[12]*m.e[5]*m.e[10] +
		m.e[12]*m.e[6]*m.e[9]

	inv[1] = -m.e[1]*m.e[10]*m.e[15] +
		m.e[1]*m.e[11]*m.e[14] +
		m.e[9]*m.e[2]*m.e[15] -
		m.e[9]*m.e[3]*m.e[14] -
		m.e[13]*m.e[2]*m.e[11] +
		m.e[13]*m.e[3]*m.e[10]

	inv[5] = m.e[0]*m.e[10]*m.e[15] -
		m.e[0]*m.e[11]*m.e[14] -
		m.e[8]*m.e[2]*m.e[15] +
		m.e[8]*m.e[3]*m.e[14] +
		m.e[12]*m.e[2]*m.e[11] -
		m.e[12]*m.e[3]*m.e[10]

	inv[9] = -m.e[0]*m.e[9]*m.e[15] +
		m.e[0]*m.e[11]*m.e[13] +
		m.e[8]*m.e[1]*m.e[15] -
		m.e[8]*m.e[3]*m.e[13] -
		m.e[12]*m.e[1]*m.e[11] +
		m.e[12]*m.e[3]*m.e[9]

	inv[13] = m.e[0]*m.e[9]*m.e[14] -
		m.e[0]*m.e[10]*m.e[13] -
		m.e[8]*m.e[1]*m.e[14] +
		m.e[8]*m.e[2]*m.e[13] +
		m.e[12]*m.e[1]*m.e[10] -
		m.e[12]*m.e[2]*m.e[9]

	inv[2] = m.e[1]*m.e[6]*m.e[15] -
		m.e[1]*m.e[7]*m.e[14] -
		m.e[5]*m.e[2]*m.e[15] +
		m.e[5]*m.e[3]*m.e[14] +
		m.e[13]*m.e[2]*m.e[7] -
		m.e[13]*m.e[3]*m.e[6]

	inv[6] = -m.e[0]*m.e[6]*m.e[15] +
		m.e[0]*m.e[7]*m.e[14] +
		m.e[4]*m.e[2]*m.e[15] -
		m.e[4]*m.e[3]*m.e[14] -
		m.e[12]*m.e[2]*m.e[7] +
		m.e[12]*m.e[3]*m.e[6]

	inv[10] = m.e[0]*m.e[5]*m.e[15] -
		m.e[0]*m.e[7]*m.e[13] -
		m.e[4]*m.e[1]*m.e[15] +
		m.e[4]*m.e[3]*m.e[13] +
		m.e[12]*m.e[1]*m.e[7] -
		m.e[12]*m.e[3]*m.e[5]

	inv[14] = -m.e[0]*m.e[5]*m.e[14] +
		m.e[0]*m.e[6]*m.e[13] +
		m.e[4]*m.e[1]*m.e[14] -
		m.e[4]*m.e[2]*m.e[13] -
		m.e[12]*m.e[1]*m.e[6] +
		m.e[12]*m.e[2]*m.e[5]

	inv[3] = -m.e[1]*m.e[6]*m.e[11] +
		m.e[1]*m.e[7]*m.e[10] +
		m.e[5]*m.e[2]*m.e[11] -
		m.e[5]*m.e[3]*m.e[10] -
		m.e[9]*m.e[2]*m.e[7] +
		m.e[9]*m.e[3]*m.e[6]

	inv[7] = m.e[0]*m.e[6]*m.e[11] -
		m.e[0]*m.e[7]*m.e[10] -
		m.e[4]*m.e[2]*m.e[11] +
		m.e[4]*m.e[3]*m.e[10] +
		m.e[8]*m.e[2]*m.e[7] -
		m.e[8]*m.e[3]*m.e[6]

	inv[11] = -m.e[0]*m.e[5]*m.e[11] +
		m.e[0]*m.e[7]*m.e[9] +
		m.e[4]*m.e[1]*m.e[11] -
		m.e[4]*m.e[3]*m.e[9] -
		m.e[8]*m.e[1]*m.e[7] +
		m.e[8]*m.e[3]*m.e[5]

	inv[15] = m.e[0]*m.e[5]*m.e[10] -
		m.e[0]*m.e[6]*m.e[9] -
		m.e[4]*m.e[1]*m.e[10] +
		m.e[4]*m.e[2]*m.e[9] +
		m.e[8]*m.e[1]*m.e[6] -
		m.e[8]*m.e[2]*m.e[5]

	det := m.e[0]*inv[0] + m.e[1]*inv[4] + m.e[2]*inv[8] + m.e[3]*inv[12]

	if det == 0 {
		return false
	}

	det = 1.0 / det

	for i := 0; i < 16; i++ {
		m.e[i] = inv[i] * det
	}

	return true
}

func (m matrix4) String() string {
	s := fmt.Sprintf("[%7.3f, %7.3f, %7.3f, %7.3f]\n", m.e[M00], m.e[M01], m.e[M02], m.e[M03])
	s += fmt.Sprintf("[%7.3f, %7.3f, %7.3f, %7.3f]\n", m.e[M10], m.e[M11], m.e[M12], m.e[M13])
	s += fmt.Sprintf("[%7.3f, %7.3f, %7.3f, %7.3f]\n", m.e[M20], m.e[M21], m.e[M22], m.e[M23])
	s += fmt.Sprintf("[%7.3f, %7.3f, %7.3f, %7.3f]", m.e[M30], m.e[M31], m.e[M32], m.e[M33])
	return s
}
