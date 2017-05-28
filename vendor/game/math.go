package game

import (
	"engine"
	"math"
	"phys"

	"phys/vect"

	"github.com/go-gl/mathgl/mgl32"
)

// func SubPoint(angle, x, y float32) (float32, float32) {
// 	b := vect.Vect{x, y}
// 	scale := vect.Dist(vect.Vect{}, b)

// 	v := vect.FromAngle(angle + b.Angle())
// 	v.Mult(scale)

// 	return v.X, v.Y
// }

func GetPlayerInPoint(x, y float32) *Player {
	shape := phys.Hit(x, y, phys.GROUP_PLAYER)
	if shape == nil {
		return nil
	}

	return shape.UserData.(*engine.Object).UserData.(*Player)

	// if userData := engine.Hit(x, y); userData != nil {
	// 	return userData.(*engine.Object)
	// }
	// return nil
}

func GetNearObjectByRay(x0, y0, x1, y1 float32, ignoreObject *engine.Object) (o *engine.Object, shortDist float32) {

	objects := ignoreObject.Raycast(x0, y0, x1, y1)

	shortDist = 999
	for _, obj := range objects {

		if obj != nil {
			if dist := obj.DistancePoint(x0, y0); dist < shortDist {
				shortDist = dist
				o = obj
			}
		}
	}

	return
}

// func AngleBetweenPoints(px, py, cx, cy float32) float32 {
// 	return float32(math.Atan2(float64(cy-py), float64(cx-px)))
// }

// func AngleBetweenAngles(pa, ca float32) float32 {
// 	cv := vect.FromAngle(ca)
// 	pv := vect.FromAngle(pa)

// 	sin := pv.X*cv.Y - cv.X*pv.Y
// 	cos := pv.X*cv.X + pv.Y*cv.Y

// 	return float32(math.Atan2(float64(sin), float64(cos)))
// }

func getCursorPos(x, y, w, h float32, campos mgl32.Vec3) (float32, float32) {
	// m := engine.Camera.GetViewMatrix()

	// p := mgl32.Perspective(mgl32.DegToRad(50), w/h, 1.0, 100.0)
	// rowy, rowx, _, _ := p.Rows()
	// px := rowx[1]
	// py := rowy[0]

	// curX := x/w - 0.5
	// curY := y/h - 0.5

	// log.Println(curX, curY, px, py)

	// normalised_x := 2*x/w - 1
	// normalised_y := 1 - 2*y/h
	// log.Println(normalised_x, normalised_y)

	// unviewMat = (projectionMat * modelViewMat).inverse()

	// near_point = unviewMat * Vec(normalised_x, normalised_y, 0, 1)
	// camera_pos = ray_origin = modelViewMat.inverse().col(4)
	// ray_dir = near_point - camera_pos

	// log.Println(curX, curY)
	//that`s black magic constant!
	// var d float32 = 840
	var d = h + campos.Z()

	x = (x-w/2)/d*campos.Z() + campos.X()
	y = (h/2-y)/d*campos.Z() + campos.Y()

	// mgl32.ScreenToGLCoords(x, y, screenWidth, screenHeight)
	// y = (y-float32(h)/2)/820*campos.Z() + campos.Y()

	// log.Println(x, y)
	return x, y
}

// //LookAtTarget slowly turns object to target 2d-point(cpx,cpy)
// func LookAtTarget(o *engine.Object, cpx, cpy, dt float32) (angle float32) {
// 	pp := o.Shape.Body.Position()
// 	ca := AngleBetweenPoints(pp.X, pp.Y, cpx, cpy)
// 	pa := o.Shape.Body.Angle()
// 	// log.Println(pp.X, pp.Y, cpx, cpy, " | ", ca, pa)

// 	// //get rotspeed
// 	// var rotspeed float32
// 	// if o.Player != nil {
// 	// 	rotspeed = o.Player.RotSpeed * 0.1
// 	// }
// 	// if o.Bullet != nil {
// 	// 	rotspeed = o.Bullet.Param.RotSpeed * 0.1
// 	// }

// 	//get angle
// 	return AngleBetweenAngles(pa, ca)

// 	// if ang > rotspeed {
// 	// 	ang = rotspeed
// 	// } else if ang < -rotspeed {
// 	// 	ang = -rotspeed
// 	// }

// 	// // add angluar velocity
// 	// avel := o.Shape.Body.AngularVelocity()
// 	// if avel < rotspeed && avel > -rotspeed {
// 	// 	o.Shape.Body.AddAngularVelocity(ang * 0.001)
// 	// }

// 	// //roll angle
// 	// if o.Player != nil {
// 	// 	// subAngle := AngleBetweenAngles(pa, ca)

// 	// 	if ang > -o.Player.SubAngle && o.Player.SubAngle > -1.5 {
// 	// 		o.Player.SubAngle -= dt * o.Player.RotSpeed * 0.05
// 	// 	}

// 	// 	if ang < -o.Player.SubAngle && o.Player.SubAngle < 1.5 {
// 	// 		o.Player.SubAngle += dt * o.Player.RotSpeed * 0.05
// 	// 	}
// 	// }

// 	// return
// }

//AngleObjectPoint calculate angle(rad) between object(o) angle and 2d point(b)
func SubAngleObjectPoint(o *engine.Object, b mgl32.Vec2) (angle float32) {
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

func AngleObjectPoint(o *engine.Object, b mgl32.Vec2) float32 {
	a := o.PositionVec2()
	//angle between points
	return float32(math.Atan2(float64(b.Y()-a.Y()), float64(b.X()-a.X())))
	// abAngleVec := vect.FromAngle(abAngle)

	// sin := v.X*abAngleVec.Y - abAngleVec.X*v.Y
	// cos := v.X*abAngleVec.X + v.Y*abAngleVec.Y

	// return float32(math.Atan2(float64(sin), float64(cos)))
}

//check if time is nil
// func timeNil(t time.Time) bool {
// 	if t.Equal(time.Time{}) {
// 		return true
// 	}
// 	return false
// }
