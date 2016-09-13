package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"engine"
	"param"
	"phys"
)

type Bullet struct {
	Param  param.Bullet
	Object *engine.Object

	Weapon *param.Weapon
	Player *Player

	Shoot bool
}

//Fire is main function for make shoot
func (p *Player) Fire(w *param.Weapon) {
	if !w.Shoot {
		w.DelayTime = time.Time{}
		return
	}

	if w.NextShot.After(time.Now()) {
		return
	}

	if w.Delay > 0 {
		if timeNil(w.DelayTime) {
			w.DelayTime = time.Now().Add(w.Delay)
			return
		}

		if w.DelayTime.After(time.Now()) {
			//wait - no shot was
			return
		}
	}

	b := p.Shoot(w)
	if b.Shoot {
		w.NextShot = time.Now().Add(w.AttackRate)
	}
}

//Shoot create new bullet
func (p *Player) Shoot(w *param.Weapon) *Bullet {
	b := &Bullet{
		Player: p,
		Weapon: w,
		Param:  w.BulletParam,
	}
	// *b.Param = w.BulletParam
	switch b.Param.Type {
	case "gun":
		b.Gun()
	case "rocket":
		b.Rocket()
	case "laser":
		b.Laser()
	}
	return b
}

//CreateObject create bullet for gun and rocket
func (b *Bullet) CreateObject() {
	vx, vy := b.Player.Object.VectorSide(1)
	x, y := b.Player.Object.Position()

	b.Object = engine.NewObject(b.Weapon.BulletObject)
	b.Object.SetPosition(x+b.Weapon.X*vx, y+b.Weapon.X*vy)
	b.Object.SetRotation(b.Player.Object.Rotation())

	b.Object.Shape.Body.CallBackCollision = b.Collision

	b.Param.TimePoint = time.Now().Add(b.Param.Lifetime)
}

// INITIALIZE SHOOT from other types of weapons

//Gun create simple bullet
func (b *Bullet) Gun() {
	b.CreateObject()
	b.Object.Callback = b.GunCallback
	b.Object.SetVelocity(b.Player.Object.VectorForward(b.Param.MovSpeed))
	b.Shoot = true
}

//Rocket create rocket bullet
func (b *Bullet) Rocket() {
	if b.Param.SubType == "homing" {
		target := engine.Hit(cursor.Position())
		if target == nil {
			b.Weapon.DelayTime = time.Time{}
			return
		}

		b.Param.TargetObject = target
	}

	b.CreateObject()
	b.Object.Callback = b.RocketCallback
	b.Object.SetVelocity(b.Player.Object.VectorSide(b.Weapon.BulletObject.PH.Mass * 5 * b.Weapon.X))
	b.Param.Target = cursor.PositionVec2()

	b.Shoot = true
	return
}

//Laser create laser line
func (b *Bullet) Laser() {
	//start position
	x, y := b.Player.Object.Position()
	wX, wY := b.Player.Object.VectorSide(b.Weapon.X)
	x += wX
	y += wY

	//length from lifetime
	h := float32(b.Param.Lifetime.Seconds() * 10)
	tx, ty := b.Player.Object.VectorForward(h)
	tx += x
	ty += y

	if hit := engine.Raycast(x, y, tx, ty, 1, b.Player.Object.Shape.Body); hit != nil {
		h = hit.Distance

		if hit.Body.UserData != nil {
			ApplyDamage(hit.Body.UserData.(*engine.Object), b.Param.Damage)
			hit.Body.AddVelocity(b.Player.Object.VectorForward(b.Weapon.BulletObject.PH.Mass * 10 / hit.Body.Mass()))
			hit.Body.AddAngularVelocity((rand.Float32() - 0.5) * 10 / hit.Body.Mass())
		}
	}

	//draw laser
	b.Object = engine.NewPlane(b.Weapon.BulletObject, b.Weapon.BulletObject.PH.W, h)
	b.Object.SetPosition(x, y)
	b.Object.SetRotation(b.Player.Object.Rotation())

	b.Param.TimePoint = time.Now().Add(b.Param.Lifetime)
	b.Object.Callback = b.LaserCallback

	b.Shoot = true
	return
}

