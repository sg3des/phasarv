package db

import (
	"engine"
	"game"
	"materials"
	"phys"
	"point"
	"render"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

func GetPlayer(name string) *game.Player {
	var v mgl32.Vec3

	v[0] = 1

	player := &game.Player{
		Name: name,
		Object: &engine.Object{
			Name: "player",
			PI:   &phys.Instruction{W: 2, H: 2, Mass: 12, Group: 1, ShapeType: phys.ShapeType_Box},
			RI: &render.Instruction{
				MeshName:    "trapeze",
				Material:    &materials.Instruction{Name: "player", Texture: "TestCube", Shader: "basic", SpecLevel: 1},
				Shadow:      true,
				Transparent: false,
			},
		},

		InitParam: game.PlayerParam{
			Health:    100,
			MovSpeed:  20,
			RotSpeed:  50,
			RollAngle: 1.5,
		},

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
		LeftWeapon: &game.Weapon{
			Bullet: game.Bullet{
				Type:     "laser",
				Lifetime: 2500 * time.Millisecond,
				Damage:   50,
				Object: &engine.Object{
					Name: "bullet",

					P: point.Param{Size: point.P{1, 1, 1}},
					RI: &render.Instruction{
						MeshName:    "plane",
						Material:    &materials.Instruction{Name: "laser", Texture: "laser", Shader: "blend"},
						Transparent: true,
					},
				},
			},
			X: -1,
			// Delay:      500 * time.Millisecond,
			AttackRate: 100 * time.Millisecond,
		},
		RightWeapon: &game.Weapon{
			Bullet: game.Bullet{
				Type:     "rocket",
				SubType:  "aimed",
				MovSpeed: 30,
				RotSpeed: 100,
				Lifetime: 20000 * time.Millisecond,
				Damage:   200,

				Object: &engine.Object{
					Name: "bullet",

					PI: &phys.Instruction{W: 0.1, H: 0.1, Mass: 0.5},
					RI: &render.Instruction{
						MeshName:    "rocket",
						Material:    &materials.Instruction{Name: "bullet", Texture: "gray", Shader: "color"},
						Transparent: true,
					},
				},
			},
			X: 1,
			// Delay:      500 * time.Millisecond,
			AttackRate: 1000 * time.Millisecond,
		},
	}
	return player
}
