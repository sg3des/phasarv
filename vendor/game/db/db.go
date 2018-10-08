package db

import (
	"engine"
	"fmt"
	"game"
	"game/equip"
	"game/ships"
	"game/weapons"
	"log"
	"materials"
	"phys"
	"point"
	"render"
	"time"

	"github.com/1lann/cete"
	"github.com/go-gl/mathgl/mgl32"
)

var db *cete.DB

func init() {
	var err error
	db, err = cete.Open("database")
	if err != nil {
		log.Fatalln("failed open database", err)
	}

	if db.Table("players") == nil {
		SetInitialValues()
	}

	// table := db.Table("ships").All()
	// for table.Next() {
	// 	var item ships.Ship
	// 	err := table.Decode(&item)
	// 	if err != nil {
	// 		log.Println("ERROR:", err)
	// 		break
	// 	}
	// 	log.Printf("SHIP: %+v", item)
	// }

	// table = db.Table("weapons").All()
	// for table.Next() {
	// 	var item weapons.Weapon
	// 	err := table.Decode(&item)
	// 	if err != nil {
	// 		log.Println("ERROR:", err)
	// 		break
	// 	}
	// 	log.Printf("WEAPON: %+v", item)
	// }

	// table = db.Table("equip").All()
	// for table.Next() {
	// 	var item equip.Equip
	// 	err := table.Decode(&item)
	// 	if err != nil {
	// 		log.Println("ERROR:", err)
	// 		break
	// 	}
	// 	log.Printf("EQUIP: %+v", item)
	// }

}

func GetPlayer(name string) *game.Player {
	return &game.Player{
		Name: name,
		Ship: GetShip("ship0"),
		Inventory: []*equip.Equip{
			&equip.Equip{Name: "gun0", EquipType: equip.TypeWeapon},
			&equip.Equip{Name: "rocket0", EquipType: equip.TypeWeapon},
			&equip.Equip{Name: "laser0", EquipType: equip.TypeWeapon},
			&equip.Equip{Name: "engine0", EquipType: equip.TypeEngine},
		},
		// WeaponsIDs: []string{"gun0", "rocket0", "laser0"},
		// EquipIDs:   []string{"engine0"},
	}
}

// func GetPlayerWeapons(IDs []string) (list []*weapons.Weapon) {
// 	table := db.Table("weapons")
// 	for _, id := range IDs {
// 		var w *weapons.Weapon
// 		table.Get(id, &w)
// 		list = append(list, w)
// 	}
// 	return
// }

// func GetPlayerEquip(IDs []string) (list []*equip.Equip) {
// 	table := db.Table("equip")
// 	for _, id := range IDs {
// 		var e *equip.Equip
// 		table.Get(id, &e)
// 		list = append(list, e)
// 	}
// 	return
// }

func GetShip(id string) (s *ships.Ship) {
	_, err := db.Table("ships").Get(id, &s)
	if err != nil {
		log.Printf("ship by id: '%s' not found\n", id)
	}
	return
}

func GetWeapon(id string) (w *weapons.Weapon) {
	_, err := db.Table("weapons").Get(id, &w)
	if err != nil {
		log.Printf("weapon by id: '%s' not found\n", id)
	}
	return
}

func GetEquip(id string) (e *equip.Equip) {
	_, err := db.Table("equip").Get(id, &e)
	if err != nil {
		log.Printf("equip by id: '%s' not found\n", id)
	}
	return
}

