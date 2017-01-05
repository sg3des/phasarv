package render

import (
	"assets"
	"log"
	"materials"
	"math"
	"point"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
)

type Renderable struct {
	Body  *fizzle.Renderable
	Shape *fizzle.Renderable

	Shadow      bool
	Transparent bool

	ArtStatic []*Art
	ArtRotate []*Art

	arts []*Art

	P  point.Param
	RI *Instruction

	needDestroy bool
}

func (r *Renderable) AppendArt(a *Art) {
	if a.Static {
		r.ArtStatic = append(r.ArtStatic, a)
	} else {
		r.ArtRotate = append(r.ArtRotate, a)
	}
}

type Instruction struct {
	MeshName string
	Material *materials.Instruction

	Line bool //only for art

	Shadow      bool
	Transparent bool
}

func (i *Instruction) Create(p point.Param) *Renderable {
	r := &Renderable{
		// Body: &fizzle.Renderable{},
		// Body:        i.createBody(p),
		Shadow:      i.Shadow,
		Transparent: i.Transparent,
	}

	if p.Static {
		Scene = append(Scene, r)
	} else {
		Renderables[r] = true
	}

	r.P = p
	r.RI = i

	return r
}

// func (r *Renderable) createBody() {
// 	switch r.RI.MeshName {
// 	case "plane":
// 		r.Body = fizzle.CreatePlaneV(mgl32.Vec3{0, -r.P.Size.X / 2, 0}, mgl32.Vec3{r.P.Size.Y, r.P.Size.X / 2, 0})
// 	case "box":
// 		log.Println("warning: fixed size")
// 		r.Body = fizzle.CreateCube(-2, -2, -2, 2, 2, 2)
// 	default:
// 		r.Body = assets.GetModel(r.RI.MeshName)
// 	}

// 	r.Body.Material = r.RI.Material.Create()
// 	r.Body.Location = r.P.Pos.Vec3()
// 	r.Body.LocalRotation = mgl32.AnglesToQuat(0, 0, r.P.Angle, 1)
// }

func (i *Instruction) createBody(p point.Param) (body *fizzle.Renderable) {
	switch i.MeshName {
	case "plane":
		body = fizzle.CreatePlaneV(mgl32.Vec3{0, -p.Size.X / 2, 0}, mgl32.Vec3{p.Size.Y, p.Size.X / 2, 0})
	case "box":
		log.Println("warning: fixed size")
		body = fizzle.CreateCube(-2, -2, -2, 2, 2, 2)
	default:
		body = assets.GetModel(i.MeshName)
	}

	body.Material = i.Material.Create()
	body.Location = p.Pos.Vec3()
	body.LocalRotation = mgl32.AnglesToQuat(0, 0, p.Angle, 1)

	return body
}

type Art struct {
	Body   *fizzle.Renderable
	Line   bool
	Static bool

	Pos   mgl32.Vec3
	Angle float32

	P  point.Param
	RI *Instruction
}

func (i *Instruction) CreateArt(p point.Param) *Art {
	a := &Art{
		// Body:   i.createBody(p),
		Line:   i.Line,
		Static: p.Static,
		Pos:    p.Pos.Vec3(),
		Angle:  p.Angle,
		P:      p,
		RI:     i,
	}
	return a
}

func (r *Renderable) Render() {
	if r.Body == nil {
		log.Println("create ", r.RI.MeshName)
		r.Body = r.RI.createBody(r.P)
	}

	if r.needDestroy {
		log.Println("destroy", r.RI.MeshName)
		r.destroy()
		return
	}

	if r.Body == nil {
		log.Println("WTF?!?!", r.RI.MeshName, r.P.Pos)
		return
	}

	// log.Println("render", r.RI.MeshName, fmt.Sprintf("%#v", r.Body.Core))
	render.DrawRenderable(r.Body, nil, perspective, view, camera)

	if r.Shape != nil {
		r.Shape.Location = mgl32.Vec3{0, 0, 1}.Add(r.Body.Location)
		render.DrawLines(r.Shape, r.Shape.Material.Shader, nil, perspective, view, camera)
	}

	r.RenderArts()
}

func (r *Renderable) RenderArts() {
	for _, a := range r.ArtStatic {
		if a.Body == nil {
			a.Body = a.RI.createBody(a.P)
		}

		a.Body.Location = mgl32.Vec3{}.Add(r.Body.Location).Add(a.Pos)
		if a.Line {
			render.DrawLines(a.Body, a.Body.Material.Shader, nil, perspective, view, camera)
		} else {
			render.DrawRenderable(a.Body, nil, perspective, view, camera)
		}
	}

	for _, a := range r.ArtRotate {
		a.Body.Location = mgl32.Vec3{}.Add(r.Body.Location).Add(a.Pos)
		a.Body.LocalRotation = mgl32.AnglesToQuat(0, 0, r.Angle(), 1)

		if a.Line {
			render.DrawLines(a.Body, a.Body.Material.Shader, nil, perspective, view, camera)
		} else {
			render.DrawRenderable(a.Body, nil, perspective, view, camera)
		}
	}
}

func (r *Renderable) Angle() float32 {
	return float32(2 * math.Acos(float64(r.Body.Rotation.W)))
}

func (r *Renderable) Destroy() {
	r.needDestroy = true
	// delete(Renderables, r)
}

func (r *Renderable) destroy() {
	// r.needDestroy = true
	// r.Body.Destroy()
	delete(Renderables, r)
	// r = nil
}
