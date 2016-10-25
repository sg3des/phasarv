package param

import (
	"phys"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	Name   string
	Object Object

	Health, MovSpeed, RotSpeed, RollAngle float32

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
	LocalPos mgl32.Vec3
	Type     int //ArtStatic, ArtRotate
	Line     bool

	Material
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
	Mesh
	Material
	Pos
	*Phys

	MaxRollAngle float32

	Transparent bool
}

type Mesh struct {
	Model, Texture, Shader string
	Shadow                 bool
}

type Material struct {
	Name, Shader, Texture string
	DiffColor             mgl32.Vec4
	SpecLevel             float32
}

type Pos struct {
	X, Y, Z float32
}

type Phys struct {
	W, H, Mass float32
	Type       phys.ShapeType
	Group      int
}

type Light struct {
	Name      string
	Shadow    bool
	Intensity float32
	Pos       Pos
}