func SetInitialValues() {
	db.NewTable("players")
	db.NewTable("ships")
	db.NewTable("weapons")
	db.NewTable("equip")

	db.Table("ships").Set("ship0", &ships.Ship{
		Name: "red fighter",
		Img:  "red",
		Mesh: "trapeze",
		Type: ships.Fighter,
		Size: mgl32.Vec3{2, 2, 2},

		InitParam: equip.Param{
			Weight:    12,
			Health:    40,
			MovSpeed:  5,
			RotSpeed:  50,
			RollAngle: 1.57,
			Energy:    10,
			EnergyAcc: 1,
			Metal:     22,
			MetalAcc:  0.1,
		},

		LeftWpnPos:  mgl32.Vec3{-1, 0, 0},
		RightWpnPos: mgl32.Vec3{1, 0, 0},

		Slots: []equip.Slot{
			equip.Slot{"engine0", "41%", "70%", "18%", "18%", equip.TypeEngine, 0, 0},
			equip.Slot{"wpn-l", "10%", "35%", "15%", "15%", equip.TypeWeapon, equip.Left, 0},
			equip.Slot{"wpn-r", "75%", "35%", "15%", "15%", equip.TypeWeapon, equip.Right, 0},
		},
	})

	db.Table("equip").Set("engine0", &equip.Equip{
		EquipType: equip.TypeEngine,
		Name:      "engine-w12.m15",
		InitParam: equip.Param{
			Weight:   12,
			MovSpeed: 15,
		},
		Img: "engine-00",
	})

	// tableWpns :=
	err := db.Table("weapons").Set("gun0", &weapons.Weapon{
		Equip: equip.Equip{
			EquipType: equip.TypeWeapon,
			Name:      "gun-d8",
			Img:       "gun-00",
			InitParam: equip.Param{
				Weight: 3,
				WeaponParam: equip.WeaponParam{
					Damage:         8,
					Rate:           3e8, //300ms
					Range:          7e8,
					Angle:          0.3,
					Ammunition:     20,
					ReloadTime:     2e9, //2sec
					ReloadCost:     20,
					BulletMovSpeed: 30,
				},
			},
		},
		Type: weapons.Gun,
		BulletObj: &engine.Object{
			Name: "bullet",
			P:    &point.Param{Size: point.P{0.1, 0.1, 0.1}},
			PI:   &phys.Instruction{Mass: 0.5},
			RI: &render.Instruction{
				MeshName:    "bullet",
				Material:    &materials.Instruction{Name: "bullet", Texture: "gray", Shader: "color"},
				Transparent: true,
			},
		},
	})
	if err != nil {
		log.Println(err)
	}

	err = db.Table("weapons").Set("rocket0", &weapons.Weapon{
		Equip: equip.Equip{
			EquipType: equip.TypeWeapon,
			Name:      "rocket-d15",
			Img:       "rocket-00",
			InitParam: equip.Param{
				Weight: 3,
				WeaponParam: equip.WeaponParam{
					Damage:         15,
					Rate:           1e9, //300ms
					Range:          30e9,
					Angle:          3,
					Ammunition:     3,
					ReloadTime:     5e9, //2sec
					ReloadCost:     30,
					BulletMovSpeed: 20,
					BulletRotSpeed: 5,
				},
			},
		},

		Type:    weapons.Rocket,
		SubType: weapons.TypeHoming,

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
	})
	if err != nil {
		log.Println(err)
	}

	err = db.Table("weapons").Set("laser0", &weapons.Weapon{
		Equip: equip.Equip{
			EquipType: equip.TypeWeapon,
			Name:      "laser-d15",
			Img:       "laser-00",
			InitParam: equip.Param{
				Weight: 3,
				WeaponParam: equip.WeaponParam{
					Damage:     15,
					Rate:       1e8,
					Range:      5e9,
					Angle:      0.9,
					Ammunition: 5,
					ReloadTime: 5e9,
					ReloadCost: 25,
				},
			},
		},

		Type: weapons.Laser,
		BulletObj: &engine.Object{
			Name: "bullet",
			P:    &point.Param{Size: point.P{1, 1, 1}},
			RI: &render.Instruction{
				MeshName:    "vector",
				Material:    &materials.Instruction{Name: "laser", Texture: "laser", Shader: "blend"},
				Transparent: true,
			},
		},
	})
	if err != nil {
		log.Println(err)
	}

}

func uuid(s string) string {
	return fmt.Sprintf("%s-%x", s, time.Now().UnixNano()-150202233e10)
}
