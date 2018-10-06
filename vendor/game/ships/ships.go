package ships

import (
	"engine"
	"game/equip"
	"game/weapons"
	"log"
	"materials"
	"math/rand"
	"phys"
	"phys/vect"
	"point"
	"render"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Type string

var Fighter Type = "interceptor"

type Ship struct {
	Object *engine.Object

	Name  string
	Class string
	Img   string
	Mesh  string
	Type  Type
	Size  mgl32.Vec3

	InitParam equip.Param
	CurrParam equip.Param
	Slots     []equip.Slot

	LeftWpnPos, RightWpnPos mgl32.Vec3
	LeftWeapon, RightWeapon *weapons.Weapon
	Equipment               []*equip.Equip

	Cursor       *engine.Object
	CursorOffset vect.Vect //only for network

	fires []*engine.Art

	Local bool
}

func (s *Ship) Create() {
	s.CurrParam = s.InitParam
	// log.Printf("%+v : %+v", s.CurrParam, s.InitParam)
	// prm := s.InitParam

	s.Object = &engine.Object{
		Name: s.Name,
		P: &point.Param{
			Pos:  point.PFromVec3(s.InitParam.Pos),
			Size: point.PFromVec3(s.Size),
		},
		PI: &phys.Instruction{
			Mass:      s.InitParam.Weight,
			Group:     phys.GROUP_PLAYER,
			ShapeType: phys.ShapeType_Box,
		},
		RI: &render.Instruction{
			MeshName: s.Mesh,
			Material: &materials.Instruction{Name: "player", Texture: "TestCube", Shader: "basic", SpecLevel: 1},
			Shadow:   true,
		},
	}
	// p.Object.PI.Group = 2

	if engine.NeedRender {
		var hb *engine.Art
		hb = NewHealthBar(s.InitParam.Health)
		s.Object.Create(hb)
	}

	s.Object.MaxRollAngle = s.InitParam.RollAngle

	s.createWeapon(s.LeftWeapon, s.LeftWpnPos)
	s.createWeapon(s.RightWeapon, s.RightWpnPos)

	// if s.Cursor == nil {
	// s.Cursor = &engine.Object{}
	// s.Cursor.Create()
	// }

	if s.InitParam.MovSpeed > 5 && engine.NeedRender {
		lf := NewEngineFire(point.P{-0.8, 0.3, 0.1})
		s.Object.AppendArt(lf)
		rf := NewEngineFire(point.P{-0.8, -0.3, 0.1})
		s.Object.AppendArt(rf)

		s.fires = []*engine.Art{lf, rf}

		s.Object.AddTrail(mgl32.Vec3{-1.2, 0.3, 0}, int(s.InitParam.MovSpeed)*2, point.P{0.7, 0.3, 1}, 1)
		s.Object.AddTrail(mgl32.Vec3{-1.2, -0.3, 0}, int(s.InitParam.MovSpeed)*2, point.P{0.7, 0.3, 1}, 1)

		// p.trails = []*render.Particle{lt, rt}
	}

	s.Object.SetUserData(s)
}

func (s *Ship) createWeapon(w *weapons.Weapon, pos mgl32.Vec3) {
	if w == nil {
		return
	}

	w.InitParam.Pos = pos
	w.CurrParam = w.InitParam
	w.ShipObj = s.Object
	w.SetBulletCollisionCallback(s.BulletCollision)

	if s.Local {
		w.Aim = w.NewAim()
		s.Object.AppendArt(w.Aim)
	}

	// s.Object.AddCallback(w.Callback)
}

func (s *Ship) CreateCursor(color mgl32.Vec4) {
	cursor := &engine.Object{
		Name: "cursor",
		P:    &point.Param{Size: point.P{1, 1, 0}},
		RI: &render.Instruction{
			MeshName: "plane",
			Material: &materials.Instruction{Name: "cursor", Shader: "colortext2", DiffColor: color},
		},
	}

	cursor.Create()
	s.Cursor = cursor
}

func (s *Ship) Collision(arb *phys.Arbiter) bool {
	log.Println("ship Collision")

	return true
}

func (s *Ship) Destroy() {
	// log.Println("Destroy", p.Name)

	s.Object.SetPosition(engine.GetRandomPoint(20, 20))
	s.Object.SetVelocity(0, 0)
	s.Object.SetRotation(0)

	s.CurrParam = s.InitParam
	s.updateArt("health", s.CurrParam.Health)

	// hp, ok := p.Object.GetArt("health")
	// if !ok {
	// 	log.Fatalln("helth bar not found", p.Object.Name)
	// }

	// hp.Value = hp.MaxValue
	// hp.Resize()
}

func (s *Ship) ApplyDamage(damage float32) {
	// log.Println("ApplyDamage", p.Name, damage)

	s.CurrParam.Health -= damage
	s.updateArt("health", s.CurrParam.Health)
	if s.CurrParam.Health <= 0 {
		s.Object.Destroy()
	}
}

func (s *Ship) updateArt(name string, value float32) {
	if art := s.Object.GetArt(name); art != nil {
		art.Resize(value)
		return
	}
	log.Printf("warning: art by name: `%s` not found", name)
}

//ClientCursor is update cursor position on server side by cursor offset
func (s *Ship) ClientCursor(dt float32) {
	pos := s.Object.PositionVect()
	pos.Add(s.CursorOffset)
	s.Cursor.SetPosition(pos.X, pos.Y)
}

func (s *Ship) Attack(dt float32) {
	if s.LeftWeapon != nil {
		s.LeftWeapon.Fire()
		s.WeaponDelay(s.LeftWeapon, "leftDelay")
	}

	if s.RightWeapon != nil {
		s.RightWeapon.Fire()
		s.WeaponDelay(s.RightWeapon, "rightDelay")
	}
}

func (s *Ship) BulletCollision(b *weapons.Bullet, target *engine.Object) bool {
	if target == nil {
		return false
	}
	if target == s.Object {
		return false
	}

	return true
}

func (s *Ship) WeaponDelay(w *weapons.Weapon, name string) {
	if w.CurrParam.Delay == 0 {
		return
	}

	var value float32
	if w.DelayTime.IsZero() {
		value = 1
	} else {
		value = float32(w.DelayTime.Sub(time.Now()).Seconds())
		if value < 0 {
			value = 0
		}
		value = value / float32(w.CurrParam.Delay.Seconds())
	}

	delayBar := s.Object.GetArt(name)
	if delayBar == nil {
		log.Printf("WARINING: art by name: %s not found", name)
		return
	}

	delayBar.Art.Body.Scale = mgl32.Vec3{1, value, 1}
}

func (s *Ship) Rotation(dt float32) {
	s.Rotate(dt, s.Cursor.PositionVec2())

	angVel := s.Object.AngularVelocity() / 2
	if angVel > s.Object.MaxRollAngle {
		angVel = s.Object.MaxRollAngle
	}
	if angVel < -s.Object.MaxRollAngle {
		angVel = -s.Object.MaxRollAngle
	}
	s.Object.RollAngle = -angVel
}

func (s *Ship) Rotate(dt float32, target mgl32.Vec2) float32 {
	angle := s.Object.SubAngleObjectPoint(target)

	if vect.FAbs(s.Object.AngularVelocity()) < vect.FAbs(s.CurrParam.RotSpeed/10) {
		s.Object.AddAngularVelocity(s.CurrParam.RotSpeed * 0.05 * angle * dt)
	}

	return angle
}

func (s *Ship) Movement(dt float32) {
	speed := s.Object.Velocity().Length()

	// log.Println(speed, s.CurrParam.MovSpeed)

	if speed < s.CurrParam.MovSpeed {
		dist := s.Object.Distance(s.Cursor)
		s.Object.AddVelocity(s.Object.VectorForward(s.CurrParam.MovSpeed * 0.05 * dist * dt))
	}

	if engine.NeedRender {
		for _, f := range s.fires {
			if f.Art.Body != nil {
				scale := 0.5 + speed*0.2 + rand.Float32()
				f.Art.Body.Scale = mgl32.Vec3{scale, 0.8 + scale*0.1, scale}
			}
		}
	}
}

func (p *Ship) MouseControl(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {

	switch {
	case button == 0 && p.LeftWeapon != nil:
		p.LeftWeapon.ToShoot = action == 1
	case button == 1 && p.RightWeapon != nil:
		p.RightWeapon.ToShoot = action == 1
	}
}

// func keyboardControl(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
// 	log.Println("key", scancode)
// }

func NewEngineFire(pos point.P) *engine.Art {
	fire := &engine.Art{
		Name:     "fire",
		MaxValue: 10,
		P: &point.Param{
			Pos:   pos,
			Size:  point.P{1, 1, 1},
			Angle: 3.14159,
		},
		RI: &render.Instruction{
			MeshName: render.MESH_VECTOR,
			Material: &materials.Instruction{
				Name:      "fire",
				Texture:   "fire",
				Shader:    "blend",
				DiffColor: mgl32.Vec4{0.9, 0.9, 0.9, 0.9},
			},
			Transparent: true,
		},
	}

	return fire
}

func NewHealthBar(value float32) *engine.Art {
	return &engine.Art{
		Name:     "health",
		Value:    value,
		MaxValue: value,
		P: &point.Param{
			Pos:    point.P{-1, 1, 1.1},
			Size:   point.P{2, 0.2, 0},
			Static: true,
		},
		RI: &render.Instruction{
			MeshName: render.MESH_VECTOR,
			Material: &materials.Instruction{
				Name:      "healthBar",
				DiffColor: mgl32.Vec4{0, 0.6, 0, 0.7},
			},
		},
	}
}
