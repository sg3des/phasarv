package vect

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	Vector_Zero = Vect{0, 0}
)

func FMin(a, b float32) float32 {
	if a > b {
		return b
	}
	return a
}

func FAbs(a float32) float32 {
	if a < 0 {
		return -a
	}
	return a
}

func FMax(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func FClamp(val, min, max float32) float32 {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

//basic 2d vector.
type Vect struct {
	X, Y float32
}

//adds v2 to the given vector.
func (v1 *Vect) Add(v2 Vect) {
	v1.X += v2.X
	v1.Y += v2.Y
}

//subtracts v2 rom the given vector.
func (v1 *Vect) Sub(v2 Vect) {
	v1.X -= v2.X
	v1.Y -= v2.Y
}

//subtracts v2 rom the given vector.
func (v *Vect) Angle() float32 {
	return float32(math.Atan2(float64(v.Y), float64(v.X)))
}

//returns the squared length of the vector.
func (v Vect) LengthSqr() float32 {
	//length of a vector: distance to origin
	return DistSqr(v, Vect{})
}

//returns the length of the vector.
func (v Vect) Length() float32 {
	//length of a vector: distance to origin
	return Dist(v, Vect{})
}

//multiplies the vector by the scalar.
func (v *Vect) Mult(s float32) {
	v.X *= s
	v.Y *= s
}

//normalizes the vector to a length of 1.
func (v *Vect) Normalize() {
	f := 1.0 / v.Length()
	v.X *= f
	v.Y *= f
}

//compare two vectors by value.
func Equals(v1, v2 Vect) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}

//adds the input vectors and returns the result.
func Add(v1, v2 Vect) Vect {
	return Vect{v1.X + v2.X, v1.Y + v2.Y}
}

//subtracts the input vectors and returns the result.
func Sub(v1, v2 Vect) Vect {
	return Vect{v1.X - v2.X, v1.Y - v2.Y}
}

//multiplies a vector by a scalar and returns the result.
func Mult(v1 Vect, s float32) Vect {
	return Vect{v1.X * s, v1.Y * s}
}

//returns the square distance between two vectors.
func DistSqr(v1, v2 Vect) float32 {
	return (v1.X-v2.X)*(v1.X-v2.X) + (v1.Y-v2.Y)*(v1.Y-v2.Y)
}

//returns the distance between two vectors.
func Dist(v1, v2 Vect) float32 {
	return float32(math.Sqrt(float64(DistSqr(v1, v2))))
}

//returns the squared length of the vector.
func LengthSqr(v Vect) float32 {
	//length of a vector: distance to origin
	return DistSqr(v, Vect{})
}

//returns the length of the vector.
func Length(v Vect) float32 {
	//length of a vector: distance to origin
	return Dist(v, Vect{})
}

/*
func Min(v1, v2 Vect) Vect {
	return INLI_Min(v1, v2)
}

func Max(v1, v2 Vect) Vect {
	return INLI_Max(v1, v2)
}
*/
func Min(v1, v2 Vect) (out Vect) {
	if v1.X < v2.X {
		out.X = v1.X
	} else {
		out.X = v2.X
	}

	if v1.Y < v2.Y {
		out.Y = v1.Y
	} else {
		out.Y = v2.Y
	}
	return
}

//returns a new vector with its x/y values set to the bigger one from the two input values.
//e.g. Min({2, 10}, {8, 3}) would return {2, 3}
func Max(v1, v2 Vect) (out Vect) {
	if v1.X > v2.X {
		out.X = v1.X
	} else {
		out.X = v2.X
	}

	if v1.Y > v2.Y {
		out.Y = v1.Y
	} else {
		out.Y = v2.Y
	}
	return
}

//returns the normalized input vector.
func Normalize(v Vect) Vect {
	f := 1.0 / Length(v)
	return Vect{v.X * f, v.Y * f}
}

//dot product between two vectors.
func Dot(v1, v2 Vect) float32 {
	return (v1.X * v2.X) + (v1.Y * v2.Y)
}

//same as CrossVV.
func Cross(a, b Vect) float32 {
	return (a.X * b.Y) - (a.Y * b.X)
}

func Clamp(v Vect, l float32) Vect {
	if Dot(v, v) > l*l {
		return Mult(Normalize(v), l)
	}
	return v
}

//cross product of two vectors.
func CrossVV(a, b Vect) float32 {
	return (a.X * b.Y) - (a.Y * b.X)
}

//cross product between a vector and a float64.
//result = {s * a.Y, -s * a.X}
func CrossVF(a Vect, s float32) Vect {
	return Vect{s * a.Y, -s * a.X}
}

//cross product between a float64 and a vector.
//Not the same as CrossVD
//result = {-s * a.Y, s * a.X}
func CrossFV(s float32, a Vect) Vect {
	return Vect{-s * a.Y, s * a.X}
}

//linear interpolation between two vectors by the given scalar
func Lerp(v1, v2 Vect, s float32) Vect {
	return Vect{
		v1.X + (v1.X-v2.X)*s,
		v1.Y + (v1.Y-v2.Y)*s,
	}
}

//Returns v rotated by 90 degrees
func Perp(v Vect) Vect {
	return Vect{-v.Y, v.X}
}

func FromAngle(angle float32) Vect {
	return Vect{float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))}
}

