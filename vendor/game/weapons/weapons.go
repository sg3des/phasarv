package weapons

import (
	"engine"
	"game/equip"
	"materials"
	"phys"
	"phys/vect"
	"point"
	"render"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type Type byte
type SubType byte

var Gun Type = 'G'
var Rocket Type = 'R'
var Laser Type = 'L'

var TypeDirect SubType = 'd'
var TypeAimed SubType = 'a'
var TypeGuided SubType = 'g'
var TypeHoming SubType = 'h'

type Weapon struct {
	ShipObj *engine.Object

	equip.Equip

	Type    Type
	SubType SubType

	BulletObj *engine.Object

	absAngle float32
	nextShot time.Time

	DelayTime time.Time
	ToShoot   bool

	Target    *engine.Object
	CursorPos mgl32.Vec2

	bulletCollisionCallback BulletCollisionCallback

	Turret *engine.Art
	Aim    *engine.Art
}

func (w *Weapon) SetBulletCollisionCallback(f BulletCollisionCallback) {
	w.bulletCollisionCallback = f
}

func (w *Weapon) UpdateCursor(x, y float32) {
	w.absAngle = w.ShipObj.Rotation()
	wpnpos := w.ShipObj.PositionVect()
	wpnpos = wpnpos.SubPoint(w.absAngle-1.5704, vect.FromVec3(w.InitParam.Pos))
	w.CurrParam.Pos = wpnpos.Vec3()

	cPos := vect.Vect{x, y}

	w.CurrParam.Angle = wpnpos.SubAngle(w.absAngle, cPos)
	if w.CurrParam.Angle > 0 && w.CurrParam.Angle > w.InitParam.Angle {
		w.CurrParam.Angle = w.InitParam.Angle
	}
	if w.CurrParam.Angle < 0 && w.CurrParam.Angle < w.InitParam.Angle*-1 {
		w.CurrParam.Angle = w.InitParam.Angle * -1
	}
	w.absAngle += w.CurrParam.Angle

	cPos.Sub(wpnpos)
	dist := cPos.Length()
	if ar := w.GetAttackRange(w.CurrParam); dist > ar {
		dist = ar
	}

	av := vect.FromAngle(w.absAngle)
	av.Mult(dist)
	av.Add(wpnpos)

	w.CursorPos = av.Vec2()

	if w.Aim != nil && w.Aim.Art != nil {
		w.Aim.Art.Angle = w.CurrParam.Angle
	}
}

//Fire is main function for make shoot
func (w *Weapon) Fire() {
	if !w.ToShoot {
		w.DelayTime = time.Time{}
		return
	}

	if w.nextShot.After(time.Now()) {
		return
	}

	if w.CurrParam.Delay > 0 {

		if w.DelayTime.IsZero() {
			w.DelayTime = time.Now().Add(w.CurrParam.Delay)
			return
		}

		if w.DelayTime.After(time.Now()) {
			//wait - no shot was
			return
		}
	}

	b := w.Shoot()
	if b.Shoot {
		w.nextShot = time.Now().Add(w.CurrParam.Rate)
	}
}

//Shoot create new bullet
func (w *Weapon) Shoot() *Bullet {
	b := &Bullet{
		Object:      new(engine.Object),
		Weapon:      w,
		TargetPoint: w.CursorPos,
		RotSpeed:    w.CurrParam.BulletRotSpeed,
		MovSpeed:    w.CurrParam.BulletMovSpeed,
		Lifetime:    w.CurrParam.Range,
		Damage:      w.CurrParam.Damage,
	}

	*b.Object = *w.BulletObj
	if b.Object.PI != nil {
		b.Object.PI.Group = phys.GROUP_BULLET
	}

	switch w.Type {
	case Gun:
		b.Gun()
	case Rocket:
		b.Rocket()
	case Laser:
		b.Laser()
	}

	return b
}

// func (w *Weapon) GetAngle() (ang float32) {
// 	return w.GetSubAngle() + w.ShipObj.Rotation()
// }

// func (w *Weapon) GetSubAngle() (ang float32) {
// 	v1 := w.GetPosition()
// 	angle := w.ShipObj.Rotation()

// 	ang = v1.SubAngle(angle, vect.FromVec2(w.CursorPos))

// 	if ang > w.Angle {
// 		ang = w.Angle
// 	}
// 	if ang < -w.Angle {
// 		ang = -w.Angle
// 	}

// 	return ang
// }

func (w *Weapon) NewAim() *engine.Art {
	ar := w.GetAttackRange(w.InitParam)

	wX, wY := w.ShipObj.VectorSide(w.InitParam.Pos[0], -1.5704)

	return &engine.Art{
		Name:     "aim",
		Value:    ar,
		MaxValue: ar,
		P: &point.Param{
			Pos:  point.P{wX, wY, 0},
			Size: point.P{ar, 0.1, 0},
		},
		RI: &render.Instruction{
			MeshName: render.MESH_VECTOR,
			Material: &materials.Instruction{
				Name:      "aim",
				Texture:   "laser",
				Shader:    "colortext2",
				DiffColor: mgl32.Vec4{0.9, 0.9, 0.9, 0.5},
			},
		},
	}
}

func (w *Weapon) GetAttackRange(p equip.Param) (ar float32) {
	ar = float32(p.Range.Seconds())
	if p.BulletMovSpeed > 0 {
		if w.Type == Rocket {
			ar--
		}
		ar *= p.BulletMovSpeed
	}

	return
}
