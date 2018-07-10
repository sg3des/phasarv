package weapons

import (
	"log"
	"phys"
	"phys/vect"
	"point"
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"engine"
)

type BulletCollisionCallback func(b *Bullet, target *engine.Object) bool

type Bullet struct {
	Object *engine.Object
	Weapon *Weapon

	TargetPoint mgl32.Vec2
	Target      *engine.Object

	RotSpeed  float32
	MovSpeed  float32
	Lifetime  time.Duration
	TimePoint time.Time

	Damage float32
	Shoot  bool
}

//CreateObject create bullet for gun and rocket
func (b *Bullet) CreateObject() {
	b.Object.P.Pos = point.PFromVec3(b.Weapon.CurrParam.Pos)
	b.Object.P.Angle = b.Weapon.absAngle

	b.Object.Create()

	b.Object.SetCallbackCollision(b.collision)
	b.Object.SetUserData(b)

	b.TimePoint = time.Now().Add(b.Lifetime)
}

// INITIALIZE SHOOT from other types of weapons

//Gun create simple bullet
func (b *Bullet) Gun() {
	b.CreateObject()
	b.Object.AddCallback(b.gunCallback)
	b.Object.SetVelocity(b.Weapon.ShipObj.VectorForward(b.MovSpeed))

	b.Shoot = true
}

//Rocket create rocket bullet
func (b *Bullet) Rocket() {
	if b.Weapon.SubType == TypeHoming {
		if b.Target == nil {
			b.Target = engine.GetObjectInPoint(b.Weapon.CursorPos)
			if b.Target == nil {
				b.Weapon.DelayTime = time.Time{}
			}
		}
	}

	b.CreateObject()
	b.Object.AddCallback(b.rocketCallback)
	b.Object.SetVelocity(b.Weapon.ShipObj.VectorSide(b.Object.PI.Mass*5*b.Weapon.InitParam.Pos[0], -1.5704))

	b.Object.AddTrail(mgl32.Vec3{-0.2, 0, 0}, int(b.MovSpeed), point.P{0.3, 0.2, 1}, 1)

	b.Shoot = true
	return
}

//Laser create laser line
func (b *Bullet) Laser() {
	h := b.Weapon.GetAttackRange(b.Weapon.CurrParam)
	x, y, _ := b.Weapon.CurrParam.Pos.Elem()
	// tx, ty := b.Weapon.CursorPos.Elem()

	av := vect.FromAngle(b.Weapon.absAngle)
	av.Mult(h)
	av.Add(vect.FromVec2(b.Weapon.CursorPos))
	tx, ty := av.Elem()

	if o, dist := b.Weapon.ShipObj.GetNearObjectByRay(x, y, tx, ty); o != nil {
		log.Println(o, dist)
		if b.Weapon.bulletCollisionCallback != nil {
			b.Weapon.bulletCollisionCallback(b, o)
		}
		h = dist
	}

	b.Object.P = &point.Param{
		Size:  point.P{h, b.Object.P.Size.Y, b.Object.P.Size.Z},
		Pos:   point.P{x, y, 0},
		Angle: b.Weapon.absAngle,
	}

	b.Object.Create()

	b.TimePoint = time.Now().Add(time.Second)
	b.Object.AddCallback(b.laserCallback)

	b.Shoot = true
	return
}

//
//
//CALLBACKS

//gunCallback callback each frame
func (b *Bullet) gunCallback(dt float32) {
	if b.TimePoint.Before(time.Now()) {
		b.Destroy()
		return
	}

	b.Object.SetVelocity(b.Object.VectorForward(b.MovSpeed))
}

//rocketCallback callback each frame
func (b *Bullet) rocketCallback(dt float32) {
	if b.TimePoint.Before(time.Now()) {
		b.Destroy()
		return
	}

	if b.Lifetime-b.TimePoint.Sub(time.Now()) < 5e8 { //500ms
		return
	}

	var angle float32
	switch b.Weapon.SubType {
	case TypeDirect:
	case TypeAimed:
		angle = b.Object.SubAngleObjectPoint(b.TargetPoint)
	case TypeGuided:
		angle = b.Object.SubAngleObjectPoint(b.Weapon.CursorPos)
	case TypeHoming:
		tp := b.TargetPoint
		if b.Target == nil {
			target := engine.GetObjectInPoint(b.Weapon.CursorPos)
			if target != nil {
				b.Target = target
				tp = target.PositionVec2()
			}

		} else {
			tp = b.Target.PositionVec2()
		}
		angle = b.Object.SubAngleObjectPoint(mgl32.Vec2{tp.X(), tp.Y()})
	}

	if angle != 0 {
		b.Object.AddAngularVelocity(angle * b.RotSpeed * 0.05 * dt)
	}

	if b.Object.Velocity().Length() < b.MovSpeed {
		b.Object.AddVelocity(b.Object.VectorForward(dt * b.MovSpeed))
	}
}

//collision event bullet collision
func (b *Bullet) collision(arb *phys.Arbiter) bool {
	b, target := resolveCollisionShapes(arb.ShapeA, arb.ShapeB)
	if b == nil {
		log.Fatalln("WTF?!? bullet in collision not found")
	}
	if target == nil {
		return false
	}

	if b.Weapon.bulletCollisionCallback != nil {
		des := b.Weapon.bulletCollisionCallback(b, target)
		if des {
			b.Destroy()
		}
		return des
	}
	return false

	// if p, ok := target.UserData.(*Ship); ok {
	// 	if p == b.Ship {
	// 		return false
	// 	}
	// 	p.ApplyDamage(b.Damage)
	// }
	// b.Destroy()
	// return true
}

func resolveCollisionShapes(shapes ...*phys.Shape) (b *Bullet, target *engine.Object) {
	for _, shape := range shapes {
		o := shape.UserData.(*engine.Object)
		switch o.UserData.(type) {
		case *Bullet:
			b = o.UserData.(*Bullet)
		default:
			target = o
		}
	}
	return
}

//LaserCallback each frame
func (b *Bullet) laserCallback(dt float32) {
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
	case Gun:
		b.Object.Destroy()
	case Rocket:
		b.Object.Destroy()
	case Laser:
		b.Object.Destroy()
		// b.Object.Shape.Body.Enabled = false
		// b.LaserCallback(0.1)
	}
}
