package point

import (
	"phys/vect"

	"github.com/go-gl/mathgl/mgl32"
)

type P struct {
	X, Y, Z float32
}

func (p P) Vec3() mgl32.Vec3 {
	return mgl32.Vec3{p.X, p.Y, p.Z}
}

func (p P) Vec2() mgl32.Vec2 {
	return mgl32.Vec2{p.X, p.Y}
}

func (p P) Vect() vect.Vect {
	return vect.Vect{p.X, p.Y}
}

type Param struct {
	Pos    P
	Size   P
	Angle  float32
	Static bool
}
