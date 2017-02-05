package game

import (
	"log"
	"phys"
	"point"
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"engine"
)

type Bullet struct {
	Object *engine.Object
	Weapon *Weapon
	Player *Player

	TargetPoint  mgl32.Vec2
	TargetPlayer *Player

	RotSpeed  float32
	MovSpeed  float32
	Lifetime  time.Duration
	TimePoint time.Time

	Damage float32
	Shoot  bool
}

//CreateObject create bullet for gun and rocket
func (b *Bullet) CreateObject() {
	// vx, vy := b.Player.Object.VectorSide(1, -1.5704)
	// x, y := b.Player.Object.Position()

	// log.Println(ang, b.Weapon.Angle)
	// if ang > b.Weapon.Angle {
	// 	ang = b.Weapon.Angle
	// } else if ang < -b.Weapon.Angle {
	// 	ang = -b.Weapon.Angle
	// }

	// log.Println(ang)

	// b.Object.P = point.Param{
	// 	Pos:   point.P{x + b.Weapon.X*vx, y + b.Weapon.X*vy, 0},
	// 	Angle: b.Player.Object.Rotation() + ang,
	// 	// Angle: ang,
	// }

	b.Object.P = point.Param{
		Pos:   point.PFromVect(b.Weapon.GetPosition()),
		Angle: b.Weapon.GetAngle(),
	}

	b.Object.Create()
	// log.Println(b.Object.Shape, b.Object.Param.Phys)

	b.Object.SetCallbackCollision(b.Collision)
	b.Object.SetUserData(b)

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
	if b.Weapon.SubType == Weapons.RocketType.Homing {

		if b.TargetPlayer == nil {
			target := GetPlayerInPoint(b.Player.Cursor.Position())
			if target == nil {
				b.Weapon.DelayTime = time.Time{}
			}
			b.TargetPlayer = target
			// target := engine.Hit(b.Player.Cursor.Position())
			// if target == nil || target.Name == "bullet" {
			// 	b.Weapon.DelayTime = time.Time{}
			// 	return
			// }
			// b.TargetObject = target
		}

	}

	b.CreateObject()
	b.Object.AddCallback(b.RocketCallback)
	b.Object.SetVelocity(b.Player.Object.VectorSide(b.Object.PI.Mass*5*b.Weapon.Pos.X, -1.5704))
	// b.Target = b.Player.Cursor.PositionVec2()

	createTrail(b.Object, 0.75, int(b.MovSpeed), mgl32.Vec2{-0.5})

	b.Shoot = true
	return
}

//Laser create laser line
func (b *Bullet) Laser() {
	//start position
	x, y := b.Player.Object.Position()
	wX, wY := b.Player.Object.VectorSide(b.Weapon.Pos.X, -1.5704)
	x += wX
	y += wY

	//length from lifetime
	h := float32(b.Lifetime.Seconds() * 10)
	tx, ty := b.Player.Object.VectorForward(h)
	tx += x
	ty += y

	if target := GetNearPlayerByRay(x, y, tx, ty, b.Player); target != nil {
		target.ApplyDamage(b.Damage)
	}

	// if hit := engine.Raycast(x, y, tx, ty, b.Player.Object.Shape.Body); hit != nil {
	// 	h = hit.Distance

	// 	if hit.Body.UserData != nil {

	// 		// ApplyDamage(hit.Body.UserData.(*engine.Object), b.Damage)
	// 		// // hit.Body.AddVelocity(b.Player.Object.VectorForward(b.Weapon.BulletObject.Phys.Mass * 10 / hit.Body.Mass()))
	// 		// hit.Body.AddAngularVelocity((rand.Float32() - 0.5) * 10 / hit.Body.Mass())
	// 	}
	// }

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
	switch b.Weapon.SubType {
	case Weapons.RocketType.Direct:
	case Weapons.RocketType.Aimed:
		angle = SubAngleObjectPoint(b.Object, b.TargetPoint)
	case Weapons.RocketType.Guided:
		angle = SubAngleObjectPoint(b.Object, b.Player.Cursor.PositionVec2())
	case Weapons.RocketType.Homing:
		tp := b.TargetPoint
		if b.TargetPlayer == nil {
			target := GetPlayerInPoint(b.Player.Cursor.Position())
			if target != nil {
				b.TargetPlayer = target
				tp = target.Object.PositionVec2()
			}

			// if target := engine.Hit(b.Player.Cursor.Position()); target != nil {
			// 	b.TargetObject = target
			// 	tp = target.PositionVec2()
			// }
		} else {
			tp = b.TargetPlayer.Object.PositionVec2()
		}
		angle = SubAngleObjectPoint(b.Object, mgl32.Vec2{tp.X(), tp.Y()})
	}

	if angle != 0 {
		b.Object.AddAngularVelocity(angle * b.RotSpeed * 0.05 * dt)
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
	switch b.Weapon.Type {
	case Weapons.Gun:
		b.Object.Destroy()
	case Weapons.Rocket:
		b.Object.Destroy()
	case Weapons.Laser:
		b.Object.Destroy()
		// b.Object.Shape.Body.Enabled = false
		// b.LaserCallback(0.1)
	}
}

//Collision event bullet collision
func (b *Bullet) Collision(arb *phys.Arbiter) bool {
	if arb.ShapeA.UserData == nil || arb.ShapeB.UserData == nil {
		//then one of objects is scene object
		b.Destroy()
		return true
	}

	p, b := resolveCollisionShapes(arb.ShapeA, arb.ShapeB)
	if p == nil {
		//then both shapes is bullets
		return false
	}
	if p == b.Player {
		//then player is parent of bullet
		return false
	}
	if b == nil {
		log.Fatalln("WTF?!? bullet not found")
	}

	p.ApplyDamage(b.Damage)

	b.Destroy()
	return true
}

func resolveCollisionShapes(shapes ...*phys.Shape) (p *Player, b *Bullet) {
	for _, shape := range shapes {
		if shape.UserData == nil {
			continue
		}
		switch shape.UserData.(type) {
		case *Bullet:
			b = shape.UserData.(*Bullet)
		case *Player:
			p = shape.UserData.(*Player)
		}
	}
	return
}

// //ApplyDamage to object
// func ApplyDamage(target *engine.Object, damage float32) {
// 	hp, ok := target.GetArt("health")
// 	if !ok {
// 		// log.Printf("WARINING: art by name: %s not found", "health")
// 		return
// 	}

// 	hp.Value -= damage
// 	if hp.Value <= 0 {
// 		target.Destroy()
// 		return
// 	}
// 	hp.Resize()

// }
