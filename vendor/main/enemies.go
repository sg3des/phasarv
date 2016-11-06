package main

import (
	"engine"
	"log"
	"param"
	"phys"
	"phys/vect"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

func createEnemy(x, y float32) {
	log.Println("createEnemy", x, y)
	var p = &Player{}
	p.Param = param.Player{
		Name: "enemy0",
		Object: param.Object{
			Name:     "enemy",
			Mesh:     param.Mesh{Model: "trapeze"},
			Material: param.Material{Name: "enemy", Texture: "TestCube", Shader: "basic", SpecLevel: 1},
			Phys:     &param.Phys{W: 2, H: 2, Mass: 12, Group: 1, Type: phys.ShapeType_Box},
			Shadow:   true,
		},
		Health:   100,
		MovSpeed: 0,
		RotSpeed: 50,

		LeftWeapon: &param.Weapon{
			BulletParam: param.Bullet{
				Type:     "rocket",
				SubType:  "homing",
				MovSpeed: 20,
				RotSpeed: 50,
				Lifetime: 10000 * time.Millisecond,
				Damage:   20,
			},
			BulletObject: param.Object{
				Name:     "bullet",
				Mesh:     param.Mesh{Model: "rocket"},
				Material: param.Material{Name: "color", Texture: "gray", Shader: "color"},
				Phys:     &param.Phys{W: 0.1, H: 0.1, Mass: 1},
			},
			AttackRate: 1000 * time.Millisecond,
		},
	}

	p.CreatePlayer()
	p.Object.SetPosition(x, y)
	p.respawnPoint = mgl32.Vec2{x, y}
	// p.Object.Callback = p.EnemyAttack

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

	if p.Target == nil {
		p.Param.LeftWeapon.Shoot = false
		return true
	}

	p.Param.LeftWeapon.BulletParam.TargetObject = p.Target.Object

	// log.Println(p.Param.LeftWeapon.BulletParam.Type, vect.FAbs(p.targetAngle))
	if p.Param.LeftWeapon.BulletParam.Type != "rocket" && vect.FAbs(p.targetAngle) > 0.2 {
		p.Param.LeftWeapon.Shoot = false
		return true
	}

	// if p.Param.LeftWeapon.BulletParam.Type == "rocket" {
	// 	 = p.Target.Object
	// }

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
