package render

import (
	"assets"
	"log"
	"materials"
	"math"
	"phys"
	"phys/vect"
	"point"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
)

type Renderable struct {
	Body  *fizzle.Renderable
	Shape *fizzle.Renderable

	Shadow      bool
	Transparent bool
	needDestroy bool

	ArtStatic []*Art
	ArtRotate []*Art

	arts []*Art

	P  point.Param
	RI *Instruction
	PI *phys.Instruction

	particles []*Particle
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
		Shadow:      i.Shadow,
		Transparent: i.Transparent,
		P:           p,
		RI:          i,
	}

	if p.Static {
		Scene = append(Scene, r)
	} else {
		Renderables = append(Renderables, r)
		// Renderables[r] = true
	}

	return r
}

func (r *Renderable) AddShape(pi *phys.Instruction) {
	r.PI = pi
}

func (i *Instruction) createBody(p point.Param) (body *fizzle.Renderable) {
	switch i.MeshName {
	case "trail":
		body = fizzle.CreatePlaneV(mgl32.Vec3{-p.Size.Y / 2, -p.Size.X / 2, 0}, mgl32.Vec3{p.Size.Y / 2, p.Size.X / 2, 0})
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

func (r *Renderable) createRenderableShape(i *phys.Instruction) (renderShape *fizzle.Renderable) {
	switch i.ShapeType {
	case phys.ShapeType_Box:
		// shape := physShape.GetAsBox()
		// w := shape.Width
		// h := shape.Height
		renderShape = fizzle.CreateWireframeCube(-i.H, -i.W, -0.1, i.H, i.W, 0.1)
	case phys.ShapeType_Circle:
		renderShape = fizzle.CreateWireframeCircle(0, 0, 0, i.W, 32, fizzle.X|fizzle.Y)
	case phys.ShapeType_Polygon:
		log.Println("shapetype polygon not yet ready")
		// renderShape = createTriangle(o.Param.Phys.W, o.Param.Phys.H, 1)
	default:

		log.Fatalf("WARNING: shape type `%s` not yet!", i.ShapeType)
	}

	renderShape.Material = (&materials.Instruction{Name: "shape", Shader: "color", DiffColor: mgl32.Vec4{1, 0.1, 0.1, 0.75}}).Create()

	return
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

func (r *Renderable) render() {
	if r.Body == nil {
		r.Body = r.RI.createBody(r.P)
	}

	if r.Body == nil {
		log.Println("WTF?!?!", r.RI.MeshName, r.P.Pos)
		return
	}

	render.DrawRenderable(r.Body, nil, perspective, view, camera)

	if r.PI != nil {
		r.renderShape()
	}

	r.renderArts()
}

func (r *Renderable) renderShape() {
	//create shape if need
	if r.Shape == nil && r.PI != nil {
		r.Shape = r.createRenderableShape(r.PI)
	}

	// if r.physShape == nil || r.physShape.Body == nil {
	// 	//then shape or body already deleted
	// 	return
	// }

	r.Shape.Location = mgl32.Vec3{0, 0, 0}.Add(mgl32.Vec3{r.Body.Location.X(), r.Body.Location.Y(), 0})

	//set rotation
	// r.Shape.LocalRotation = mgl32.AnglesToQuat(0, 0, r.PI.Angle, 1)

	//resize width if it box... crunch!!!!
	// if r.PI.ShapeType == phys.ShapeType_Box {
	// 	r.Shape.Scale = mgl32.Vec3{1, r.PI.W / 2, 1}
	// }

	//render lines
	render.DrawLines(r.Shape, r.Shape.Material.Shader, nil, perspective, view, camera)

}

func (r *Renderable) renderArts() {
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
		if a.Body == nil {
			a.Body = a.RI.createBody(a.P)
		}

		a.Body.Location = vect.FromVec3(r.Body.Location).SubPoint(r.Angle(), vect.FromVec3(a.Pos)).Vec3()
		// a.Body.Location = mgl32.Vec3{}.Add(r.Body.Location).Add(a.Pos)

		// r.Body.LocalRotation
		a.Body.LocalRotation = mgl32.AnglesToQuat(0, 0, r.Angle()+a.Angle, 1)

		if a.Line {
			render.DrawLines(a.Body, a.Body.Material.Shader, nil, perspective, view, camera)
		} else {
			render.DrawRenderable(a.Body, nil, perspective, view, camera)
		}
	}
}

func (r *Renderable) Angle() float32 {
	q := r.Body.LocalRotation

	var ysqr = q.Y() * q.Y()

	// yaw (z-axis rotation)
	var t3 = 2.0 * float64(q.W*q.Z()+q.X()*q.Y())
	var t4 = 1.0 - 2.0*float64(ysqr+q.Z()*q.Z())
	return float32(math.Atan2(t3, t4))

	// return float32(2 * math.Acos(float64(r.Body.LocalRotation.W)))
}

func (r *Renderable) Destroy() {
	// r.needDestroy = true
	r.needDestroy = true
	for _, p := range r.particles {
		p.Destroy()
	}
	r = nil

	// r = nil

	// delete(Renderables, r)
	// delete(Renderables, r)
}

// func (r *Renderable) destroy() {
// 	// r.needDestroy = true
// 	// r.Body.Destroy()
// 	delete(Renderables, r)
// 	// r = nil
// }
