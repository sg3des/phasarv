package db

import (
	"engine"
	"game"
	"game/equip"
	"game/ships"
	"game/weapons"
	"materials"
	"phys"
	"phys/vect"
	"point"
	"render"

	"github.com/go-gl/mathgl/mgl32"
)

func GetPlayer(name string) *game.Player {
	return &game.Player{Name: name, Ship: GetShip(name)}
}

func GetShip(name string) *ships.Ship {
	s := &ships.Ship{
		Object: &engine.Object{
			Name: "player",
			P:    &point.Param{Size: point.P{X: 2, Y: 2}},
			PI:   &phys.Instruction{Mass: 12, Group: phys.GROUP_PLAYER, ShapeType: phys.ShapeType_Box},
			RI: &render.Instruction{
				MeshName:    "trapeze",
				Material:    &materials.Instruction{Name: "player", Texture: "TestCube", Shader: "basic", SpecLevel: 1},
				Shadow:      true,
				Transparent: false,
			},
		},

		InitParam: ships.Param{
			Size: mgl32.Vec3{1, 1, 1},
			Param: equip.Param{
				Health:    100,
				MovSpeed:  0,
				RotSpeed:  0,
				RollAngle: 1.57,
			},
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
		LeftWeapon: &weapons.Weapon{
			Type: weapons.Laser,
			InitParam: weapons.Param{
				Pos:    vect.Vect{-1, 0},
				Rate:   1e8,
				Range:  20e9,
				Damage: 50,
				Angle:  0.8,
			},
			BulletObj: &engine.Object{
				Name: "bullet",
				P:    &point.Param{Size: point.P{1, 1, 1}},
				RI: &render.Instruction{
					MeshName:    "vector",
					Material:    &materials.Instruction{Name: "laser", Texture: "laser", Shader: "blend"},
					Transparent: true,
				},
			},
		},

		RightWeapon: &weapons.Weapon{
			Type:    weapons.Rocket,
			SubType: weapons.TypeAimed,
			InitParam: weapons.Param{
				Pos:            vect.Vect{1, 0},
				Angle:          1.28,
				Rate:           1e9,
				BulletMovSpeed: 10,
				BulletRotSpeed: 50,
				Range:          3e9,
				Damage:         50,
			},

			BulletObj: &engine.Object{
				Name: "bullet",
				P:    &point.Param{Size: point.P{0.1, 0.1, 0.1}},
				PI:   &phys.Instruction{Mass: 0.5},
				RI: &render.Instruction{
					MeshName:    "rocket",
					Material:    &materials.Instruction{Name: "bullet", Texture: "gray", Shader: "color"},
					Transparent: true,
				},
			},
		},
	}

	return s
}