//Intersection of two lines ab - cd incoming 4 points vect.Vect
func Intersection(a, b, c, d Vect) bool {
	baX := b.X - a.X
	baY := b.Y - a.Y
	dcX := d.X - c.X
	dcY := d.Y - c.Y

	denom := baX*dcY - dcX*baY
	if denom == 0 {
		return false //parallel
	}

	denomPositive := denom > 0

	acX := a.X - c.X
	acY := a.Y - c.Y
	s := baX*acY - baY*acX
	if (s < 0) == denomPositive {
		return false
	}

	t := dcX*acY - dcY*acX
	if (t < 0) == denomPositive {
		return false
	}

	if ((s >= denom) == denomPositive) || ((t >= denom) == denomPositive) {
		return false
	}

	return true
}

// phasarv

func FromVec3(v mgl32.Vec3) Vect {
	return Vect{v.X(), v.Y()}
}

func FromVec2(v mgl32.Vec2) Vect {
	return Vect{v.X(), v.Y()}
}

//Vec2 return mgl32.Vec2 from Vect
func (v Vect) Vec2() mgl32.Vec2 {
	return mgl32.Vec2{v.X, v.Y}
}

//Vec3 return mgl32.Vec3 from Vect with zero in Z axis
func (v Vect) Vec3() mgl32.Vec3 {
	return mgl32.Vec3{v.X, v.Y, 0}
}

//SubPoint return absolute point from relative offset of v1 and angle
func (v1 Vect) SubPoint(angle float32, v2 Vect) Vect {
	scale := Dist(Vect{}, v2)

	v := FromAngle(angle + v2.Angle())
	v.Mult(scale)
	v.Add(v1)

	return v
}

func (v1 Vect) SubAngle(angle float32, v2 Vect) float32 {

	// a := o.PositionVec2()

	// var oAngleVec float32
	// if o.Shape != nil {
	oAngleVec := FromAngle(angle)
	// }else{
	// 	oAngleVec =
	// }

	//angle between points
	abAngle := float32(math.Atan2(float64(v2.Y-v1.Y), float64(v2.X-v1.X)))
	abAngleVec := FromAngle(abAngle)

	sin := oAngleVec.X*abAngleVec.Y - abAngleVec.X*oAngleVec.Y
	cos := oAngleVec.X*abAngleVec.X + oAngleVec.Y*abAngleVec.Y

	return float32(math.Atan2(float64(sin), float64(cos)))

}
