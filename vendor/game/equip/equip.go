package equip

import (
	"log"
	"path/filepath"
	"time"

	"github.com/go-gl/mathgl/mgl32"
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

//
//
//Param is main structure contains common parameters
type Param struct {
	Pos                           mgl32.Vec3
	Weight                        float32
	MovSpeed, RotSpeed, RollAngle float32
	Health, Shield                float32
	Energy, EnergyAcc             float32
	Metal, MetalAcc               float32

	WeaponParam
}

//WeaponParam is parameters of weapons
type WeaponParam struct {
	Delay time.Duration
	Rate  time.Duration
	Range time.Duration

	Angle float32

	Ammunition int
	ReloadTime time.Duration
	ReloadCost float32

	Damage float32

	BulletMovSpeed float32
	BulletRotSpeed float32
}

func (p Param) Summ(p2 Param) Param {
	p.Weight += p2.Weight
	p.Health += p2.Health
	p.Shield += p2.Shield
	p.MovSpeed += p2.MovSpeed
	p.RotSpeed += p2.RotSpeed
	p.RollAngle += p2.RollAngle
	p.Energy += p2.Energy
	p.EnergyAcc += p2.EnergyAcc
	return p
}

//
//
//Equip is equipment for ships, sush as engine, generators, shields etc...
type Equip struct {
	Name      string
	SlotName  string
	Img       string
	EquipType Type

	InitParam Param
	CurrParam Param
}

func (e Equip) LoadImgUI() *fizzgui.Texture {
	return loadTexture(AssetsItems, e.Img)
}

func (e Equip) ImgPath() string {
	return pathJoin(AssetsItems, e.Img)
}

//
//
//Player

//
//
//Slot equipment slot on ships
type Slot struct {
	Name       string //should be unique
	X, Y, W, H string
	Type       Type
	Side       Side
	Size       Size
}

//
//
//Type of equipment
type Type byte

func (t Type) Str() string {
	switch t {
	case TypeWeapon:
		return "weapon"
	case TypeEngine:
		return "engine"
	case TypeGenerator:
		return "generator"
	case TypeShield:
		return "shield"
	case TypeRadar:
		return "radar"
	}

	log.Println("WARNING! unknown type: '%s'", t)
	return "unknown"
}

const (
	TypeWeapon    Type = 'w'
	TypeEngine    Type = 'e'
	TypeGenerator Type = 'g'
	TypeShield    Type = 's'
	TypeRadar     Type = 'r'
)

type Side byte

const (
	Front Side = 'f'
	Back  Side = 'b'
	Left  Side = 'l'
	Right Side = 'r'
)

func (s Side) String() string {
	switch s {
	case Front:
		return "front"
	case Back:
		return "back"
	case Left:
		return "left"
	case Right:
		return "right"
	}

	return ""
}

type Size uint8

const (
	SizeAny    Size = 0
	SizeSmall  Size = 1
	SizeMedium Size = 2
	SizeBig    Size = 3
)
