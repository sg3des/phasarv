package render

import (
	"assets"
	"log"
	"materials"
	"math"
	"phys"
	"point"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
)

type Renderable struct {
	Body      *fizzle.Renderable
	Shape     *fizzle.Renderable
	physShape *phys.Shape

	Shadow      bool
	Transparent bool

	ArtStatic []*Art
	ArtRotate []*Art

	arts []*Art

	P  point.Param
	RI *Instruction
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
		Renderables[r] = true
	}

	return r
}

func (r *Renderable) AddShape(physShape *phys.Shape) {
	r.physShape = physShape
}

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

func (r *Renderable) createRenderableShape(physShape *phys.Shape) (renderShape *fizzle.Renderable) {
	switch physShape.ShapeType() {
	case phys.ShapeType_Box:
		shape := physShape.GetAsBox()
		w := shape.Width
		h := shape.Height
		renderShape = fizzle.CreateWireframeCube(-h/2, -w/2, -0.1, h/2, w/2, 0.1)
	case phys.ShapeType_Circle:
		shape := physShape.GetAsCircle()
		renderShape = fizzle.CreateWireframeCircle(0, 0, 0, shape.Radius, 32, fizzle.X|fizzle.Y)
	case phys.ShapeType_Polygon:
		log.Println("shapetype polygon not yet ready")
		// renderShape = createTriangle(o.Param.Phys.W, o.Param.Phys.H, 1)
	default:

		log.Fatalf("WARNING: shape type `%s` not yet!", physShape.ShapeType())
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

func (r *Renderable) Render() {
	if r.Body == nil {
		// log.Println("create ", r.RI.MeshName)
		r.Body = r.RI.createBody(r.P)
	}

	// if r.needDestroy {
	// 	// log.Println("destroy", r.RI.MeshName)
	// 	r.destroy()
	// 	return
	// }

	if r.Body == nil {
		log.Println("WTF?!?!", r.RI.MeshName, r.P.Pos)
		return
	}

	render.DrawRenderable(r.Body, nil, perspective, view, camera)

	if r.physShape != nil {
		r.renderShape()
	}

	r.renderArts()
}

func (r *Renderable) renderShape() {
	//create shape if need
	if r.Shape == nil {
		r.Shape = r.createRenderableShape(r.physShape)
	}

	if r.physShape == nil || r.physShape.Body == nil {
		//then shape or body already deleted
		return
	}

	r.Shape.Location = mgl32.Vec3{0, 0, 0}.Add(mgl32.Vec3{r.Body.Location.X(), r.Body.Location.Y(), 0})

	//set rotation
	r.Shape.LocalRotation = mgl32.AnglesToQuat(0, 0, r.physShape.Body.Angle(), 1)

	//resize width if it box... crunch!!!!
	if r.physShape.ShapeType() == phys.ShapeType_Box {
		r.Shape.Scale = mgl32.Vec3{1, r.physShape.GetAsBox().Width / 2, 1}
	}

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
	// r.needDestroy = true

	delete(Renderables, r)
	// delete(Renderables, r)
}

// func (r *Renderable) destroy() {
// 	// r.needDestroy = true
// 	// r.Body.Destroy()
// 	delete(Renderables, r)
// 	// r = nil
// }
