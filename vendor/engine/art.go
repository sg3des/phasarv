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

	P  point.Param
	RI *render.Instruction
}

func (o *Object) AppendArt(a *Art) {
	if o.Arts == nil {
		o.Arts = make(map[string]*Art)
	}

	a.Art = a.RI.CreateArt(a.P)
	o.renderable.AppendArt(a.Art)

	o.Arts[a.Name] = a
}

func (o *Object) GetArt(name string) (*Art, bool) {
	if art, ok := o.Arts[name]; ok {
		return art, true
	}

	return nil, false
}

//Resize bar
func (a *Art) Resize(value float32) {
	a.Value = value
	if a.Art.Line {
		a.Art.Body.FaceCount = uint32(a.MaxValue * a.Value)
	} else {
		percent := a.Value / a.MaxValue
		a.Art.Body.Scale = mgl32.Vec3{percent, 1, 1}
	}
}
