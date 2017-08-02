package weapons

import (
	"engine"
	"game/equip"
	"materials"
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

type Param struct {
	Pos vect.Vect

	Delay time.Duration
	Rate  time.Duration
	Range time.Duration

	Angle float32

	Ammunition int
	ReloadTime time.Duration
	ReloadCost float32

	Damage float32

	BulletMovSpeed float32
	BulletRotSpeed float32

	equip.Param
}

type Weapon struct {
	ShipObj *engine.Object

	Name      string
	Img       string
	Type      Type
	SubType   SubType
	EquipType equip.Type

	InitParam Param
	CurrParam Param
	BulletObj *engine.Object

	absAngle float32
	nextShot time.Time

	DelayTime time.Time
	ToShoot   bool

	Target    *engine.Object
	CursorPos mgl32.Vec2

	BulletCollisionCallback

	Turret *engine.Art
	Aim    *engine.Art
}

func (w *Weapon) Callback(dt float32) {
	w.update()

	if w.Aim != nil && w.Aim.Art != nil {
		w.Aim.Art.Angle = w.CurrParam.Angle
	}

}

func (w *Weapon) update() {
	w.absAngle = w.ShipObj.Rotation()
	w.CurrParam.Pos = w.ShipObj.PositionVect()
	w.CurrParam.Pos = w.CurrParam.Pos.SubPoint(w.absAngle-1.5704, w.InitParam.Pos)

	cPos := vect.FromVec2(w.CursorPos)

	w.CurrParam.Angle = w.CurrParam.Pos.SubAngle(w.absAngle, cPos)
	if w.CurrParam.Angle > w.InitParam.Angle {
		w.CurrParam.Angle = w.InitParam.Angle
	}
	if w.CurrParam.Angle < -w.InitParam.Angle {
		w.CurrParam.Angle = -w.InitParam.Angle
	}
	w.absAngle += w.CurrParam.Angle

	cPos.Sub(w.CurrParam.Pos)
	dist := cPos.Length()
	if ar := w.GetAttackRange(); dist > ar {
		dist = ar
	}

	av := vect.FromAngle(w.CurrParam.Angle)
	av.Mult(dist)
	av.Add(w.CurrParam.Pos)

	w.CursorPos = av.Vec2()
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
	ar := w.GetAttackRange()

	wX, wY := w.ShipObj.VectorSide(w.InitParam.Pos.X, -1.5704)

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

func (w *Weapon) GetAttackRange() (ar float32) {
	ar = float32(w.CurrParam.Range.Seconds())
	if w.CurrParam.BulletMovSpeed > 0 {
		if w.Type == Rocket {
			ar--
		}
		ar *= w.CurrParam.BulletMovSpeed
	}

	return
}
