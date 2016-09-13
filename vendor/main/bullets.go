package main

import (
	"log"
	"math/rand"
	"time"

	"engine"
	"param"
	"phys"
)

//check if time is nil
func timeNil(t time.Time) bool {
	if t.Equal(time.Time{}) {
		return true
	}
	return false
}

//Fire is main function for make shoot
func Fire(w *param.Weapon, p *engine.Object) {
	if !w.Shoot {
		w.DelayTime = time.Time{}
		return
	}

	if w.NextShot.After(time.Now()) {
		return
	}

	var shoot bool
	switch w.BulletParam.Type {
	case "gun":
		shoot = gun(w, p)
	case "rocket":
		shoot = rocket(w, p)
	case "laser":
		shoot = laser(w, p)
	}

	if shoot {
		w.NextShot = time.Now().Add(w.AttackRate)
	}
}

// INITIALIZE SHOOT from other types of weapons

//create simple bullet from gun
func gun(w *param.Weapon, p *engine.Object) bool {
	ob := createBullet(w, p)
	// w.NextShot = time.Now().Add(w.AttackRate)
	ob.Callback = gunCallback
	ob.SetVelocity(p.VectorForward(w.BulletParam.MovSpeed))
	return true
}

//create rocket
func rocket(w *param.Weapon, p *engine.Object) bool {
	var ob *engine.Object

	if w.BulletParam.SubType == "homing" {
		target := engine.Hit(cursor.Position())
		if target == nil {
			w.DelayTime = time.Time{}
			return false
		} else {
			w.BulletParam.TargetObject = target
		}
	}

	if w.Delay > 0 {
		if timeNil(w.DelayTime) {
			w.DelayTime = time.Now().Add(w.Delay)
		}

		if w.DelayTime.After(time.Now()) {
			//wait - no shot was
			return false
		}
		w.DelayTime = time.Time{}
	}

	ob = createBullet(w, p)
	ob.Callback = rocketCallback
	ob.SetVelocity(p.VectorSide(w.BulletObject.PH.Mass * 5 * w.X))
	ob.Bullet.Param.Target = cursor.PositionVec2()

	return true
}

//create bullet for gun and rocket
func createBullet(w *param.Weapon, p *engine.Object) (ob *engine.Object) {
	vx, vy := p.VectorSide(1)
	x, y := p.Position()

	ob = engine.NewObject(w.BulletObject, nil)
	ob.SetPosition(x+w.X*vx, y+w.X*vy)
	ob.SetRotation(p.Rotation())

	ob.Shape.Body.CallBackCollision = BulletCollision

	ob.Bullet = &engine.Bullet{
		Parent: p,
		Param:  &param.Bullet{},
	}
	*ob.Bullet.Param = w.BulletParam //copy bullet parameters
	ob.Bullet.Param.TimePoint = time.Now().Add(w.BulletParam.Lifetime)

	return
}

var tbox *engine.Object

//create laser
func laser(w *param.Weapon, p *engine.Object) bool {
	if w.Delay > 0 {
		if timeNil(w.DelayTime) {
			w.DelayTime = time.Now().Add(w.Delay)
		}

		if w.DelayTime.After(time.Now()) {
			return false
		}
		w.DelayTime = time.Time{}
	}

	//start position
	x, y := p.Position()
	wX, wY := p.VectorSide(w.X)
	x += wX
	y += wY

	//length from lifetime
	h := float32(w.BulletParam.Lifetime.Seconds() * 10)
	tx, ty := p.VectorForward(h)
	tx += x
	ty += y

	if hit := engine.Raycast(x, y, tx, ty, 1, p.Shape.Body); hit != nil {
		h = hit.Distance
		tx, ty = p.VectorForward(h)
		tx += x
		ty += y

		if hit.Body.UserData != nil {
			ApplyDamage(hit.Body.UserData.(*engine.Object), w.BulletParam.Damage)
			hit.Body.AddVelocity(p.VectorForward(w.BulletObject.PH.Mass * 10 / hit.Body.Mass()))
			hit.Body.AddAngularVelocity((rand.Float32() - 0.5) * 10 / hit.Body.Mass())
		}
	}

	//draw laser
	laser := engine.NewPlane(w.BulletObject, w.BulletObject.PH.W, h)
	laser.SetPosition(x, y)
	laser.SetRotation(p.Rotation())

	laser.Bullet = &engine.Bullet{
		Parent: p,
		Param:  &param.Bullet{},
	}
	*laser.Bullet.Param = w.BulletParam //copy bullet parameters
	laser.Bullet.Param.TimePoint = time.Now().Add(w.BulletParam.Lifetime)

	laser.Callback = laserCallback

	// w.NextShot = time.Now().Add(w.AttackRate)

	return true
}

