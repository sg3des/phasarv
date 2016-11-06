package main

import (
	"engine"
	"log"
	"math/rand"
	"param"
	"phys"
	"phys/vect"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	Param  param.Player
	Object *engine.Object

	Target *Player

	targetAngle  float32
	respawnPoint mgl32.Vec2
}

func (p *Player) CreatePlayer() {
	p.Object = engine.NewObject(p.Param.Object, engine.NewHealthBar(p.Param.Health))

	aim := param.Art{Name: "aim", Type: param.ArtRotate, W: 1, H: 1, LocalPos: mgl32.Vec3{15, 0, 1}, Material: param.Material{Name: "aim", Texture: "cursor", Shader: "colortext2", DiffColor: mgl32.Vec4{1, 1, 1, 0.5}}}
	p.Object.NewArt(aim)

	weapondelay := param.Art{MaxValue: 100, Type: param.ArtRotate, W: 0.1, Material: param.Material{Name: "weapondelay", Shader: "colortext2", DiffColor: mgl32.Vec4{1, 1, 1, 0.5}}}

	if p.Param.RightWeapon != nil && p.Param.LeftWeapon.Delay > 0 {
		weapondelay.Name = "leftDelay"
		weapondelay.H = 1.5
		weapondelay.LocalPos = mgl32.Vec3{0, -2, 1}
		p.Object.NewArt(weapondelay)
	}

	if p.Param.RightWeapon != nil && p.Param.RightWeapon.Delay > 0 {
		weapondelay.Name = "rightDelay"
		weapondelay.H = -1.5
		weapondelay.LocalPos = mgl32.Vec3{0, 2, 1}
		p.Object.NewArt(weapondelay)
	}

	if p.Param.MovSpeed > 5 {
		createTrail(p.Object, 0.5, int(p.Param.MovSpeed), mgl32.Vec2{1.4, 2.95})
		createTrail(p.Object, 0.5, int(p.Param.MovSpeed), mgl32.Vec2{1.4, -2.95})
	}

	p.Object.DestroyFunc = p.Destroy
	// engine.NewParticles()
}

func CreateLocalPlayer(paramPlayer param.Player) {
	log.Println("CreateLocalPlayer")

	p := &Player{Param: paramPlayer}
	p.CreatePlayer()

	players = append(players, p)
	// p.Object.Shape.Body.CallBackCollision = p.Collision

	engine.AddCallback(p.Movement, p.PlayerRotation, p.CameraMovement, p.Attack)
	engine.SetMouseCallback(p.MouseControl)
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
	if p.Object.Velocity().Length() < p.Param.MovSpeed {
		dist := p.Object.Distance(cursor)
		// log.Println(dist)
		// if dist > p.Param.MovSpeed {
		// 	dist = p.Param.MovSpeed
		// }

		p.Object.AddVelocity(p.Object.VectorForward(p.Param.MovSpeed * 0.05 * dist * dt))
	}

	return true
}

func (p *Player) PlayerRotation(dt float32) bool {
	p.rotation(dt, cursor.PositionVec2())

	angVel := p.Object.Shape.Body.AngularVelocity() / 2
	if angVel > p.Param.Object.MaxRollAngle {
		angVel = p.Param.Object.MaxRollAngle
	}
	if angVel < -p.Param.Object.MaxRollAngle {
		angVel = -p.Param.Object.MaxRollAngle
	}
	p.Object.RollAngle = -angVel

	return true
}

func (p *Player) rotation(dt float32, target mgl32.Vec2) float32 {
	angle := SubAngleObjectPoint(p.Object, target)

	if vect.FAbs(p.Object.Shape.Body.AngularVelocity()) < vect.FAbs(p.Param.RotSpeed/10) {
		p.Object.Shape.Body.AddAngularVelocity(p.Param.RotSpeed * 0.05 * angle * dt)
	}

	return angle
}

func (p *Player) Attack(dt float32) bool {
	if p.Param.LeftWeapon != nil {
		p.Fire(p.Param.LeftWeapon)
		p.WeaponDelay(p.Param.LeftWeapon, "leftDelay")
	}

	if p.Param.RightWeapon != nil {
		p.Fire(p.Param.RightWeapon)
		p.WeaponDelay(p.Param.RightWeapon, "rightDelay")
	}

	return true
}

func (p *Player) WeaponDelay(w *param.Weapon, name string) {
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

	delayBar.Art.Scale = mgl32.Vec3{1, value, 1}
}

func (p *Player) MouseControl(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {

	if button == 0 {
		p.Param.LeftWeapon.Shoot = action == 1
	}

	if button == 1 {
		p.Param.RightWeapon.Shoot = action == 1
	}

	if action == 1 {
		object := engine.Hit(cursor.Position())
		if object != nil {
			log.Println(object.Name)
		}
	}
}

func (p *Player) CameraMovement(dt float32) bool {
	pp := p.Object.Node.Location

	cp := camera.GetPosition()
	camera.SetPosition(pp.X(), pp.Y(), cp.Z())

	// sun.Position = pp.Add(mgl32.Vec3{-30, 30, 100})

	x, y := engine.CursorPosition()
	w, h := engine.WindowSize()

	cursor.SetPosition(getCursorPos(x, y, w, h, cp))

	return true
}

// func keyboardControl(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
// 	log.Println("key", scancode)
// }
