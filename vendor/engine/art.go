package engine

import (
	"point"
	"render"

	"github.com/go-gl/mathgl/mgl32"
)

type Art struct {
	Name string

	Value    float32
	MaxValue float32

	Art *render.Art

	P  *point.Param
	RI *render.Instruction
}

func (o *Object) AppendArt(a *Art) {
	if a == nil {
		return
	}
	a.Art = a.RI.CreateArt(a.P)
	o.renderable.AppendArt(a.Art)

	o.Arts = append(o.Arts, a)
}

func (o *Object) GetArt(name string) *Art {
	for _, a := range o.Arts {
		if a.Name == name {
			return a
		}
	}
	return nil
}

//Resize bar
func (a *Art) Resize(value float32) {
	a.Value = value
	if a.Art.Body != nil {
		if a.Art.Line {
			a.Art.Body.FaceCount = uint32(a.MaxValue * a.Value)
		} else {
			percent := a.Value / a.MaxValue
			a.Art.Body.Scale = mgl32.Vec3{percent, 1, 1}
		}
	}
}
