package main

import (
	"engine"
	"param"
	"phys"
	"phys/vect"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

func initEnemies() {

	var p = &Player{}
	p.Param = param.Player{
		Name: "enemy0",
		Object: param.Object{
			Name:     "enemy",
			Mesh:     param.Mesh{Model: "trapeze", Shadow: true},
			Material: param.Material{Name: "enemy", Texture: "TestCube", Shader: "diffuse_texbumped_shadows", SpecLevel: 1},
			Phys:     &param.Phys{W: 2, H: 2, Mass: 12, Group: 1, Type: phys.ShapeType_Box},
		},
		Health:   100,
		MovSpeed: 0,
		RotSpeed: 30,

		LeftWeapon: &param.Weapon{
			BulletParam: param.Bullet{
				Type:     "gun",
				MovSpeed: 30,
				Lifetime: 10000 * time.Millisecond,
				Damage:   20,
			},
			BulletObject: param.Object{
				Name: "bullet",
				Mesh: param.Mesh{Model: "bullet", Texture: "gray", Shader: "color"},
				Phys: &param.Phys{W: 0.1, H: 0.1, Mass: 1},
			},
			AttackRate: 500 * time.Millisecond,
		},
	}

	p.CreatePlayer()
	p.Object.SetPosition(10, 15)
	p.respawnPoint = mgl32.Vec2{10, 15}

	engine.AddCallback(p.EnemyRotation, p.EnemyAttack)
}

func (p *Player) EnemyRotation(dt float32) bool {
	if p == nil {
		return false
	}

	if p.Target == nil || p.Object.Distance(p.Target.Object) > 20 {
		p.Target = p.FindClosestPlayer(players, 20)
	}

	if p.Target == nil {
		return true
	}

	p.targetAngle = p.rotation(dt, p.Target.Object.PositionVec2())

	return true
}

func (p *Player) EnemyAttack(dt float32) bool {
	if p == nil {
		return false
	}

	if p.Target == nil || vect.FAbs(p.targetAngle) > 0.2 {
		p.Param.LeftWeapon.Shoot = false
		return true
	}

	// log.Println("attack:", vect.FAbs(p.targetAngle), p.Target.Object.Name)
	p.Param.LeftWeapon.Shoot = true
	p.Fire(p.Param.LeftWeapon)

	return true
}

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
