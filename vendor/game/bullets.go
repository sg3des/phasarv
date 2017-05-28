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
	b.Object.P.Pos = point.PFromVect(b.Weapon.GetPosition())
	b.Object.P.Angle = b.Weapon.GetAngle()

	// b.Object.P = &point.Param{

	// 	Pos:   point.PFromVect(b.Weapon.GetPosition()),
	// 	Angle: b.Weapon.GetAngle(),
	// }

	b.Object.Create()

	b.Object.SetCallbackCollision(b.Collision)
	b.Object.SetUserData(b)

	b.TimePoint = time.Now().Add(b.Lifetime)
}

// INITIALIZE SHOOT from other types of weapons

//Gun create simple bullet
func (b *Bullet) Gun() {
	b.CreateObject()
	b.Object.AddCallback(b.GunCallback)
	b.Object.SetVelocity(b.Player.Object.VectorForward(b.MovSpeed))

	b.Shoot = true
}

//Rocket create rocket bullet
func (b *Bullet) Rocket() {
	if b.Weapon.SubType == Weapons.RocketType.Homing {

		if b.TargetPlayer == nil {
			b.TargetPlayer = GetPlayerInPoint(b.Player.Cursor.Position())
			if b.TargetPlayer == nil {
				b.Weapon.DelayTime = time.Time{}
			}
		}

	}

	b.CreateObject()
	b.Object.AddCallback(b.RocketCallback)
	b.Object.SetVelocity(b.Player.Object.VectorSide(b.Object.PI.Mass*5*b.Weapon.Pos.X, -1.5704))
	// b.Target = b.Player.Cursor.PositionVec2()

	createTrail(b.Object, 0.3, int(b.MovSpeed), mgl32.Vec2{-0.2, 0})

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

	o, dist := GetNearObjectByRay(x, y, tx, ty, b.Player.Object)
	if o != nil {
		if o.UserData != nil {
			o.UserData.(*Player).ApplyDamage(b.Damage)
		}
		h = dist
	}

	b.Object.P = &point.Param{
		Size:  point.P{h, b.Object.P.Size.Y, b.Object.P.Size.Z},
		Pos:   point.P{x, y, 0},
		Angle: b.Player.Object.Rotation(),
	}

	b.Object.Create()

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
	// if arb.ShapeA.UserData == nil || arb.ShapeB.UserData == nil {
	// 	//then one of objects is scene object
	// 	b.Destroy()
	// 	return true
	// }

	b, target := resolveCollisionShapes(arb.ShapeA, arb.ShapeB)
	// log.Println(b, p)
	if b == nil {
		log.Fatalln("WTF?!? bullet in collision not found")
	}

	if p, ok := target.UserData.(*Player); ok {
		if p == b.Player {
			return false
		}
		p.ApplyDamage(b.Damage)
	}
	b.Destroy()
	return true

	// if p == nil {
	// 	//then both shapes is bullets
	// 	return false
	// }
	// if p == b.Player {
	// 	//then player is parent of bullet
	// 	return false
	// }
	// if b == nil {
	// 	log.Fatalln("WTF?!? bullet not found")
	// }

	// p.ApplyDamage(b.Damage)

	// b.Destroy()
	// return true
}

func resolveCollisionShapes(shapes ...*phys.Shape) (b *Bullet, target *engine.Object) {
	for _, shape := range shapes {
		// if shape.UserData == nil {
		// 	continue
		// }
		o := shape.UserData.(*engine.Object)
		switch o.UserData.(type) {
		case *Bullet:

			b = o.UserData.(*Bullet)
		// case *Player:
		// 	p = o.UserData.(*Player)
		default:
			target = o
		}
	}
	return
}
