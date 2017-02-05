package game

import (
	"engine"
	"log"
	"materials"
	"phys"
	"phys/vect"
	"point"
	"render"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

func CreateEnemy() {
	x, y := GetRandomPoint(50, 50)
	log.Println("createEnemy", x, y)
	var p = &Player{
		Name: "enemy0",
		Object: &engine.Object{
			Name: "enemy",
			P: point.Param{
				Pos: point.P{x, y, 0},
			},
			RI: &render.Instruction{
				MeshName: "trapeze",
				Material: &materials.Instruction{Name: "enemy", Texture: "TestCube", Shader: "basic", SpecLevel: 1},

				Shadow: true,
			},
			PI: &phys.Instruction{W: 1, H: 1, Mass: 12, Group: 1, ShapeType: phys.ShapeType_Box},
		},

		InitParam: PlayerParam{
			Health:   100,
			MovSpeed: 0,
			RotSpeed: 50,
		},

		LeftWeapon: &Weapon{
			Type:    Weapons.Rocket,
			SubType: Weapons.RocketType.Homing,

			Bullet: Bullet{
				MovSpeed: 20,
				RotSpeed: 50,
				Lifetime: 10000 * time.Millisecond,
				Damage:   20,
				Object: &engine.Object{
					Name: "bullet",
					RI: &render.Instruction{
						MeshName: "rocket",
						Material: &materials.Instruction{Name: "color", Texture: "gray", Shader: "color"},
					},
					PI: &phys.Instruction{W: 0.1, H: 0.1, Mass: 1},
				},
			},
			AttackRate: 1000 * time.Millisecond,
		},
		respawnPoint: mgl32.Vec2{x, y},
	}

	p.CreateCursor(mgl32.Vec4{1, 0.1, 0.1, 0.7})
	p.CreatePlayer()

	p.Object.AddCallback(p.EnemyRotation, p.EnemyAttack)
	// engine.AddCallback(p.EnemyRotation, p.EnemyAttack)
}

func (p *Player) EnemyRotation(dt float32) {
	if p == nil {
		return
	}

	if p.Target == nil || p.Object.Distance(p.Target.Object) > 20 {
		p.Target = p.FindClosestPlayer(Players, 20)

	}

	if p.Target == nil {
		return
	}

	p.Cursor.SetPosition(p.Target.Object.Position())

	p.targetAngle = p.rotation(dt, p.Target.Object.PositionVec2())
}

func (p *Player) EnemyAttack(dt float32) {
	if p == nil {
		return
	}

	if p.Target == nil {
		p.LeftWeapon.ToShoot = false
		return
	}

	p.LeftWeapon.TargetPlayer = p.Target

	if p.LeftWeapon.Type != Weapons.Rocket && vect.FAbs(p.targetAngle) > 0.2 {
		p.LeftWeapon.ToShoot = false
		return
	}

	p.LeftWeapon.ToShoot = true
	p.LeftWeapon.Fire()

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
		dist := p.Object.Distance(player.Object)
		if dist < length && dist < mindist {
			mindist = dist
			closest = player
		}
	}

	return closest
}
