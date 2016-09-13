package param

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	Name   string
	Object Object

	Health, MovSpeed, RotSpeed, SubAngle float32

	LeftWeapon, RightWeapon *Weapon
}

const (
	ArtStatic = 0
	ArtRotate = 1
)

type Art struct {
	Name     string
	Value    float32
	MaxValue float32
	W, H     float32
	Seg      int
	Color    mgl32.Vec4
	LocalPos mgl32.Vec3
	Type     int //ArtStatic, ArtRotate
	Line     bool

	Shader  string
	Texture string
}

type Bullet struct {
	Type    string
	SubType string

	TimePoint    time.Time
	Target       mgl32.Vec2
	TargetObject interface{} //crutch - this is must be engine.Object

	RotSpeed float32
	MovSpeed float32
	Lifetime time.Duration

	Damage float32
}

type Weapon struct {
	NextShot  time.Time
	Shoot     bool
	Delay     time.Duration
	DelayTime time.Time

	BulletParam  Bullet
	BulletObject Object

	X float32

	AttackRate time.Duration
}

type Object struct {
	Name string
	Mesh Mesh
	Pos  Pos
	PH   Phys

	Transparent bool
}

type Mesh struct {
	Model, Texture, Shader string
	Shadow                 bool
}

type Pos struct {
	X, Y, Z float32
}

type Phys struct {
	W, H, Mass float32
	Group      int
}

type Light struct {
	Name      string
	Shadow    bool
	Intensity float32
	Pos       Pos
}
