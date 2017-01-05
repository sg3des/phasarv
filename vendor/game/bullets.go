package game

import (
	"log"
	"math/rand"
	"phys"
	"point"
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"engine"
)

type Weapon struct {
	NextShot  time.Time
	Shoot     bool
	Delay     time.Duration
	DelayTime time.Time

	Bullet

	X float32

	AttackRate time.Duration
}

type Bullet struct {
	Type    string
	SubType string

	Object *engine.Object
	Weapon *Weapon
	Player *Player

	TargetPoint  mgl32.Vec2
	TargetObject *engine.Object

	RotSpeed  float32
	MovSpeed  float32
	Lifetime  time.Duration
	TimePoint time.Time

	Damage float32
	Shoot  bool
}

//Fire is main function for make shoot
func (p *Player) Fire(w *Weapon) {
	if !w.Shoot {
		w.DelayTime = time.Time{}
		return
	}

	if w.NextShot.After(time.Now()) {
		return
	}

	if w.Delay > 0 {

		if w.DelayTime.IsZero() {
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
func (p *Player) Shoot(w *Weapon) Bullet {
	// b := w.Bullet
	b := Bullet{
		Type:    w.Bullet.Type,
		SubType: w.Bullet.SubType,

		Object: &engine.Object{},
		Weapon: &Weapon{},
		Player: &Player{},

		// TargetPoint: mgl32.Vec2{},
		// TargetObject: &engine.Object{},

		RotSpeed:  w.Bullet.RotSpeed,
		MovSpeed:  w.Bullet.MovSpeed,
		Lifetime:  w.Bullet.Lifetime,
		TimePoint: time.Time{},

		Damage: w.Bullet.Damage,
		// Shoot  bool
	}

	// b.Object = &engine.Object{}
	*b.Object = *w.Bullet.Object
	*b.Weapon = *w
	*b.Player = *p
	b.TargetPoint = p.Cursor.PositionVec2()

	log.Println(&b.Object.RI, &w.Bullet.Object.RI)

	switch b.Type {
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
	vx, vy := b.Player.Object.VectorSide(1, -1.5704)
	x, y := b.Player.Object.Position()

	b.Object.P = point.Param{
		Pos:   point.P{x + b.Weapon.X*vx, y + b.Weapon.X*vy, 0},
		Angle: b.Player.Object.Rotation(),
	}

	b.Object.Create()
	// log.Println(b.Object.Shape, b.Object.Param.Phys)

	b.Object.Shape.Body.CallBackCollision = b.Collision

	b.TimePoint = time.Now().Add(b.Lifetime)
}

// INITIALIZE SHOOT from other types of weapons

//Gun create simple bullet
func (b *Bullet) Gun() {
	b.CreateObject()
	b.Object.AddCallback(b.GunCallback)
	// b.Object.Callback = b.GunCallback
	b.Object.SetVelocity(b.Player.Object.VectorForward(b.MovSpeed))

	b.Shoot = true
}

//Rocket create rocket bullet
func (b *Bullet) Rocket() {
	if b.SubType == "homing" {

		if b.TargetObject == nil {
			target := engine.Hit(b.Player.Cursor.Position())
			if target == nil || target.Name == "bullet" {
				b.Weapon.DelayTime = time.Time{}
				return
			}
			b.TargetObject = target
		}

	}

	b.CreateObject()
	b.Object.AddCallback(b.RocketCallback)
	b.Object.SetVelocity(b.Player.Object.VectorSide(b.Object.PI.Mass*5*b.Weapon.X, -1.5704))
	// b.Target = b.Player.Cursor.PositionVec2()

	createTrail(b.Object, 0.75, int(b.MovSpeed), mgl32.Vec2{-0.5})

	b.Shoot = true
	return
}

//Laser create laser line
func (b *Bullet) Laser() {
	//start position
	x, y := b.Player.Object.Position()
	wX, wY := b.Player.Object.VectorSide(b.Weapon.X, -1.5704)
	x += wX
	y += wY

	//length from lifetime
	h := float32(b.Lifetime.Seconds() * 10)
	tx, ty := b.Player.Object.VectorForward(h)
	tx += x
	ty += y

	if hit := engine.Raycast(x, y, tx, ty, 1, b.Player.Object.Shape.Body); hit != nil {
		h = hit.Distance

		if hit.Body.UserData != nil {
			ApplyDamage(hit.Body.UserData.(*engine.Object), b.Damage)
			// hit.Body.AddVelocity(b.Player.Object.VectorForward(b.Weapon.BulletObject.Phys.Mass * 10 / hit.Body.Mass()))
			hit.Body.AddAngularVelocity((rand.Float32() - 0.5) * 10 / hit.Body.Mass())
		}
	}

	b.Object.P.Size.Y = h
	b.Object.P.Pos = point.P{x, y, 0}
	b.Object.P.Angle = b.Player.Object.Rotation()
	b.Object.Create()

	// b.TimePoint = time.Now().Add(b.Lifetime)
	b.TimePoint = time.Now().Add(time.Second)
	b.Object.AddCallback(b.LaserCallback)

	b.Shoot = true
	return
}

//
//
//CALLBACKS

//GunCallback callback each frame
func (b *Bullet) GunCallback(dt float32) {
	if b.TimePoint.Before(time.Now()) {
		b.Destroy()
		// b.Object.Destroy()
		// b = nil
		return
	}

	b.Object.SetVelocity(b.Object.VectorForward(b.MovSpeed))
}

//RocketCallback callback each frame
func (b *Bullet) RocketCallback(dt float32) {
	if b.TimePoint.Before(time.Now()) {
		b.Destroy()
		return
	}

	dur, _ := time.ParseDuration("500ms")
	if b.Lifetime-b.TimePoint.Sub(time.Now()) < dur {
		return
	}

	var angle float32
	switch b.SubType {
	case "direct":
	case "aimed":
		angle = SubAngleObjectPoint(b.Object, b.TargetPoint)
	case "guided":
		angle = SubAngleObjectPoint(b.Object, b.Player.Cursor.PositionVec2())
	case "homing":
		tp := b.TargetPoint
		if b.TargetObject == nil {
			if target := engine.Hit(b.Player.Cursor.Position()); target != nil {
				b.TargetObject = target
				tp = target.PositionVec2()
			}
		} else {
			tp = b.TargetObject.PositionVec2()
		}
		angle = SubAngleObjectPoint(b.Object, mgl32.Vec2{tp.X(), tp.Y()})
	}

	if angle != 0 {
		b.Object.Shape.Body.AddAngularVelocity(angle * b.RotSpeed * 0.05 * dt)
	}

	if b.Object.Velocity().Length() < b.MovSpeed {
		b.Object.AddVelocity(b.Object.VectorForward(dt * b.MovSpeed))
	}

}

//LaserCallback each frame
func (b *Bullet) LaserCallback(dt float32) {
	if b.TimePoint.Before(time.Now()) {
		b.Destroy()
		return
	}
	// log.Println("laserCallbacl", dt)
	// if b.Object.Body == nil {
	// 	// log.Println("laser object node is nil - return")
	// 	return
	// }

	// color := b.Object.Body.Material.DiffuseColor

	// color[3] = color[3] - dt
	// b.Object.Body.Material.DiffuseColor = color
	// if color[3] <= 0.1 {
	// 	b.Object.Destroy()
	// }
}

//Destroy handler for bullet destroy
func (b *Bullet) Destroy() {
	switch b.Type {
	case "gun":
		b.Object.Destroy()
	case "rocket":
		b.Object.Destroy()
	case "laser":
		b.Object.Destroy()
		// b.Object.Shape.Body.Enabled = false
		// b.LaserCallback(0.1)
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

	if target.Name == "bullet" {
		return false
	}

	if _, ok := target.Arts["health"]; ok {
		ApplyDamage(target, b.Damage)
	}

	b.Destroy()

	return true
}

//ApplyDamage to object
func ApplyDamage(target *engine.Object, damage float32) {
	hp, ok := target.GetArt("health")
	if !ok {
		// log.Printf("WARINING: art by name: %s not found", "health")
		return
	}

	hp.Value -= damage
	if hp.Value <= 0 {
		target.Destroy()
		return
	}
	hp.Resize()

}