//
//
//CALLBACKS

//GunCallback callback each frame
func (b *Bullet) GunCallback(ob *engine.Object, dt float32) {
	if b.Param.TimePoint.Before(time.Now()) {
		b.Destroy()
		// b.Object.Destroy()
		// b = nil
		return
	}

	b.Object.SetVelocity(b.Object.VectorForward(b.Param.MovSpeed))
}

//RocketCallback callback each frame
func (b *Bullet) RocketCallback(ob *engine.Object, dt float32) {
	if b.Param.TimePoint.Before(time.Now()) {
		b.Destroy()
		return
	}

	dur, _ := time.ParseDuration("500ms")
	if b.Param.Lifetime-b.Param.TimePoint.Sub(time.Now()) < dur {
		return
	}

	var angle float32
	switch b.Param.SubType {
	case "direct":
	case "aimed":
		tp := b.Param.Target
		angle = AngleObjectPoint(ob, mgl32.Vec2{tp.X(), tp.Y()})
	case "guided":
		angle = AngleObjectPoint(ob, cursor.PositionVec2())
	case "homing":
		tp := b.Param.Target
		if b.Param.TargetObject == nil {
			if target := engine.Hit(cursor.Position()); target != nil {
				b.Param.TargetObject = target
				tp = target.PositionVec2()
			}
		} else {
			tp = b.Param.TargetObject.(*engine.Object).PositionVec2()
		}
		angle = AngleObjectPoint(ob, mgl32.Vec2{tp.X(), tp.Y()})
	}

	if angle != 0 {
		ob.Shape.Body.SetAngularVelocity(angle * b.Param.RotSpeed * 0.1)
	}

	if ob.Velocity().Length() < b.Param.MovSpeed {
		ob.AddVelocity(ob.VectorForward(dt * b.Param.MovSpeed))
	}
}

//LaserCallback each frame
func (b *Bullet) LaserCallback(ob *engine.Object, dt float32) {
	color := b.Object.Node.Core.DiffuseColor
	// color[4] =
	color[3] = color[3] - dt
	b.Object.Node.Core.DiffuseColor = color
	if color[3] <= 0.1 {
		b.Object.Destroy()
	}
}

//Destroy handler for bullet destroy
func (b *Bullet) Destroy() {
	switch b.Param.Type {
	case "gun":
		b.Object.Destroy()
	case "rocket":
		b.Object.Destroy()
	case "laser":
		b.Object.Shape.Body.Enabled = false
		b.LaserCallback(b.Object, 0.1)
	}
}

//Collision event bullet collision
func (b *Bullet) Collision(arb *phys.Arbiter) bool {
	if arb.BodyA.UserData == nil || arb.BodyB.UserData == nil {
		log.Println("unset bodies")
		return true
	}

	var target *engine.Object

	if arb.BodyA.UserData.(*engine.Object) == b.Object {
		target = arb.BodyB.UserData.(*engine.Object)
	} else if arb.BodyB.UserData.(*engine.Object) == b.Object {
		target = arb.BodyA.UserData.(*engine.Object)
	} else {
		log.Println("WTF?!")
		return false
	}

	if b.Player.Object == target {
		return false
	}

	if _, ok := target.ArtStatic["health"]; ok {
		ApplyDamage(target, b.Param.Damage)
	}

	b.Destroy()

	return true
}

//ApplyDamage to object
func ApplyDamage(target *engine.Object, damage float32) {
	hp, ok := target.GetArt("health")
	if !ok {
		log.Printf("WARINING: art by name: %s not found", "health")
		return
	}

	hp.Value -= damage
	if hp.Value <= 0 {
		ObjectDestroy(target)
		return
	}
	hp.Resize()

}

//ObjectDestroy need add explosion
func ObjectDestroy(o *engine.Object) {
	log.Println("not yet ready")
	o.Destroy()
}
