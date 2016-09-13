package engine

import (
	"fmt"
	"log"
	"math"

	"phys/vect"

	"github.com/go-gl/mathgl/mgl32"
)

func (o *Object) Position() (x, y float32) {
	// if o.Shape != nil && o.Node != nil {
	// 	log.Println(o.Shape.Body.Position(), o.Node.Location)
	// }

	if o.Shape != nil {
		v := o.Shape.Body.Position()
		return v.X, v.Y
	}

	return o.Node.Location.X(), o.Node.Location.Y()
}

func (o *Object) PositionVec2() mgl32.Vec2 {
	if o.Shape != nil {
		v := o.Shape.Body.Position()
		return mgl32.Vec2{v.X, v.Y}
	}

	return mgl32.Vec2{o.Node.Location.X(), o.Node.Location.Y()}
}

func (o *Object) PositionVec3() mgl32.Vec3 {
	if o.Shape != nil {
		v := o.Shape.Body.Position()
		return mgl32.Vec3{v.X, v.Y, 0}
	}

	return mgl32.Vec3{o.Node.Location.X(), o.Node.Location.Y(), o.Node.Location.Z()}
}

func (e *Object) SetPosition(x, y float32) {
	if e.Shape != nil {
		e.Shape.Body.SetPosition(vect.Vect{x, y})
	} else {
		e.Node.Location = mgl32.Vec3{x, y, 0}
	}
}

func (e *Object) VectorForward(scale float32) (float32, float32) {
	pvx, pvy := e.Shape.Body.Rot()
	return pvx * scale, pvy * scale
}

func (e *Object) VectorSide(scale float32) (float32, float32) {
	// ang := e.Shape.Body.Angle() - 1.5708 // ~90 deg
	// if scale < 0 {
	// 	ang = e.Shape.Body.Angle() + 1.5708 // ~90 deg
	// }

	v := vect.FromAngle(e.Shape.Body.Angle() - 1.5708)
	v.Mult(scale)

	return v.X, v.Y
}

// func (e *Object) VectorRight(scale float32) (float32, float32) {
// 	ang := e.Shape.Body.Angle() + 1.5708 // ~90 deg
// 	v := vect.FromAngle(ang)
// 	v.Mult(scale)
// 	return v.X, v.Y
// }

func (e *Object) Rotation() float32 {
	if e.Shape != nil {
		return e.Shape.Body.Angle()
	}
	log.Println("need check get rotation fron fizzle node, it`s may be not correct!")
	return e.Node.Rotation.V.Z()
}

func (e *Object) SetRotation(ang float32) {
	if e.Shape != nil {
		e.Shape.Body.SetAngle(ang)
	} else {
		// e.Node.LocalRotation = mgl32.AnglesToQuat(0, 0, vect.Vect{x, y}.Angle(), 1)
		e.Node.LocalRotation = mgl32.AnglesToQuat(0, 0, ang, 1)
	}
}

//Quat2Rad convert Quaternion to rad euler angles
func Quat2Rad(q mgl32.Quat) (x, y, z, a float64) {
	qw := float64(q.W)
	qx := float64(q.X())
	qy := float64(q.Y())
	qz := float64(q.Z())

	// x = math.Atan2(2*qx*qw-2*qy*qz, 1-2*qx*qx-2*qz*qz)
	// y = math.Atan2(2*qy*qw-2*qx*qz, 1-2*qy*qy-2*qz*qz)
	// z = math.Asin(2*qx*qy + 2*qz*qw)

	fmt.Println(qw, qx, qy, qz)
	a = 2 * math.Acos(qw)
	x = math.Acos(qx / math.Sin(a/2))
	y = math.Acos(qy / math.Sin(a/2))
	z = math.Acos(qz / math.Sin(a/2))

	// a = 2 * math.Acos(qw)

	return
}

//Rad2Quat convert rad euler angles to quternion
func Rad2Quat(x, y, z float32) mgl32.Quat {
	qx := float64(x)
	qy := float64(y)
	qz := float64(z)

	q0 := math.Cos(qx/2)*math.Cos(qy/2)*math.Cos(qz/2) + math.Sin(qx/2)*math.Sin(qy/2)*math.Sin(qz/2)
	q1 := math.Sin(qx/2)*math.Cos(qy/2)*math.Cos(qz/2) - math.Cos(qx/2)*math.Sin(qy/2)*math.Sin(qz/2)
	q2 := math.Cos(qx/2)*math.Sin(qy/2)*math.Cos(qz/2) + math.Sin(qx/2)*math.Cos(qy/2)*math.Sin(qz/2)
	q3 := math.Cos(qx/2)*math.Cos(qy/2)*math.Sin(qz/2) - math.Sin(qx/2)*math.Sin(qy/2)*math.Cos(qz/2)

	return mgl32.Quat{W: float32(q0), V: mgl32.Vec3{float32(q1), float32(q2), float32(q3)}}
}

func Rad2Deg(a float32) float32 {
	return a * math.Pi * 180
}

func Distance(x0, y0, x1, y1 float32) float32 {
	dist := math.Sqrt(math.Pow(float64(x0)-float64(x1), 2) + math.Pow(float64(y0)-float64(y1), 2))
	return float32(dist)
}

func (o1 *Object) Distance(o2 *Object) float32 {
	cpx, cpy := o1.Position()
	ppx, ppy := o2.Position()

	dist := math.Sqrt(math.Pow(float64(cpx)-float64(ppx), 2) + math.Pow(float64(cpy)-float64(ppy), 2))
	return float32(dist)
}

func (e *Object) SetVelocity(x, y float32) {
	if e.Shape == nil {
		log.Println("velocity can not be set - shape is nil, name of object:", e.Name)
		return
	}
	e.Shape.Body.SetVelocity(x, y)

}

func (e *Object) AddVelocity(x, y float32) {
	if e.Shape == nil {
		log.Println("velocity can not be set - shape is nil, name of object:", e.Name)
		return
	}
	e.Shape.Body.AddVelocity(x, y)
}

func (e *Object) Velocity() vect.Vect {
	if e.Shape == nil {
		log.Println("velocity can not be set - shape is nil, name of object:", e.Name)
		return vect.Vect{}
	}
	return e.Shape.Body.Velocity()
}

func (e *Object) Destroy() {
	// log.Println("destroy:", e.Name)
	if e.Shape != nil {
		e.Shape.Body.Enabled = false
		space.RemoveBody(e.Shape.Body)
	}

	Objects[e] = false

	delete(Objects, e)
	e = nil
}