//
//
//CALLBACKS

func gunCallback(ob *engine.Object, dt float32) {
	if ob.Bullet.Param.TimePoint.Before(time.Now()) {
		ob.Destroy()
	}

	ob.SetVelocity(ob.VectorForward(ob.Bullet.Param.MovSpeed))
}

func rocketCallback(ob *engine.Object, dt float32) {
	if ob.Bullet.Param.TimePoint.Before(time.Now()) {
		BulletDestroy(ob)
	}

	dur, _ := time.ParseDuration("500ms")
	if ob.Bullet.Param.Lifetime-ob.Bullet.Param.TimePoint.Sub(time.Now()) < dur {
		return
	}
	// if   {

	// }

	switch ob.Bullet.Param.SubType {
	case "direct":
	case "aimed":
		tp := ob.Bullet.Param.Target
		LookAtTarget(ob, tp.X(), tp.Y(), dt)
	case "guided":
		cpx, cpy := cursor.Position()
		LookAtTarget(ob, cpx, cpy, dt)
	case "homing":
		tp := ob.Bullet.Param.Target
		if ob.Bullet.Param.TargetObject == nil {
			if target := engine.Hit(cursor.Position()); target != nil {
				ob.Bullet.Param.TargetObject = target
				tp = target.PositionVec2()
			}
		} else {
			tp = ob.Bullet.Param.TargetObject.(*engine.Object).PositionVec2()
		}
		LookAtTarget(ob, tp.X(), tp.Y(), dt)
	}

	if ob.Velocity().Length() < ob.Bullet.Param.MovSpeed {
		ob.AddVelocity(ob.VectorForward(dt * ob.Bullet.Param.MovSpeed))
	}
}

func laserCallback(ob *engine.Object, dt float32) {
	// if ob.Bullet.Param.TimePoint.Before(time.Now()) {

	// }

	color := ob.Node.Core.DiffuseColor
	// color[4] =
	color[3] = color[3] - dt
	ob.Node.Core.DiffuseColor = color
	if color[3] <= 0.1 {
		ob.Destroy()
	}
}

// handler for bullet destroy
func BulletDestroy(ob *engine.Object) {
	switch ob.Bullet.Param.Type {
	case "gun":
		ob.Destroy()
	case "rocket":
		ob.Destroy()
	case "laser":
		ob.Shape.Body.Enabled = false
		laserCallback(ob, 0.1)
	}
}

// event bullet collision
func BulletCollision(arb *phys.Arbiter) bool {
	// log.Println("bullet Collision")
	if arb.BodyA.UserData == nil || arb.BodyB.UserData == nil {
		log.Println("unset bodies")
		return true
	}

	var bullet *engine.Object
	var target *engine.Object

	if arb.BodyA.UserData.(*engine.Object).Name == "bullet" {
		bullet = arb.BodyA.UserData.(*engine.Object)
		target = arb.BodyB.UserData.(*engine.Object)
	} else {
		bullet = arb.BodyB.UserData.(*engine.Object)
		target = arb.BodyA.UserData.(*engine.Object)
	}

	if bullet.Bullet != nil && bullet.Bullet.Parent != nil && bullet.Bullet.Parent == target {
		return false
	}

	if _, ok := target.ArtStatic["health"]; ok {
		if bullet.Bullet != nil {
			ApplyDamage(target, bullet.Bullet.Param.Damage)
		} else {
			log.Println("damage unknown")
			ApplyDamage(target, 10)
		}

	}

	BulletDestroy(bullet)

	return true
}

func ApplyDamage(target *engine.Object, damage float32) {
	healthbar, ok := target.ArtStatic["health"]
	if !ok {
		return
	}

	healthbar.Value -= damage
	if healthbar.Value <= 0 {
		ObjectDestroy(target)
	}

	healthbar.Resize()
}

func ObjectDestroy(o *engine.Object) {
	log.Println("not yet ready")
	o.Destroy()
}
