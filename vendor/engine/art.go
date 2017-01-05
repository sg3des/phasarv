package engine

import (
	"point"
	"render"

	"github.com/go-gl/mathgl/mgl32"
)

type Art struct {
	Name string

	Art *render.Art

	Value    float32
	MaxValue float32

	P  point.Param
	RI *render.Instruction
	// RenderLine bool

	// Node *fizzle.Renderable
	// Material *fizzle.Material

	// Param ObjectParam
}

//AddArt to object
// func (o *Object) AddArt(a *Art) {

// 	// if o.Body.ArtStatic == nil {
// 	// 	o.Body.ArtStatic = make(map[string]*render.Art)
// 	// 	o.Body.ArtRotate = make(map[string]*render.Art)
// 	// }

// 	// if a.Param.StaticType {
// 	// 	o.Body.ArtStatic[a.Name] = a.Art
// 	// } else {
// 	// 	o.Body.ArtRotate[a.Name] = a.Art
// 	// }
// }

// func (a *Art) Create() {
// 	a.Art = a.ri.CreateArt(a.p)

// 	// log.Println("create art node", a.Name, a.ri.MeshName)
// 	// switch a.ri.Node {
// 	// case "plane":
// 	// 	a.Art.Body = fizzle.CreatePlaneV(mgl32.Vec3{-p.Size.X / 2, 0, 0}, mgl32.Vec3{p.Size.X / 2, p.Size.Y, 0})
// 	// case "box":
// 	// 	log.Println("warning: fixed size")
// 	// 	a.Art.Body = fizzle.CreateCube(-2, -2, -2, 2, 2, 2)
// 	// default:
// 	// 	log.Printf("%++v\n\n", a)
// 	// 	a.Art.Body = assets.GetModel(p.Node)
// 	// }

// 	// a.Art.Body.Location = p.Pos.Vec3() //mgl32.Vec3{p.Pos.X(), p.Pos.Y, p.Pos.Z}
// 	// a.Art.Body.Material = NewMaterial(p.Material)
// }

//Resize bar
func (a *Art) Resize() {
	if a.Art.Line {
		a.Art.Body.FaceCount = uint32(a.MaxValue * a.Value)
	} else {
		percent := a.Value / a.MaxValue
		a.Art.Body.Scale = mgl32.Vec3{percent, 1, 1}
	}
}

func (o *Object) GetArt(name string) (*Art, bool) {
	if art, ok := o.Arts[name]; ok {
		return art, true
	}

	// if art, ok := o.ArtStatic[name]; ok {
	// 	return art, true
	// }

	// if art, ok := o.ArtRotate[name]; ok {
	// 	return art, true
	// }

	return nil, false
}

// --------
