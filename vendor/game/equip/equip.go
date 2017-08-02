package equip

import (
	"log"
	"path/filepath"

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

//Param is main structure contains common parameters
type Param struct {
	Weight                        float32
	MovSpeed, RotSpeed, RollAngle float32
	Health, Shield                float32
	Energy, EnergyAcc             float32
	Metal, MetalAcc               float32
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

//Equip is equipment for ships, sush as engine, generators, shields etc...
type Equip struct {
	Name string
	Img  string
	Type Type

	Param Param
}

func (e Equip) LoadImgUI() *fizzgui.Texture {
	return loadTexture(AssetsItems, e.Img)
}

func (e Equip) ImgPath() string {
	return pathJoin(AssetsItems, e.Img)
}

//Slot equipment slot on ships
type Slot struct {
	X, Y, W, H string
	ET         Type
}

//Type of equipment
type Type string

func (et Type) Str() string {
	return string(et)
}

const (
	TypeWeapon    Type = "weapon"
	TypeEngine    Type = "engine"
	TypeGenerator Type = "generator"
	TypeShield    Type = "shield"
	TypeRadar     Type = "radar"
)
