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

var (
	Players []*Player
)

func CreateLocalPlayer(p *Player) {
	log.Println("CreateLocalPlayer")
	p.Local = true
	// p := &Player{Param: paramPlayer}
	p.CreateCursor(mgl32.Vec4{0.3, 0.3, 0.9, 0.7})
	p.CreatePlayer()

	p.Object.DestroyFunc = p.Destroy

	Players = append(Players, p)
	// p.Object.Shape.Body.CallBackCollision = p.Collision

	engine.AddCallback(p.Movement, p.PlayerRotation, p.CameraMovement, p.Attack)
	engine.SetMouseCallback(p.MouseControl)
}

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

	Cursor *engine.Object
}

func (p *Player) CreatePlayer() {
	p.CurrParam = p.InitParam
	p.Object.PI.Group = 2

	hb := NewHealthBar(p.InitParam.Health)
	p.Object.Create(hb)
	p.Object.MaxRollAngle = p.InitParam.RollAngle

	p.createWeapon(p.LeftWeapon)
	p.createWeapon(p.RightWeapon)

	if p.InitParam.MovSpeed > 5 {
		createTrail(p.Object, 0.5, int(p.InitParam.MovSpeed), mgl32.Vec2{1.4, 2.95})
		createTrail(p.Object, 0.5, int(p.InitParam.MovSpeed), mgl32.Vec2{1.4, -2.95})
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
		P:    point.Param{Size: point.P{1, 1, 0}},
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

func (p *Player) Movement(dt float32) bool {
	// log.Println(p.Object.Velocity().Length())
	if p.Object.Velocity().Length() < p.CurrParam.MovSpeed {
		dist := p.Object.Distance(p.Cursor)
		// log.Println(dist)
		// if dist > p.Param.MovSpeed {
		// 	dist = p.Param.MovSpeed
		// }

		p.Object.AddVelocity(p.Object.VectorForward(p.CurrParam.MovSpeed * 0.05 * dist * dt))
	}

	return true
}

func (p *Player) PlayerRotation(dt float32) bool {
	p.rotation(dt, p.Cursor.PositionVec2())

	angVel := p.Object.Shape.Body.AngularVelocity() / 2
	if angVel > p.Object.MaxRollAngle {
		angVel = p.Object.MaxRollAngle
	}
	if angVel < -p.Object.MaxRollAngle {
		angVel = -p.Object.MaxRollAngle
	}
	p.Object.RollAngle = -angVel

	return true
}

func (p *Player) rotation(dt float32, target mgl32.Vec2) float32 {
	angle := SubAngleObjectPoint(p.Object, target)

	if vect.FAbs(p.Object.Shape.Body.AngularVelocity()) < vect.FAbs(p.CurrParam.RotSpeed/10) {
		p.Object.Shape.Body.AddAngularVelocity(p.CurrParam.RotSpeed * 0.05 * angle * dt)
	}

	return angle
}

func (p *Player) Attack(dt float32) bool {
	if p.LeftWeapon != nil {
		p.LeftWeapon.Fire()
		p.WeaponDelay(p.LeftWeapon, "leftDelay")
	}

	if p.RightWeapon != nil {
		p.RightWeapon.Fire()
		p.WeaponDelay(p.RightWeapon, "rightDelay")
	}

	return true
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

	if button == 0 {
		p.LeftWeapon.shoot = action == 1
	}

	if button == 1 {
		p.RightWeapon.shoot = action == 1
	}

	if action == 1 {
		target := GetPlayerInPoint(p.Cursor.Position())
		if target != nil {
			log.Println(target.Name)
		}
		// object := engine.Hit(p.Cursor.Position())
		// if object != nil {
		// 	log.Println(object.Name)
		// }
	}
}

func (p *Player) CameraMovement(dt float32) bool {
	engine.Camera.SetPosition(p.Object.Position())

	// sun.Position = pp.Add(mgl32.Vec3{-30, 30, 100})

	x, y := engine.CursorPosition()
	w, h := engine.WindowSize()

	p.Cursor.SetPosition(getCursorPos(x, y, w, h, engine.Camera.GetPosition()))

	return true
}

// func keyboardControl(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
// 	log.Println("key", scancode)
// }

func NewHealthBar(value float32) *engine.Art {
	return &engine.Art{
		Name:     "health",
		Value:    value,
		MaxValue: value,
		P: point.Param{
			Pos:    point.P{-1, 1, 1.1},
			Size:   point.P{0.2, 2, 0},
			Static: true,
		},
		RI: &render.Instruction{
			MeshName: "plane",
			Material: &materials.Instruction{Name: "healthBar", DiffColor: mgl32.Vec4{0, 0.6, 0, 0.7}},
		},
	}
}
