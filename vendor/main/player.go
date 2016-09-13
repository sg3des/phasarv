package main

import (
	"engine"
	"log"
	"param"
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

	p.Object.NewArt(param.Art{Name: "aim", Texture: "cursor", Shader: "colortext2", Color: mgl32.Vec4{1, 1, 1, 0.5}, Type: param.ArtRotate, W: 1, H: 1, LocalPos: mgl32.Vec3{15, 0, 1}})

	if p.Param.LeftWeapon.Delay > 0 {
		p.Object.NewArt(param.Art{Name: "leftDelay", Shader: "colortext2", Color: mgl32.Vec4{1, 1, 1, 0.5}, MaxValue: 100, Type: param.ArtRotate, W: 0.1, H: 1.5, LocalPos: mgl32.Vec3{0, -2, 1}})
	}
	if p.Param.RightWeapon.Delay > 0 {
		p.Object.NewArt(param.Art{Name: "rightDelay", Shader: "colortext2", Color: mgl32.Vec4{1, 1, 1, 0.5}, MaxValue: 100, Type: param.ArtRotate, W: 0.1, H: -1.5, LocalPos: mgl32.Vec3{0, 2, 1}})
	}

}

func initLocalPlayer() {
	var p = &Player{}
	p.Param = &param.Player{
		Name: "player0",
		Object: param.Object{
			Name: "player",
			Mesh: param.Mesh{Model: "trapeze", Texture: "TestCube", Shader: "diffuse_texbumped_shadows", Shadow: true},
			PH:   param.Phys{W: 2, H: 2, Mass: 12, Group: 1},
		},

		Health:   100,
		MovSpeed: 15,
		RotSpeed: 20,

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
		// 		PH:   param.Phys{W: 0.1, H: 0.1, Mass: 1},
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
				Mesh:        param.Mesh{Model: "bullet", Texture: "laser", Shader: "blend"},
				PH:          param.Phys{W: 0.5, Mass: 0.5},
				Transparent: true,
			},

			X:          -1,
			Delay:      500 * time.Millisecond,
			AttackRate: 1000 * time.Millisecond,
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
				Name: "bullet",
				Mesh: param.Mesh{Model: "bullet", Texture: "TestCube", Shader: "diffuse"},
				PH:   param.Phys{W: 0.1, H: 0.1, Mass: 0.5},
			},
			X:          1,
			Delay:      500 * time.Millisecond,
			AttackRate: 1000 * time.Millisecond,
		},
	}

	p.CreateLocalPlayer()

	engine.AddCallback(p.Movement, p.Rotation, p.CameraMovement, p.Attack)
	engine.Window.SetMouseButtonCallback(p.MouseControl)
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

	if p.Object.RollAngle >= 0 {
		if angle < 0 && p.Object.RollAngle < 1.5 {
			p.Object.RollAngle -= p.Param.RotSpeed * dt * 0.05 * angle
		}

		if angle > 0 {
			p.Object.RollAngle -= p.Param.RotSpeed * dt * 0.05
		}
	} else {
		if angle > 0 && p.Object.RollAngle > -1.5 {
			p.Object.RollAngle -= p.Param.RotSpeed * dt * 0.05 * angle
		}

		if angle < 0 {
			p.Object.RollAngle += p.Param.RotSpeed * dt * 0.05
		}
	}
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

	light0.Position = mgl32.Vec3{pp.X() - 5, pp.Y() + 5, 30}

	xfloat, yfloat := engine.Window.GetCursorPos()
	width, height := engine.Window.GetSize()

	x, y := getCursorPos(float32(xfloat), float32(yfloat), width, height, cp)

	cursor.Node.Location = mgl32.Vec3{x, y, 0}
	// log.Println(x, y, xfloat, yfloat)
}

// func keyboardControl(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
// 	log.Println("key", scancode)
// }
