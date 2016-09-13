package transform

import (
	"phys/vect"
	"math"
)

type Rotation struct {
	//sine and cosine.
	C, S float32
}

func NewRotation(angle float32) Rotation {
	return Rotation{
		C: float32(math.Cos(float64(angle))),
		S: float32(math.Sin(float64(angle))),
	}
}

func (rot *Rotation) SetIdentity() {
	rot.S = 0
	rot.C = 1
}

func (rot *Rotation) SetAngle(angle float32) {
	rot.C = float32(math.Cos(float64(angle)))
	rot.S = float32(math.Sin(float64(angle)))
}

func (rot *Rotation) Angle() float32 {
	return float32(math.Atan2(float64(rot.S), float64(rot.C)))
}

//rotates the input vector.
func (rot *Rotation) RotateVect(v vect.Vect) vect.Vect {
	return vect.Vect{
		X: (v.X * rot.C) - (v.Y * rot.S),
		Y: (v.X * rot.S) + (v.Y * rot.C),
	}
}

//rotates the input vector.
func (rot *Rotation) RotateVectPtr(v *vect.Vect) vect.Vect {
	return vect.Vect{
		X: (v.X * rot.C) - (v.Y * rot.S),
		Y: (v.X * rot.S) + (v.Y * rot.C),
	}
}

func (rot *Rotation) RotateVectInv(v vect.Vect) vect.Vect {
	return vect.Vect{
		X: (v.X * rot.C) + (v.Y * rot.S),
		Y: (-v.X * rot.S) + (v.Y * rot.C),
	}
}

func RotateVect(v vect.Vect, r Rotation) vect.Vect {
	return r.RotateVect(v)
}

func RotateVectPtr(v *vect.Vect, r *Rotation) vect.Vect {
	return r.RotateVectPtr(v)
}

func RotateVectInv(v vect.Vect, r Rotation) vect.Vect {
	return r.RotateVectInv(v)
}

type Transform struct {
	Position vect.Vect
	Rotation
}

func NewTransform(pos vect.Vect, angle float32) Transform {
	return Transform{
		Position: pos,
		Rotation: NewRotation(angle),
	}
}

func NewTransform2(pos vect.Vect, rot vect.Vect) Transform {
	return Transform{
		Position: pos,
		Rotation: Rotation{rot.X, rot.Y},
	}
}

func (xf *Transform) SetIdentity() {
	xf.Position = vect.Vect{}
	xf.Rotation.SetIdentity()
}

func (xf *Transform) Set(pos vect.Vect, rot float32) {
	xf.Position = pos
	xf.SetAngle(rot)
}

//moves and roates the input vector.
func (xf *Transform) TransformVect(v vect.Vect) vect.Vect {
	return vect.Add(xf.Position, xf.RotateVect(v))
}

func (xf *Transform) TransformVectInv(v vect.Vect) vect.Vect {
	return vect.Add(vect.Mult(xf.Position, -1), xf.RotateVectInv(v))
}
