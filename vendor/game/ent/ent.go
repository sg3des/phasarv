package ent

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/sg3des/fizzgui"
)

var AssetsShips = "assets/ships/"
var AssetsItems = "assets/items/"
var Ext = ".png"

func pathJoin(dir, img string) string {
	return filepath.Join(dir, img+".png")
}

func loadTexture(dir, img string) *fizzgui.Texture {
	tex, err := fizzgui.NewTextureImg(pathJoin(dir, img))
	if err != nil {
		log.Println(err)
	}
	return tex
}

type Ship struct {
	Model string
	Class string

	Width  float32
	Height float32

	ShipSpec Spec
	Spec     Spec

	Equipment []*Equipment
	Weapons   []*Weapon

	Img   string
	Mesh  string
	Slots []EquipmentSlot
}

func (s *Ship) DelWeapon(w *Weapon) {
	for i, _w := range s.Weapons {
		if _w == w {
			s.Weapons[i] = nil
			s.Weapons = append(s.Weapons[:i], s.Weapons[i+1:]...)
			break
		}
	}
}

func (s *Ship) DelEquipment(e *Equipment) {
	for i, _e := range s.Equipment {
		if _e == e {
			s.Equipment[i] = nil
			s.Equipment = append(s.Equipment[:i], s.Equipment[i+1:]...)
			break
		}
	}
}

func (s *Ship) Info() []string {
	info := []string{
		fmt.Sprintf("%s - %s", s.Class, s.Model),
		fmt.Sprintf("Weight     : %.1f t.", s.Spec.Weight),
		fmt.Sprintf("Durability : %.0f", s.Spec.Durability),
		fmt.Sprintf("Shield     : %.0f", s.Spec.Shield),
		fmt.Sprintf("Speed      : %.0f km/h", s.Spec.Speed),
		fmt.Sprintf("Controll.  : %.0f", s.Spec.Controllability),
		fmt.Sprintf("Energy     : %.0f e [+%.1f e/s] ", s.Spec.Energy, s.Spec.EnergyAcc),
		"",
	}

	for _, w := range s.Weapons {
		winfo := []string{
			fmt.Sprintf("%s %s - %s", w.Type, w.SubType, w.Model),
			fmt.Sprintf("  DPS        : %.1f", w.Damage/float32(w.AttackRate.Seconds())),
			fmt.Sprintf("  Range      : %.0f", w.AttackRange),
			fmt.Sprintf("  Ammunition : %d", w.Ammunition),
			"",
		}
		info = append(info, winfo...)
	}

	return info
}

type Spec struct {
	Weight float32

	Durability float32
	Shield     float32

	Speed           float32
	Controllability float32

	Energy    float32
	EnergyAcc float32

	Metal    float32
	MetalAcc float32
}

func (s Spec) Summ(s2 Spec) Spec {
	s.Weight += s2.Weight
	s.Durability += s2.Durability
	s.Shield += s2.Shield
	s.Speed += s2.Speed
	s.Controllability += s2.Controllability
	s.Energy += s2.Energy
	s.EnergyAcc += s2.EnergyAcc
	return s
}

func (s Ship) LoadImgUI() *fizzgui.Texture {
	return loadTexture(AssetsShips, s.Img)
}

type EquipmentSlot struct {
	X, Y, W, H string
	ET         EquipmentType
}

type EquipmentType string

func (et EquipmentType) Str() string {
	return string(et)
}

const (
	ETweapon    EquipmentType = "weapon"
	ETengine    EquipmentType = "engine"
	ETgenerator EquipmentType = "generator"
	ETshield    EquipmentType = "shield"
	ETradar     EquipmentType = "radar"
)

type Equipment struct {
	ET    EquipmentType
	Model string

	Spec Spec

	Img string
}

func (e Equipment) LoadImgUI() *fizzgui.Texture {
	return loadTexture(AssetsItems, e.Img)
}

func (e Equipment) ImgPath() string {
	return pathJoin(AssetsItems, e.Img)
}

type Weapon struct {
	ET      EquipmentType
	Model   string
	Type    string
	SubType string

	Weight float32

	Damage      float32
	AttackRate  time.Duration
	AttackAngle float32
	AttackRange float32
	AttackDelay time.Duration

	Ammunition int
	ReloadTime time.Duration
	ReloadCost float32

	BulletSpeed float32
	BulletRot   float32

	Img string
}

func (w Weapon) LoadImgUI() *fizzgui.Texture {
	return loadTexture(AssetsItems, w.Img)
}

func (w Weapon) ImgPath() string {
	return pathJoin(AssetsItems, w.Img)
}
