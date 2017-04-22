package db

import (
	"engine"
	"game"
	"materials"
	"phys"
	"phys/vect"
	"point"
	"render"
	"time"
)

func GetPlayer(name string) *game.Player {
	// var v mgl32.Vec3

	// v[0] = 1

	player := &game.Player{
		Name: name,
		Object: &engine.Object{
			Name: "player",
			PI:   &phys.Instruction{W: 1, H: 1, Mass: 12, Group: phys.GROUP_PLAYER, ShapeType: phys.ShapeType_Box},
			RI: &render.Instruction{
				MeshName:    "trapeze",
				Material:    &materials.Instruction{Name: "player", Texture: "TestCube", Shader: "basic", SpecLevel: 1},
				Shadow:      true,
				Transparent: false,
			},
		},

		InitParam: game.PlayerParam{
			Health:    100,
			MovSpeed:  6,
			RotSpeed:  30,
			RollAngle: 1.5,
		},

		// LeftWeapon: &game.Weapon{
		// 	Type: game.Weapons.Gun,
		// 	Bullet: game.Bullet{
		// 		MovSpeed: 30,
		// 		Lifetime: 1000 * time.Millisecond,
		// 		Damage:   20,
		// 		Object: &engine.Object{
		// 			Name: "bullet",
		// 			PI:   &phys.Instruction{W: 0.1, H: 0.1, Mass: 0.5},
		// 			RI: &render.Instruction{
		// 				MeshName: "bullet",
		// 				Material: &materials.Instruction{Name: "bullet", Texture: "gray", Shader: "color"},
		// 			},
		// 		},
		// 	},
		// 	Pos:        vect.Vect{0, 1},
		// 	Angle:      0.3,
		// 	AttackRate: 200 * time.Millisecond,
		// },
		LeftWeapon: &game.Weapon{
			Type: game.Weapons.Laser,
			Bullet: game.Bullet{
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
			Pos: vect.Vect{0, -1},
			// Delay:      500 * time.Millisecond,
			AttackRate: 100 * time.Millisecond,
		},

		RightWeapon: &game.Weapon{
			Type:    game.Weapons.Rocket,
			SubType: game.Weapons.RocketType.Aimed,
			Bullet: game.Bullet{
				MovSpeed: 25,
				RotSpeed: 100,
				Lifetime: 20000 * time.Millisecond,
				Damage:   50,

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
			Pos:   vect.Vect{0, -1},
			Angle: 6.28,
			// Delay:      500 * time.Millisecond,
			AttackRate: 1000 * time.Millisecond,
		},
	}
	return player
}
