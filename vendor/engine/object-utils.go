package engine

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"phys"
	"time"

	"phys/vect"

	"github.com/go-gl/mathgl/mgl32"
)

func (o *Object) Position() (x, y float32) {
	if o.shape != nil {
		v := o.shape.Body.Position()
		return v.X, v.Y
	}

	return o.P.Pos.X, o.P.Pos.Y
}

func (o *Object) PositionVect() vect.Vect {
	if o.shape != nil {
		return o.shape.Body.Position()
	}

	return o.P.Pos.Vect()
}

func (o *Object) PositionVec2() mgl32.Vec2 {
	if o.shape != nil {
		// log.Println(o.shape.Body)
		v := o.shape.Body.Position()
		return mgl32.Vec2{v.X, v.Y}
	}

	return o.P.Pos.Vec2()
}

func (o *Object) PositionVec3() mgl32.Vec3 {
	if o.shape != nil {
		v := o.shape.Body.Position()
		return mgl32.Vec3{v.X, v.Y, 0}
	}

	return o.P.Pos.Vec3()
}

func (o *Object) SetPosition(x, y float32) {
	if o.shape != nil {
		o.shape.Body.SetPosition(vect.Vect{x, y})
	} else {
		o.P.Pos.X = x
		o.P.Pos.Y = y
	}
}

func (o *Object) VectorForward(scale float32) (float32, float32) {
	var v vect.Vect
	if o.shape != nil && o.shape.Body != nil {
		v = o.shape.Body.RotVec()
	} else {
		v = vect.FromAngle(o.P.Angle)
	}

	v.Mult(scale)

	return v.X, v.Y
}

func (o *Object) VectorSide(scale, angle float32) (float32, float32) {

	var v vect.Vect
	if o.shape != nil && o.shape.Body != nil {
		v = o.shape.Body.RotVec()
		v = vect.FromAngle(v.Angle() + angle)
	} else {
		v = vect.FromAngle(o.P.Angle + angle)
	}
	v.Mult(scale)

	return v.X, v.Y
}

func (o *Object) Rotation() float32 {
	if o.shape != nil {
		return o.shape.Body.Angle()
	}
	return o.P.Angle
}

