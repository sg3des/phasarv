package main

import (
	"engine"
	"log"
	"param"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	localPlayer *engine.Object
)

func newPlayer() {
	p := &param.Player{
		Name: "player0",
		Object: param.Object{
			Name: "player",
			Mesh: param.Mesh{Model: "trapeze", Texture: "TestCube", Shader: "diffuse_texbumped_shadows", Shadow: true},
			PH:   param.Phys{W: 2, H: 2, Mass: 12, Group: 1},
		},

		Health:   100,
		MovSpeed: 5,
		RotSpeed: 10,

		// LeftWeapon: &param.Weapon{
		// 	BulletParam: param.Bullet{
		// 		Type:     "gun",
		// 		MovSpeed: 20,
		// 		Lifetime: 10000 * time.Millisecond,
		// 		Damage:   20,
		// 	},
		// 	BulletObject: param.Object{
		// 		Name: "bullet",
		// 		Mesh: param.Mesh{Model: "bullet", Texture: "TestCube", Shader: "diffuse_texbumped"},
		// 		PH:   param.Phys{W: 0.1, H: 0.1, Mass: 1},
		// 	},
		// 	X:          -1,
		// 	AttackRate: 200 * time.Millisecond,
		// },
		LeftWeapon: &param.Weapon{
			BulletParam: param.Bullet{
				Type:     "laser",
				Lifetime: 2500 * time.Millisecond,
				Damage:   18,
			},
			BulletObject: param.Object{
				Name:        "bullet",
				Mesh:        param.Mesh{Model: "bullet", Texture: "laser", Shader: "blend"},
				PH:          param.Phys{W: 0.5, Mass: 0.5},
				Transparent: true,
			},
			// Delay:      300 * time.Millisecond,
			X:          -1,
			AttackRate: 50 * time.Millisecond,
		},
		RightWeapon: &param.Weapon{
			BulletParam: param.Bullet{
				Type:     "rocket",
				SubType:  "direct",
				MovSpeed: 20,
				RotSpeed: 30,
				Lifetime: 5000 * time.Millisecond,
				Damage:   30,
			},
			BulletObject: param.Object{
				Name: "bullet",
				Mesh: param.Mesh{Model: "bullet", Texture: "TestCube", Shader: "diffuse"},
				PH:   param.Phys{W: 0.1, H: 0.1, Mass: 0.5},
			},
			X:          1,
			Delay:      1200 * time.Millisecond,
			AttackRate: 500 * time.Millisecond,
		},
	}

	localPlayer = engine.NewObject(p.Object, engine.NewHealthBar(p.Health))
	localPlayer.Player = p
	localPlayer.SetPosition(0, 0)

	rightCircle := param.Art{Name: "right", W: 16, H: -1, Seg: 1024, MaxValue: 100, Color: mgl32.Vec4{1, 1, 1, 0.5}, Type: param.ArtRotate, Line: true}
	localPlayer.NewCurve(rightCircle)

	leftCircle := param.Art{Name: "left", W: 16, H: 1, Seg: 1024, MaxValue: 100, Color: mgl32.Vec4{1, 1, 1, 0.5}, Type: param.ArtRotate, Line: true}
	localPlayer.NewCurve(leftCircle)

	// if p.LeftWeapon != nil {
	// 	// w := float32(p.LeftWeapon.BulletParam.Lifetime.Seconds() * 10)
	// 	cursorLeft := param.Art{
	// 		Name:     "left",
	// 		Value:    0,
	// 		MaxValue: float32(p.LeftWeapon.Delay.Seconds()),
	// 		W:        10, H: 0.5,
	// 		Color:    mgl32.Vec4{0.3, 0.3, 0.9, 0.5},
	// 		LocalPos: mgl32.Vec3{1, 1, 1},
	// 		Type:     param.ArtRotate,
	// 	}
	// 	localPlayer.NewArt(cursorLeft)
	// }

	// if p.RightWeapon != nil {
	// 	// w := float32(p.RightWeapon.BulletParam.Lifetime.Seconds() * 10)
	// 	// log.Println(w)
	// 	cursorRight := param.Art{
	// 		Name:     "right",
	// 		Value:    0,
	// 		MaxValue: float32(p.LeftWeapon.Delay.Seconds()),
	// 		W:        10, H: 0.1,
	// 		Color:    mgl32.Vec4{0.2, 0.2, 1, 0.5},
	// 		LocalPos: mgl32.Vec3{-1, -1, 1},
	// 		Type:     param.ArtRotate,
	// 	}
	// 	localPlayer.NewArt(cursorRight)
	// }

	engine.Movement = playerMovement
	engine.Window.SetMouseButtonCallback(mouseControl)
	// engine.Window.SetKeyCallback(keyboardControl)
}

func playerMovement(dt float32) {
	cpx, cpy := cursor.Position()
	LookAtTarget(localPlayer, cpx, cpy, dt)

	// log.Println("player:", localPlayer.Velocity().Length())

	if localPlayer.Shape.Body.Velocity().Length() < localPlayer.Player.MovSpeed {
		dist := localPlayer.Distance(cursor)
		if dist > localPlayer.Player.MovSpeed {
			dist = localPlayer.Player.MovSpeed
		}

		localPlayer.AddVelocity(localPlayer.VectorForward(dt * localPlayer.Player.MovSpeed * 0.1 * dist))
	}

	if localPlayer.Player.LeftWeapon != nil {
		Fire(localPlayer.Player.LeftWeapon, localPlayer)

		if localPlayer.Player.LeftWeapon.Delay > 0 {
			t := time.Now().Sub(localPlayer.Player.LeftWeapon.DelayTime)
			value := float32(t.Seconds())
			if value < 0 {
				value = 0
			}
			cursorWeaponDelay("left", value)
		}
	}
	if localPlayer.Player.RightWeapon != nil {
		Fire(localPlayer.Player.RightWeapon, localPlayer)

		if localPlayer.Player.RightWeapon.Delay > 0 {

			if timeNil(localPlayer.Player.RightWeapon.DelayTime) {
				cursorWeaponDelay("right", 1)
			} else {
				t := localPlayer.Player.RightWeapon.DelayTime.Sub(time.Now())
				value := t.Seconds() / localPlayer.Player.RightWeapon.Delay.Seconds()
				cursorWeaponDelay("right", float32(value))
			}
		}
	}
}

func cursorWeaponDelay(childname string, value float32) {
	delayBar, ok := localPlayer.ArtRotate[childname]
	if !ok {
		log.Println("child bar not found")
		return
	}
	delayBar.Value = value
	delayBar.Resize()
}

func mouseControl(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {

	if button == 0 {
		localPlayer.Player.LeftWeapon.Shoot = action == 1
	}

	if button == 1 {
		localPlayer.Player.RightWeapon.Shoot = action == 1
	}

	if action == 1 {
		object := engine.Hit(cursor.Position())
		if object != nil {
			log.Println(object.Name)
		}
	}
}

// func keyboardControl(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
// 	log.Println("key", scancode)
// }
