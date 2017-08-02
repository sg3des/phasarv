package db

import (
	"game/equip"
	"game/ships"
	"game/weapons"

	"github.com/go-gl/mathgl/mgl32"
)

func GetHangarShip(name string) *ships.Ship {
	return &ships.Ship{
		Name: name,
		Img:  "red",
		Type: ships.Fighter,

		InitParam: ships.Param{
			Size: mgl32.Vec3{1, 1, 1},
			Param: equip.Param{
				Weight:    20,
				Health:    40,
				MovSpeed:  10,
				RotSpeed:  10,
				Energy:    10,
				EnergyAcc: 1,
				Metal:     22,
				MetalAcc:  0.1,
			},
		},

		Slots: []equip.Slot{
			equip.Slot{"41%", "70%", "18%", "18%", equip.TypeEngine},
			equip.Slot{"10%", "35%", "15%", "15%", equip.TypeWeapon},
			equip.Slot{"75%", "35%", "15%", "15%", equip.TypeWeapon},
		},
	}
}

func GetWeapon(name string) (w *weapons.Weapon) {
	switch name {
	case "gun-104":
		w = &weapons.Weapon{
			EquipType: equip.TypeWeapon,
			Name:      name,
			Img:       "gun-00",
			Type:      weapons.Gun,
			InitParam: weapons.Param{
				Damage:         8,
				Rate:           3e8, //300ms
				Range:          30,
				Angle:          0.3,
				Ammunition:     20,
				ReloadTime:     2e9, //2sec
				ReloadCost:     20,
				BulletMovSpeed: 30,
				Param: equip.Param{
					Weight: 3,
				},
			},
		}
	case "rocket-424":
		w = &weapons.Weapon{
			EquipType: equip.TypeWeapon,
			Name:      name,
			Img:       "rocket-00",
			Type:      weapons.Rocket,
			SubType:   weapons.TypeHoming,
			InitParam: weapons.Param{
				Damage:         15,
				Rate:           5e8, //300ms
				Range:          30e9,
				Angle:          0.3,
				Ammunition:     3,
				ReloadTime:     5e9, //2sec
				ReloadCost:     30,
				BulletMovSpeed: 20,
				BulletRotSpeed: 5,
				Param: equip.Param{
					Weight: 3,
				},
			},
		}
	case "laser-89":
		w = &weapons.Weapon{
			EquipType: equip.TypeWeapon,
			Name:      name,
			Img:       "laser-00",
			Type:      weapons.Laser,
			InitParam: weapons.Param{
				Damage:     15,
				Rate:       5e8, //300ms
				Range:      30e9,
				Angle:      0,
				Ammunition: 5,
				ReloadTime: 5e9, //2sec
				ReloadCost: 25,
				Param: equip.Param{
					Weight: 3,
				},
			},
		}
	}

	return
}

func GetEquipment(name string) (e *equip.Equip) {
	switch name {
	case "engine-22":
		e = &equip.Equip{
			Type: equip.TypeEngine,
			Name: name,
			Param: equip.Param{
				Weight:   12,
				MovSpeed: 15,
			},
			Img: "engine-00",
		}
	}

	return
}