func (o *Object) SetRotation(ang float32) {
	if o.shape != nil {
		o.shape.Body.SetAngle(ang)
	} else {
		o.P.Angle = ang
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
		return 0
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
	if o.shape != nil {
		o.shape.Body.SetVelocity(x, y)
		return
	}
	log.Println("WARNING: velocity can not be set - shape is nil, name of object:", o.Name)
	return
}

func (o *Object) SetAngularVelocity(w float32) {
	if o.shape != nil {
		o.shape.Body.SetAngularVelocity(w)
		return
	}
	log.Println("WARNING: angular velocity can not be set - shape is nil, name of object:", o.Name)
	return
}

func (o *Object) AddVelocity(x, y float32) {
	if o.shape != nil {
		o.shape.Body.AddVelocity(x, y)
		return
	}
	log.Println("velocity can not be set - shape is nil, name of object:", o.Name)
	return
}

func (o *Object) Velocity() vect.Vect {
	if o.shape != nil {
		return o.shape.Body.Velocity()
	}

	log.Println("velocity can not be set - shape is nil, name of object:", o.Name)
	return vect.Vect{}
}

func (o *Object) ShapeWidthPercent() float32 {
	return vect.FAbs(o.RollAngle) / (o.MaxRollAngle * 1.1)
}

func (o *Object) SubAngleToPoint(b mgl32.Vec2, off mgl32.Vec2) float32 {
	var v mgl32.Vec2

	if off.Len() > 0 {
		vx, vy := o.VectorSide(2, -1.5704)
		v = mgl32.Vec2{vx, vy}
	}

	// log.Println(off.Len(), v)
	a := o.PositionVec2().Add(v)
	oAngleVec := vect.FromAngle(o.Rotation())

	abAngle := float32(math.Atan2(float64(b.Y()-a.Y()), float64(b.X()-a.X())))
	abAngleVec := vect.FromAngle(abAngle)

	sin := oAngleVec.X*abAngleVec.Y - abAngleVec.X*oAngleVec.Y
	cos := oAngleVec.X*abAngleVec.X + oAngleVec.Y*abAngleVec.Y

	return float32(math.Atan2(float64(sin), float64(cos)))
}

func (o *Object) SubPoint(b mgl32.Vec2) mgl32.Vec2 {
	a := o.PositionVec2()
	scale := a.Sub(b).Len()

	var v = vect.Vect{b.X(), b.Y()}
	angle := v.Angle()

	// oAngleVec := vect.FromAngle(o.Angle())

	// abAngle := float32(math.Atan2(float64(b.Y()-a.Y()), float64(b.X()-a.X())))
	// abAngleVec := vect.FromAngle(abAngle)

	// sin := oAngleVec.X*abAngleVec.Y - abAngleVec.X*oAngleVec.Y
	// cos := oAngleVec.X*abAngleVec.X + oAngleVec.Y*abAngleVec.Y

	// angle := float32(math.Atan2(float64(sin), float64(cos)))

	x, y := o.VectorSide(scale, angle)

	log.Println(x, y, scale, angle)
	return mgl32.Vec2{x, y}
}

func (o *Object) AddAngularVelocity(w float32) {
	if o.shape != nil && o.shape.Body != nil {
		o.shape.Body.AddAngularVelocity(w)
	}
}

func (o *Object) AngularVelocity() float32 {
	if o.shape == nil || o.shape.Body == nil {
		return 0
	}

	return o.shape.Body.AngularVelocity()
}

func (o *Object) Raycast(x0, y0, x1, y1 float32) (objects []*Object) {
	if o.shape == nil || o.shape.Body == nil {
		log.Println("failed raycast for non phys object", o.Name)
		return
	}

	hits := phys.Hits(x0, y0, x1, y1, phys.GROUP_STATIC, o.shape.Body)
	hits = append(hits, phys.Hits(x0, y0, x1, y1, phys.GROUP_PLAYER, o.shape.Body)...)

	dist := Distance(x0, y0, x1, y1)
	for _, hit := range hits {
		if hit.Distance <= dist {
			objects = append(objects, hit.Shape.UserData.(*Object))
		}
	}

	return
}

//GetNearObjectByRay get nearest object by ray betwean 2 points, ignore source object
func (o *Object) GetNearObjectByRay(x0, y0, x1, y1 float32) (near *Object, shortDist float32) {

	objects := o.Raycast(x0, y0, x1, y1)

	//hardcoded maximum distance
	shortDist = 9999
	for _, obj := range objects {

		if obj != nil {
			if dist := obj.DistancePoint(x0, y0); dist < shortDist {
				shortDist = dist
				near = obj
			}
		}
	}

	return
}

//SubAngleObjectPoint calculate angle(rad) between object(o) angle and 2d point(b)
func (o *Object) SubAngleObjectPoint(b mgl32.Vec2) (angle float32) {
	a := o.PositionVec2()

	// var oAngleVec float32
	// if o.Shape != nil {
	oAngleVec := vect.FromAngle(o.Rotation())
	// }else{
	// 	oAngleVec =
	// }

	//angle between points
	abAngle := float32(math.Atan2(float64(b.Y()-a.Y()), float64(b.X()-a.X())))
	abAngleVec := vect.FromAngle(abAngle)

	sin := oAngleVec.X*abAngleVec.Y - abAngleVec.X*oAngleVec.Y
	cos := oAngleVec.X*abAngleVec.X + oAngleVec.Y*abAngleVec.Y

	return float32(math.Atan2(float64(sin), float64(cos)))
}

func GetObjectInPoint(a mgl32.Vec2) *Object {
	shape := phys.Hit(a[0], a[1], phys.GROUP_PLAYER)
	if shape == nil {
		return nil
	}

	return shape.UserData.(*Object)
}

func GetRandomPoint(x, y float32) (float32, float32) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return x/2 - r.Float32()*x, y/2 - r.Float32()*y
}
