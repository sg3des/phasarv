package game

import (
	"engine"
	"log"
	"materials"
	"math/rand"
	"phys"
	"phys/vect"
	"point"
	"render"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type PlayerParam struct {
	Health, Shield                float32
	Energy, Metal                 float32
	MovSpeed, RotSpeed, RollAngle float32
}

type Player struct {
	Name   string
	Local  bool
	Object *engine.Object

	CurrParam PlayerParam
	InitParam PlayerParam

	LeftWeapon, RightWeapon *Weapon

	Target *Player

	targetAngle  float32
	respawnPoint mgl32.Vec2

	Cursor       *engine.Object
	CursorOffset vect.Vect //only for network

	trails []*render.Particle
}

func (p *Player) CreatePlayer() {
	p.CurrParam = p.InitParam
	// p.Object.PI.Group = 2

	var hb *engine.Art
	if NeedRender {
		hb = NewHealthBar(p.InitParam.Health)
	}
	p.Object.Create(hb)

	p.Object.MaxRollAngle = p.InitParam.RollAngle

	p.createWeapon(p.LeftWeapon)
	p.createWeapon(p.RightWeapon)

	if p.Cursor == nil {
		p.Cursor = &engine.Object{}
		p.Cursor.Create()
	}

	if p.InitParam.MovSpeed > 5 && NeedRender {
		lt := p.Object.AddTrail(mgl32.Vec3{1.2, 2.8, 0}, int(p.InitParam.MovSpeed), point.P{1, 0.2, 1}, 1)
		rt := p.Object.AddTrail(mgl32.Vec3{1.2, -2.8, 0}, int(p.InitParam.MovSpeed), point.P{1, 0.2, 1}, 1)

		p.trails = []*render.Particle{lt, rt}
	}

	p.Object.SetUserData(p)
}

func (p *Player) createWeapon(w *Weapon) {
	if w == nil {
		return
	}
	w.Player = p

	if p.Local {
		w.Aim = w.NewAim()
		p.Object.AppendArt(w.Aim)
	}

	p.Object.AddCallback(w.Callback)
}

func (p *Player) CreateCursor(color mgl32.Vec4) {
	cursor := &engine.Object{
		Name: "cursor",
		P:    &point.Param{Size: point.P{1, 1, 0}},
		RI: &render.Instruction{
			MeshName: "plane",
			Material: &materials.Instruction{Name: "cursor", Shader: "colortext2", DiffColor: color},
		},
	}

	cursor.Create()
	p.Cursor = cursor
}

func (p *Player) Collision(arb *phys.Arbiter) bool {
	// var player *engine.Object

	// if arb.BodyA.UserData.(*engine.Object) == p.Object {
	// 	player = arb.BodyA.UserData.(*engine.Object)
	// } else if arb.BodyB.UserData.(*engine.Object) == p.Object {
	// 	player = arb.BodyB.UserData.(*engine.Object)
	// } else {
	// 	log.Println("WTF?!")
	// 	return true
	// }

	return true
}

func (p *Player) Destroy() {
	// log.Println("Destroy", p.Name)

	p.Object.SetPosition(GetRandomPoint(20, 20))
	p.Object.SetVelocity(0, 0)
	p.Object.SetRotation(0)

	p.CurrParam = p.InitParam
	p.updateArt("health", p.CurrParam.Health)

	// hp, ok := p.Object.GetArt("health")
	// if !ok {
	// 	log.Fatalln("helth bar not found", p.Object.Name)
	// }

	// hp.Value = hp.MaxValue
	// hp.Resize()
}

func GetRandomPoint(x, y float32) (float32, float32) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return x/2 - r.Float32()*x, y/2 - r.Float32()*y
}

func (p *Player) ApplyDamage(damage float32) {
	// log.Println("ApplyDamage", p.Name, damage)

	p.CurrParam.Health -= damage
	p.updateArt("health", p.CurrParam.Health)
	if p.CurrParam.Health <= 0 {
		p.Object.Destroy()
	}
}

