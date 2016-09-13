package main

import (
	"engine"
	"math"

	"phys/vect"

	"github.com/go-gl/mathgl/mgl32"
)

func AngleBetweenPoints(px, py, cx, cy float32) float32 {
	return float32(math.Atan2(float64(cy-py), float64(cx-px)))
}

func AngleBetweenAngles(pa, ca float32) float32 {
	cv := vect.FromAngle(ca)
	pv := vect.FromAngle(pa)

	sin := pv.X*cv.Y - cv.X*pv.Y
	cos := pv.X*cv.X + pv.Y*cv.Y

	return float32(math.Atan2(float64(sin), float64(cos)))
}
func getCursorPos(x, y float32, w, h int, campos mgl32.Vec3) (float32, float32) {
	x = (x-float32(w)/2)/820*campos.Z() + campos.X()
	y = (float32(h)/2-y)/820*campos.Z() + campos.Y()
	// y = (y-float32(h)/2)/820*campos.Z() + campos.Y()

	return x, y
}

func LookAtTarget(o *engine.Object, cpx, cpy, dt float32) {
	pp := o.Shape.Body.Position()

	ca := AngleBetweenPoints(pp.X, pp.Y, cpx, cpy)
	pa := o.Shape.Body.Angle()
	// log.Println(pp.X, pp.Y, cpx, cpy, " | ", ca, pa)

	cv := vect.FromAngle(ca)
	pv := vect.FromAngle(pa)

	var rotspeed float32
	if o.Player != nil {
		rotspeed = o.Player.RotSpeed
	}
	if o.Bullet != nil {
		rotspeed = o.Bullet.Param.RotSpeed
	}

	vx := pv.X + (pv.X+cv.X)*(dt*rotspeed*0.1)
	vy := pv.Y + (pv.Y+cv.Y)*(dt*rotspeed*0.1)

	if o.Player != nil {
		subAngle := AngleBetweenAngles(pa, ca)

		if subAngle > -o.Player.SubAngle && o.Player.SubAngle > -1.5 {
			o.Player.SubAngle -= dt * o.Player.RotSpeed * 0.05
		}

		if subAngle < -o.Player.SubAngle && o.Player.SubAngle < 1.5 {
			o.Player.SubAngle += dt * o.Player.RotSpeed * 0.05
		}
	}

	//rotate
	o.Shape.Body.SetAngle(float32(math.Atan2(float64(vy), float64(vx))))
}
