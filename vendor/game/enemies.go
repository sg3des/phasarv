package game

import (
	"engine"
	"game/equip"
	"game/players"
	"game/ships"
	"game/weapons"
	"materials"
	"phys"
	"phys/vect"
	"point"
	"render"

	"github.com/go-gl/mathgl/mgl32"
)

type Enemy struct {
	*players.Player

	Battle *Battle
}

func (b *Battle) CreateEnemy() {
	x, y := engine.GetRandomPoint(50, 50)

	s := &ships.Ship{
		Name:    "enemy",
		Mesh:    "trapeze",
		Texture: "TestCube",
		Type:    ships.Interceptor,
		Size:    mgl32.Vec3{2, 2, 2},
		InitParam: equip.Param{
			Pos:       mgl32.Vec3{x, y},
			Weight:    12,
			Health:    100,
			MovSpeed:  0,
			RotSpeed:  50,
			Energy:    100,
			EnergyAcc: 100,
			Metal:     100,
			MetalAcc:  100,
		},

		LeftWpnPos:  mgl32.Vec3{-1, 0, 0},
		RightWpnPos: mgl32.Vec3{1, 0, 0},

		LeftWeapon: &weapons.Weapon{
			Type:    weapons.Rocket,
			SubType: weapons.TypeHoming,

			Equip: equip.Equip{
				InitParam: equip.Param{
					WeaponParam: equip.WeaponParam{
						Rate:             2e9,
						BulletMovSpeed:   15,
						BulletRotSpeed:   100,
						Range:            6e9,
						Damage:           10,
						Ammunition:       5000,
						ReloadTime:       0, //2sec
						ReloadEnergyCost: 1,
						ReloadMetalCost:  2,
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

	e := &Enemy{
		Player: &players.Player{
			Name: s.Name,
			Ship: s,
		},
		Battle: b,
	}

	e.Ship.Player = e

	e.Ship.CreateCursor(mgl32.Vec4{1, 0.1, 0.1, 0.7})
	// p.CreatePlayer()

	e.Ship.Object.AddCallback(e.Rotation, e.Attack)
	// engine.AddCallback(p.EnemyRotation, p.EnemyAttack)
}

func (e *Enemy) Rotation(dt float32) {
	if e == nil || e.Battle == nil || e.Battle.Finished {
		return
	}

	if e.Target == nil || e.Ship.Object.Distance(e.Target.Ship.Object) > 20 {
		e.Target = e.FindClosestPlayer(20)
	}

	if e.Target == nil {
		return
	}

	e.Ship.Cursor.SetPosition(e.Target.Ship.Object.Position())
	e.Ship.LeftWeapon.UpdateCursor(e.Target.Ship.Object.Position())

	e.TargetAngle = e.Ship.Rotate(dt, e.Target.Ship.Object.PositionVec2())
}

func (e *Enemy) Attack(dt float32) {
	if e == nil || e.Battle == nil || e.Battle.Finished {
		return
	}

	if e.Target == nil {
		e.Ship.LeftWeapon.ToShoot = false
		return
	}

	e.Ship.LeftWeapon.Target = e.Target.Ship.Object

	if e.Ship.LeftWeapon.Type != weapons.Rocket && vect.FAbs(e.TargetAngle) > 0.2 {
		e.Ship.LeftWeapon.ToShoot = false
		return
	}

	e.Ship.LeftWeapon.ToShoot = true
	e.Ship.LeftWeapon.Fire()

	return
}

// func (p *Player) EnemyCursor(dt float32) bool {
// 	if p == nil {
// 		return false
// 	}

// 	players := p.FindClosestPlayer(Players, 20)

// 	return true
// }

func (e *Enemy) FindClosestPlayer(length float32) (closest *players.Player) {
	var mindist float32 = 99

	for _, player := range e.Battle.Players {
		dist := e.Ship.Object.Distance(player.Ship.Object)
		if dist < length && dist < mindist {
			mindist = dist
			closest = player
		}
	}

	return closest
}

func (e *Enemy) Kill() {
	e.Kills++
}

func (e *Enemy) Death() {
	e.Deaths++
}
