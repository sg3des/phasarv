package game

import (
	"engine"
	"materials"
	"phys/vect"
	"point"
	"render"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

var Weapons weaponTypes

type weaponTypes struct {
	Gun byte

	Rocket     byte
	RocketType struct {
		Direct byte
		Aimed  byte
		Guided byte
		Homing byte
	}

	Laser byte
}

func init() {
	Weapons.Gun = 'G'

	Weapons.Rocket = 'R'
	Weapons.RocketType.Direct = 'd'
	Weapons.RocketType.Aimed = 'a'
	Weapons.RocketType.Guided = 'g'
	Weapons.RocketType.Homing = 'h'

	Weapons.Laser = 'L'
}

type Weapon struct {
	Player *Player

	Type    byte
	SubType byte

	NextShot  time.Time
	ToShoot   bool
	Delay     time.Duration
	DelayTime time.Time

	Bullet

	Pos   vect.Vect
	Angle float32

	AttackRate time.Duration

	Turret *engine.Art
	Aim    *engine.Art
}

func (w *Weapon) Callback(dt float32) {
	if w.Aim != nil && w.Aim.Art != nil {
		w.Aim.Art.Angle = w.GetSubAngle()
	}
}

//Fire is main function for make shoot
func (w *Weapon) Fire() {
	if !w.ToShoot {
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

	b := w.Shoot()
	if b.Shoot {
		w.NextShot = time.Now().Add(w.AttackRate)
	}
}

//Shoot create new bullet
func (w *Weapon) Shoot() Bullet {
	b := w.Bullet
	b.Object = &engine.Object{}
	*b.Object = *w.Bullet.Object
	b.Weapon = w
	b.Player = w.Player
	b.TargetPoint = w.Player.Cursor.PositionVec2()

	switch w.Type {
	case Weapons.Gun:
		b.Gun()
	case Weapons.Rocket:
		b.Rocket()
	case Weapons.Laser:
		b.Laser()
	}

	return b
}

func (w *Weapon) GetPosition() vect.Vect {
	v1 := w.Player.Object.PositionVect()
	angle := w.Player.Object.Angle()

	return v1.SubPoint(angle, w.Pos)
}

func (w *Weapon) GetAngle() (ang float32) {
	return w.GetSubAngle() + w.Player.Object.Angle()
}

func (w *Weapon) GetSubAngle() (ang float32) {
	v1 := w.GetPosition()
	angle := w.Player.Object.Angle()

	ang = v1.SubAngle(angle, w.Player.Cursor.PositionVect())

	if ang > w.Angle {
		ang = w.Angle
	}
	if ang < -w.Angle {
		ang = -w.Angle
	}

	return ang
}

// case "direct":
// 	case "aimed":
// 		angle = SubAngleObjectPoint(b.Object, b.TargetPoint)
// 	case "guided":
// 		angle = SubAngleObjectPoint(b.Object, b.Player.Cursor.PositionVec2())
// 	case "homing":

// const ROCKET = "rocket"
// const ROCKET_HOMING
// const LASER = "laser"
// const GUN = "gun"

func (w *Weapon) NewAim() *engine.Art {
	ar := w.GetAttackRange()
	// switch w.Type {
	// case "rocket":
	// }
	return &engine.Art{
		Name:     "aim",
		Value:    ar,
		MaxValue: ar,
		P: point.Param{
			Pos:  point.PFromVect(w.Pos),
			Size: point.P{0.1, ar, 0},
		},
		RI: &render.Instruction{
			MeshName: "plane",
			Material: &materials.Instruction{Name: "aim", Texture: "laser", Shader: "colortext2", DiffColor: mgl32.Vec4{0.9, 0.9, 0.9, 0.5}},
		},
	}
}

func (w *Weapon) GetAttackRange() float32 {
	return float32(w.Bullet.Lifetime.Seconds()) * w.Bullet.MovSpeed
}
