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

type Player struct {
	Name   string
	Object *engine.Object

	Health, MovSpeed, RotSpeed, RollAngle float32
	LeftWeapon, RightWeapon               *Weapon

	Target       *Player
	targetAngle  float32
	respawnPoint mgl32.Vec2

	Cursor *engine.Object
}

func (p *Player) CreatePlayer() {
	hb := NewHealthBar(p.Health)
	p.Object.Create(hb)

	// aim := &engine.Art{
	// 	Name: "aim",
	// 	Param: engine.ObjectParam{
	// 		Node:     "plane",
	// 		Size:     engine.Point{1, 1, 1},
	// 		Pos:      engine.Point{15, 0, 1},
	// 		Material: engine.Material{Name: "aim", Texture: "cursor", Shader: "colortext2", DiffColor: mgl32.Vec4{1, 1, 1, 0.5}},
	// 	},
	// }

	// p.Object.AddArt(aim)

	if p.MovSpeed > 5 {
		createTrail(p.Object, 0.5, int(p.MovSpeed), mgl32.Vec2{1.4, 2.95})
		createTrail(p.Object, 0.5, int(p.MovSpeed), mgl32.Vec2{1.4, -2.95})
	}

	p.Object.DestroyFunc = p.Destroy
}

func CreateLocalPlayer(p *Player) {
	log.Println("CreateLocalPlayer")

	// p := &Player{Param: paramPlayer}
	p.CreateCursor(mgl32.Vec4{0.3, 0.3, 0.9, 0.7})
	p.CreatePlayer()

	Players = append(Players, p)
	// p.Object.Shape.Body.CallBackCollision = p.Collision

	engine.AddCallback(p.Movement, p.PlayerRotation, p.CameraMovement, p.Attack)
	engine.SetMouseCallback(p.MouseControl)
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
	x := 20 - rand.Float32()*40
	y := 20 - rand.Float32()*40
	p.Object.SetPosition(x, y)
	p.Object.SetVelocity(0, 0)
	p.Object.SetRotation(0)

	hp, ok := p.Object.GetArt("health")
	if !ok {
		log.Fatalln("helth bar not found", p.Object.Name)
	}

	hp.Value = hp.MaxValue
	hp.Resize()
}

func (p *Player) Movement(dt float32) bool {
	// log.Println(p.Object.Velocity().Length())
	if p.Object.Velocity().Length() < p.MovSpeed {
		dist := p.Object.Distance(p.Cursor)
		// log.Println(dist)
		// if dist > p.Param.MovSpeed {
		// 	dist = p.Param.MovSpeed
		// }

		p.Object.AddVelocity(p.Object.VectorForward(p.MovSpeed * 0.05 * dist * dt))
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

	if vect.FAbs(p.Object.Shape.Body.AngularVelocity()) < vect.FAbs(p.RotSpeed/10) {
		p.Object.Shape.Body.AddAngularVelocity(p.RotSpeed * 0.05 * angle * dt)
	}

	return angle
}

func (p *Player) Attack(dt float32) bool {
	if p.LeftWeapon != nil {
		p.Fire(p.LeftWeapon)
		p.WeaponDelay(p.LeftWeapon, "leftDelay")
	}

	if p.RightWeapon != nil {
		p.Fire(p.RightWeapon)
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

	delayBar, ok := p.Object.GetArt(name)
	if !ok {
		log.Printf("WARINING: art by name: %s not found", name)
		return
	}

	delayBar.Art.Body.Scale = mgl32.Vec3{1, value, 1}
}

func (p *Player) MouseControl(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {

	if button == 0 {
		p.LeftWeapon.Shoot = action == 1
	}

	if button == 1 {
		p.RightWeapon.Shoot = action == 1
	}

	if action == 1 {
		object := engine.Hit(p.Cursor.Position())
		if object != nil {
			log.Println(object.Name)
		}
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
			Pos:    point.P{0, 1, 1.1},
			Size:   point.P{2, 0.2, 0},
			Static: true,
		},
		RI: &render.Instruction{
			MeshName: "plane",
			Material: &materials.Instruction{Name: "healthBar", DiffColor: mgl32.Vec4{0, 0.6, 0, 0.7}},
		},
	}
}
