package db

import (
	"param"
	"phys"
	"time"
)

func GetPlayer(name string) param.Player {
	player := param.Player{
		Name: "player0",
		Object: param.Object{
			Name:         "player",
			Mesh:         param.Mesh{Model: "ship", Shadow: true},
			Material:     param.Material{Name: "player", Texture: "TestCube", Shader: "diffuse_texbumped_shadows", SpecLevel: 1},
			Phys:         &param.Phys{W: 2, H: 2, Mass: 12, Group: 1, Type: phys.ShapeType_Box},
			MaxRollAngle: 1.5,
		},
		Health:   100,
		MovSpeed: 20,
		RotSpeed: 50,

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
				Material:    param.Material{Name: "laser", Texture: "laser", Shader: "blend"},
				Phys:        &param.Phys{W: 0.5, Mass: 0.5},
				Transparent: true,
			},

			X:          -1,
			Delay:      500 * time.Millisecond,
			AttackRate: 500 * time.Millisecond,
		},
		RightWeapon: &param.Weapon{
			BulletParam: param.Bullet{
				Type:     "rocket",
				SubType:  "aimed",
				MovSpeed: 20,
				RotSpeed: 60,
				Lifetime: 5000 * time.Millisecond,
				Damage:   30,
			},
			BulletObject: param.Object{
				Name:     "bullet",
				Mesh:     param.Mesh{Model: "rocket"},
				Material: param.Material{Name: "bullet", Texture: "gray", Shader: "color"},
				Phys:     &param.Phys{W: 0.1, H: 0.1, Mass: 0.5},
			},
			X: 1,
			// Delay:      500 * time.Millisecond,
			AttackRate: 1000 * time.Millisecond,
		},
	}
	return player
}
