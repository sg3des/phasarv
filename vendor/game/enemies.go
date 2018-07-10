package game

import (
	"engine"
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

func CreateEnemy() {
	x, y := engine.GetRandomPoint(50, 50)

	s := &ships.Ship{
		// Object: &engine.Object{
		// 	Name: "enemy",
		// 	P: &point.Param{
		// 		Pos: point.P{x, y, 0},
		// 	},
		// 	RI: &render.Instruction{
		// 		MeshName: "trapeze",
		// 		Material: &materials.Instruction{Name: "enemy", Texture: "TestCube", Shader: "basic", SpecLevel: 1},

		// 		Shadow: true,
		// 	},
		// 	PI: &phys.Instruction{Mass: 12, Group: 1, ShapeType: phys.ShapeType_Box},
		// },

		Mesh: "trapeze",
		Size: mgl32.Vec3{2, 2, 2},
		InitParam: equip.Param{
			Pos:      mgl32.Vec3{x, y},
			Weight:   12,
			Health:   100,
			MovSpeed: 0,
			RotSpeed: 50,
		},

		LeftWeapon: &weapons.Weapon{
			Type:    weapons.Rocket,
			SubType: weapons.TypeHoming,

			Equip: equip.Equip{
				InitParam: equip.Param{
					WeaponParam: equip.WeaponParam{
						Rate:           1e9,
						BulletMovSpeed: 20,
						BulletRotSpeed: 50,
						Range:          10e9,
						Damage:         20,
					},
				},
			},

			BulletObj: &engine.Object{
				Name: "bullet",
				RI: &render.Instruction{
					MeshName: "rocket",
					Material: &materials.Instruction{Name: "color", Texture: "gray", Shader: "color"},
				},
				P:  &point.Param{Size: point.P{0.1, 0.1, 0.1}},
				PI: &phys.Instruction{Mass: 1},
			},
		},
	}

	s.Create()

	p := &Player{
		Name: "enemy",
		Ship: s,
	}

	p.Ship.CreateCursor(mgl32.Vec4{1, 0.1, 0.1, 0.7})
	// p.CreatePlayer()

	p.Ship.Object.AddCallback(p.EnemyRotation, p.EnemyAttack)
	// engine.AddCallback(p.EnemyRotation, p.EnemyAttack)
}

func (p *Player) EnemyRotation(dt float32) {
	if p == nil {
		return
	}

	if p.Target == nil || p.Ship.Object.Distance(p.Target.Ship.Object) > 20 {
		p.Target = p.FindClosestPlayer(Players, 20)
	}

	if p.Target == nil {
		return
	}

	p.Ship.Cursor.SetPosition(p.Target.Ship.Object.Position())

	p.targetAngle = p.Ship.Rotate(dt, p.Target.Ship.Object.PositionVec2())
}

func (p *Player) EnemyAttack(dt float32) {
	if p == nil {
		return
	}

	if p.Target == nil {
		p.Ship.LeftWeapon.ToShoot = false
		return
	}

	p.Ship.LeftWeapon.Target = p.Target.Ship.Object

	if p.Ship.LeftWeapon.Type != weapons.Rocket && vect.FAbs(p.targetAngle) > 0.2 {
		p.Ship.LeftWeapon.ToShoot = false
		return
	}

	p.Ship.LeftWeapon.ToShoot = true
	p.Ship.LeftWeapon.Fire()

	return
}

// func (p *Player) EnemyCursor(dt float32) bool {
// 	if p == nil {
// 		return false
// 	}

// 	players := p.FindClosestPlayer(Players, 20)

// 	return true
// }

func (p *Player) FindClosestPlayer(players []*Player, length float32) *Player {
	var mindist float32 = 99
	var closest *Player

	for _, player := range players {
		dist := p.Ship.Object.Distance(player.Ship.Object)
		if dist < length && dist < mindist {
			mindist = dist
			closest = player
		}
	}

	return closest
}
