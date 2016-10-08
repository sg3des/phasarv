package main

import (
	"engine"
	"log"
	"param"
	"phys"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	Param  *param.Player
	Object *engine.Object
}

func (p *Player) CreateLocalPlayer() {
	p.Object = engine.NewObject(p.Param.Object, engine.NewHealthBar(p.Param.Health))

	aim := param.Art{Name: "aim", Type: param.ArtRotate, W: 1, H: 1, LocalPos: mgl32.Vec3{15, 0, 1}, Material: param.Material{Name: "aim", Texture: "cursor", Shader: "colortext2", DiffColor: mgl32.Vec4{1, 1, 1, 0.5}}}
	p.Object.NewArt(aim)

	weapondelay := param.Art{MaxValue: 100, Type: param.ArtRotate, W: 0.1, H: 1.5, LocalPos: mgl32.Vec3{0, -2, 1}, Material: param.Material{Name: "weapondelay", Shader: "colortext2", DiffColor: mgl32.Vec4{1, 1, 1, 0.5}}}
	if p.Param.LeftWeapon.Delay > 0 {
		weapondelay.Name = "leftDelay"
		p.Object.NewArt(weapondelay)
	}
	if p.Param.RightWeapon.Delay > 0 {
		weapondelay.Name = "rightDelay"
		weapondelay.H = -1.5
		weapondelay.LocalPos = mgl32.Vec3{0, 2, 1}
		p.Object.NewArt(weapondelay)
	}

}

func initLocalPlayer() {
	var p = &Player{}
	p.Param = &param.Player{
		Name: "player0",
		Object: param.Object{
			Name:     "player",
			Mesh:     param.Mesh{Model: "trapeze", Shadow: true},
			Material: param.Material{Name: "player", Texture: "TestCube", Shader: "diffuse_texbumped_shadows", SpecLevel: 1},
			Phys:     &param.Phys{W: 3, H: 2, Mass: 12, Group: 1, Type: phys.ShapeType_Box},
		},
		Health:   100,
		MovSpeed: 10,
		RotSpeed: 15,

		// LeftWeapon: &param.Weapon{
		// 	BulletParam: param.Bullet{
		// 		Type:     "gun",
		// 		MovSpeed: 20,
		// 		Lifetime: 10000 * time.Millisecond,
		// 		Damage:   20,
		// 	},
		// 	BulletObject: param.Object{
		// 		Name: "bullet",
		// 		Mesh: param.Mesh{Model: "bullet", Texture: "TestCube", Shader: "diffuse"},
		// 		Phys:   param.Phys{W: 0.1, H: 0.1, Mass: 1},
		// 	},
		// 	X:          -1,
		// 	AttackRate: 200 * time.Millisecond,
		// },
		LeftWeapon: &param.Weapon{
			BulletParam: param.Bullet{
				Type:     "laser",
				Lifetime: 2500 * time.Millisecond,
				Damage:   50,
			},
			BulletObject: param.Object{
				Name:        "bullet",
				Mesh:        param.Mesh{Model: "bullet"},
				Material:    param.Material{Name: "laser", Texture: "laser", Shader: "blend"},
				Phys:        &param.Phys{W: 0.5, Mass: 0.5},
				Transparent: true,
			},

			X: -1,
			// Delay:      500 * time.Millisecond,
			AttackRate: 50 * time.Millisecond,
		},
		RightWeapon: &param.Weapon{
			BulletParam: param.Bullet{
				Type:     "rocket",
				SubType:  "guided",
				MovSpeed: 20,
				RotSpeed: 15,
				Lifetime: 5000 * time.Millisecond,
				Damage:   30,
			},
			BulletObject: param.Object{
				Name:     "bullet",
				Mesh:     param.Mesh{Model: "bullet"},
				Material: param.Material{Name: "bullet", Texture: "gray", Shader: "color"},
				Phys:     &param.Phys{W: 0.1, H: 0.1, Mass: 0.5},
			},
			X:          1,
			Delay:      500 * time.Millisecond,
			AttackRate: 1000 * time.Millisecond,
		},
	}

	p.CreateLocalPlayer()

	engine.AddCallback(p.Movement, p.Rotation, p.CameraMovement, p.Attack)
	engine.SetMouseCallback(p.MouseControl)
	// engine.Window.SetKeyCallback(keyboardControl)
}

func (p *Player) Movement(dt float32) {
	if p.Object.Velocity().Length() < p.Param.MovSpeed {
		dist := p.Object.Distance(cursor)
		if dist > p.Param.MovSpeed {
			dist = p.Param.MovSpeed
		}

		p.Object.AddVelocity(p.Object.VectorForward(p.Param.MovSpeed * 0.0001 * dist))
	}
}

func (p *Player) Attack(dt float32) {
	if p.Param.LeftWeapon != nil {
		p.Fire(p.Param.LeftWeapon)
		p.WeaponDelay(p.Param.LeftWeapon, "leftDelay")
	}

	if p.Param.RightWeapon != nil {
		p.Fire(p.Param.RightWeapon)
		p.WeaponDelay(p.Param.RightWeapon, "rightDelay")
	}
}

func (p *Player) WeaponDelay(w *param.Weapon, name string) {
	if w.Delay == 0 {
		return
	}

	var value float32
	if timeNil(w.DelayTime) {
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

func (p *Player) Rotation(dt float32) {
	angle := AngleObjectPoint(p.Object, cursor.PositionVec2())

	if angle > -p.Param.RotSpeed && angle < p.Param.RotSpeed {
		p.Object.Shape.Body.AddAngularVelocity(angle * p.Param.RotSpeed * 0.0001)
	}

	angVel := p.Object.Shape.Body.AngularVelocity()
	if angVel > engine.MaxRollAngle {
		angVel = engine.MaxRollAngle
	}
	if angVel < -engine.MaxRollAngle {
		angVel = -engine.MaxRollAngle
	}
	p.Object.RollAngle = -angVel
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

func (p *Player) CameraMovement(dt float32) {
	pp := p.Object.Node.Location

	cp := camera.GetPosition()
	camera.SetPosition(pp.X(), pp.Y(), cp.Z())

	// light0.Position = mgl32.Vec3{pp.X() - 5, pp.Y() + 5, 30}

	x, y := engine.CursorPosition()
	w, h := engine.WindowSize()

	x, y = getCursorPos(x, y, w, h, cp)

	cursor.Node.Location = mgl32.Vec3{x, y, 0}
	// log.Println(x, y, xfloat, yfloat)
}

// func keyboardControl(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
// 	log.Println("key", scancode)
// }
