package engine

import (
	"fmt"
	"log"
	"math"

	"phys/vect"

	"github.com/go-gl/mathgl/mgl32"
)

func (o *Object) Angle() float32 {
	if o.Shape != nil {
		return o.Shape.Body.Angle()
	}
	x, y := o.Position()
	return (&vect.Vect{x, y}).Angle()
}

// func (o *Object) Length(x, y float32) float32 {
// 	ox, oy := o.Position()
// 	return vect.Dist(vect.Vect{ox, oy}, vect.Vect{x, y})
// }

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
		e.Node.Location = mgl32.Vec3{x, y, e.Node.Location.Z()}
	}
}

func (e *Object) VectorForward(scale float32) (float32, float32) {
	var v vect.Vect
	if e.Shape != nil && e.Shape.Body != nil {
		// v = vect.FromAngle(e.Shape.Body.Angle())
		v = e.Shape.Body.RotVec()
	} else {
		// v = vect.Vect{e.Node.LocalRotation.X(), e.Node.LocalRotation.Y()}

		v = vect.FromAngle(float32(2 * math.Acos(float64(e.Node.Rotation.W))))
	}

	v.Mult(scale)

	// pvx, pvy := e.Shape.Body.Rot()
	return v.X, v.Y
}

func (o *Object) VectorSide(scale, angle float32) (float32, float32) {
	// ang := e.Shape.Body.Angle() - 1.5708 // ~90 deg
	// if scale < 0 {
	// 	ang = e.Shape.Body.Angle() + 1.5708 // ~90 deg
	// }

	rotvec := o.Shape.Body.RotVec()
	// rotvec.Mult(scale)

	// log.Println(o.Shape.Body.Angle() == rotvec.Angle())

	v := vect.FromAngle(rotvec.Angle() + angle)
	v.Mult(scale)

	return v.X, v.Y
}

func (o *Object) Rotation() float32 {
	if o.Shape != nil {
		return o.Shape.Body.Angle()
	}
	// log.Println("need check get rotation fron fizzle node, it`s may be not correct!")
	return float32(2 * math.Acos(float64(o.Node.Rotation.W)))
}

func (o *Object) SetRotation(ang float32) {
	if o.Shape != nil {
		o.Shape.Body.SetAngle(ang)
	} else {
		// e.Node.LocalRotation = mgl32.AnglesToQuat(0, 0, vect.Vect{x, y}.Angle(), 1)
		o.Node.LocalRotation = mgl32.AnglesToQuat(0, 0, ang, 1)
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

func (o *Object) Distance(o2 *Object) float32 {
	if o2 == nil {
		log.Fatalln("target object is nil")
	}
	cpx, cpy := o.Position()
	ppx, ppy := o2.Position()

	dist := math.Sqrt(math.Pow(float64(cpx)-float64(ppx), 2) + math.Pow(float64(cpy)-float64(ppy), 2))
	return float32(dist)
}

func (o *Object) DistancePoint(bx, by float32) float32 {
	cpx, cpy := o.Position()

	dist := math.Sqrt(math.Pow(float64(cpx)-float64(bx), 2) + math.Pow(float64(cpy)-float64(by), 2))
	return float32(dist)
}

func (o *Object) SetVelocity(x, y float32) {
	if o.Shape != nil {
		o.Shape.Body.SetVelocity(x, y)
		return
	}
	log.Println("velocity can not be set - shape is nil, name of object:", o.Name)
	return
}

func (o *Object) AddVelocity(x, y float32) {
	if o.Shape != nil {
		o.Shape.Body.AddVelocity(x, y)
		return
	}
	log.Println("velocity can not be set - shape is nil, name of object:", o.Name)
	return
}

func (o *Object) Velocity() vect.Vect {
	if o.Shape != nil {
		return o.Shape.Body.Velocity()
	}

	log.Println("velocity can not be set - shape is nil, name of object:", o.Name)
	return vect.Vect{}
}

func (o *Object) ShapeWidthPercent() float32 {
	return vect.FAbs(o.RollAngle) / (o.Param.MaxRollAngle * 1.1)
}

func (o *Object) Destroy() {
	if o.DestroyFunc != nil {
		o.DestroyFunc()
		return
	}

	if o.Shape != nil {
		o.Shape.Body.Enabled = false
		space.RemoveBody(o.Shape.Body)
	}

	Objects[o] = false

	delete(Objects, o)
	o = nil
}

func (o *Object) Clone() *Object {
	newObject := &Object{
		Name:         o.Name,
		Node:         o.Node.Clone(),
		MaxRollAngle: o.MaxRollAngle,

		Shadow:      o.Shadow,
		Transparent: o.Transparent,

		ArtStatic: o.ArtStatic,
		ArtRotate: o.ArtRotate,

		Callback:    o.Callback,
		DestroyFunc: o.DestroyFunc,

		Param: o.Param,
	}

	newObject.SetPhys(o.Param.Phys)
	newObject.Node.Material = Material(o.Param.Material)

	Objects[newObject] = true

	return newObject
}