func (p *Player) updateArt(name string, value float32) {
	if art := p.Object.GetArt(name); art != nil {
		art.Resize(value)
		return
	}
	log.Printf("warning: art by name: `%s` not found", name)
}

func (p *Player) Movement(dt float32) {
	// log.Println(p.Object.Velocity().Length())
	if p.Object.Velocity().Length() < p.CurrParam.MovSpeed {
		dist := p.Object.Distance(p.Cursor)

		// log.Println(dist, p.Object.Velocity().Length())

		scale := dist - p.Object.Velocity().Length()
		log.Println(scale)
		if scale < 0 {
			scale = 0
		} else if scale > p.CurrParam.MovSpeed {
			scale = p.CurrParam.MovSpeed
		}

		if NeedRender {
			for _, trail := range p.trails {
				trail.Scale = scale / 16
			}
		}
		// log.Println(dist)
		// if dist > p.Param.MovSpeed {
		// 	dist = p.Param.MovSpeed
		// }

		p.Object.AddVelocity(p.Object.VectorForward(p.CurrParam.MovSpeed * 0.05 * dist * dt))
	}
}

func (p *Player) PlayerRotation(dt float32) {
	p.rotation(dt, p.Cursor.PositionVec2())

	angVel := p.Object.AngularVelocity() / 2
	if angVel > p.Object.MaxRollAngle {
		angVel = p.Object.MaxRollAngle
	}
	if angVel < -p.Object.MaxRollAngle {
		angVel = -p.Object.MaxRollAngle
	}
	p.Object.RollAngle = -angVel
}

func (p *Player) rotation(dt float32, target mgl32.Vec2) float32 {
	angle := SubAngleObjectPoint(p.Object, target)

	if vect.FAbs(p.Object.AngularVelocity()) < vect.FAbs(p.CurrParam.RotSpeed/10) {
		p.Object.AddAngularVelocity(p.CurrParam.RotSpeed * 0.05 * angle * dt)
	}

	return angle
}

//ClientCursor is update cursor position on server side by cursor offset
func (p *Player) ClientCursor(dt float32) {
	pos := p.Object.PositionVect()
	pos.Add(p.CursorOffset)
	p.Cursor.SetPosition(pos.X, pos.Y)
}

func (p *Player) Attack(dt float32) {
	if p.LeftWeapon != nil {
		p.LeftWeapon.Fire()
		p.WeaponDelay(p.LeftWeapon, "leftDelay")
	}

	if p.RightWeapon != nil {
		p.RightWeapon.Fire()
		p.WeaponDelay(p.RightWeapon, "rightDelay")
	}
}

func (p *Player) WeaponDelay(w *Weapon, name string) {
	if w.Delay == 0 {
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
		value = value / float32(w.Delay.Seconds())
	}

	delayBar := p.Object.GetArt(name)
	if delayBar == nil {
		log.Printf("WARINING: art by name: %s not found", name)
		return
	}

	delayBar.Art.Body.Scale = mgl32.Vec3{1, value, 1}
}

func (p *Player) MouseControl(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {

	switch button {
	case 0:
		p.LeftWeapon.ToShoot = action == 1
	case 1:
		p.RightWeapon.ToShoot = action == 1
	}

	// if action == 1 {
	// 	target := GetPlayerInPoint(p.Cursor.Position())
	// 	if target != nil {
	// 		log.Println(target.Name)
	// 	}
	// }
}

func (p *Player) CameraMovement(dt float32) {
	render.SetCameraPosition(p.Object.Position())

	// sun.Position = pp.Add(mgl32.Vec3{-30, 30, 100})

	x, y := engine.CursorPosition()
	w, h := engine.WindowSize()

	p.Cursor.SetPosition(getCursorPos(x, y, w, h, render.GetCameraPosition()))
}

// func keyboardControl(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
// 	log.Println("key", scancode)
// }

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
			MeshName: "vector",
			Material: &materials.Instruction{Name: "healthBar", DiffColor: mgl32.Vec4{0, 0.6, 0, 0.7}},
		},
	}
}
