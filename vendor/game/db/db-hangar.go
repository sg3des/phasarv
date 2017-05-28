package db

import "game/ent"

var Ship = &ent.Ship{
	Model:  "",
	Class:  "fighter",
	Width:  1,
	Height: 1,
	ShipSpec: ent.Spec{
		Weight:          20,
		Durability:      40,
		Controllability: 10,
		Energy:          10,
		EnergyAcc:       1,
		Metal:           22,
		MetalAcc:        0.1,
	},
	Img:  "red",
	Mesh: "trapeze",
	Slots: []ent.EquipmentSlot{
		ent.EquipmentSlot{"41%", "70%", "18%", "18%", ent.ETengine},
		ent.EquipmentSlot{"10%", "35%", "15%", "15%", ent.ETweapon},
		ent.EquipmentSlot{"75%", "35%", "15%", "15%", ent.ETweapon},
	},
}

func GetShip(model string) *ent.Ship {
	Ship.Model = model
	return Ship
}

func GetWeapon(model string) (w *ent.Weapon) {
	switch model {
	case "gun-104":
		w = &ent.Weapon{
			ET:          ent.ETweapon,
			Model:       model,
			Type:        "gun",
			Weight:      3,
			Damage:      8,
			AttackRate:  3e8, //300ms
			AttackAngle: 0.3,
			AttackRange: 30,
			Ammunition:  20,
			ReloadTime:  2e9, //2sec
			ReloadCost:  20,
			BulletSpeed: 30,
			Img:         "gun-00",
		}
	case "rocket-424":
		w = &ent.Weapon{
			ET:          ent.ETweapon,
			Model:       model,
			Type:        "rocket",
			SubType:     "homing",
			Weight:      5,
			Damage:      15,
			AttackRate:  5e8, //500ms
			AttackAngle: 1.54,
			AttackRange: 100,
			Ammunition:  3,
			ReloadTime:  5e9, //5sec
			ReloadCost:  30,
			BulletSpeed: 20,
			BulletRot:   5,
			Img:         "rocket-00",
		}
	case "laser-89":
		w = &ent.Weapon{
			ET:          ent.ETweapon,
			Model:       model,
			Type:        "laser",
			Weight:      4,
			Damage:      15,
			AttackRate:  2e9, //2sec
			AttackRange: 40,
			AttackDelay: 5e8, //500ms
			Ammunition:  5,
			ReloadTime:  5e9, //5sec
			ReloadCost:  25,
			Img:         "laser-00",
		}
	}

	return
}

func GetEquipment(model string) (e *ent.Equipment) {
	switch model {
	case "engine-22":
		e = &ent.Equipment{
			ET:    ent.ETengine,
			Model: model,
			Spec: ent.Spec{
				Weight: 12,
				Speed:  15,
			},
			Img: "engine-00",
		}
	}

	return
}
