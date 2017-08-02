package game

import (
	"engine"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// func SubPoint(angle, x, y float32) (float32, float32) {
// 	b := vect.Vect{x, y}
// 	scale := vect.Dist(vect.Vect{}, b)

// 	v := vect.FromAngle(angle + b.Angle())
// 	v.Mult(scale)

// 	return v.X, v.Y
// }

// func GetPlayerInPoint(x, y float32) *Player {
// 	shape := phys.Hit(x, y, phys.GROUP_PLAYER)
// 	if shape == nil {
// 		return nil
// 	}

// 	return shape.UserData.(*engine.Object).UserData.(*Player)

// 	// if userData := engine.Hit(x, y); userData != nil {
// 	// 	return userData.(*engine.Object)
// 	// }
// 	// return nil
// }

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
