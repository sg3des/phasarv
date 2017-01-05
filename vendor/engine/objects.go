package engine

import (
	"log"
	"materials"
	"point"
	"render"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"

	"phys"
)

var (
	Objects = make(map[*Object]bool)
	Scene   []*Object
)

type Callback func(float32)

// type ObjectParam struct {
// 	Node string
// 	Material
// 	*Phys

// 	Pos  Point
// 	Rot  float32
// 	Size Point

// 	StaticType bool

// 	Shadow      bool
// 	Transparent bool
// }

type Object struct {
	Name string

	RollAngle    float32
	MaxRollAngle float32

	// Body  *fizzle.Renderable
	Shape *phys.Shape

	renderable *render.Renderable

	Childs map[*Object]bool
	Arts   map[string]*Art

	Callbacks   map[int]Callback
	DestroyFunc func()

	P  point.Param
	RI *render.Instruction
	PI *phys.Instruction
}

func (o *Object) Create(arts ...*Art) {
	o.Childs = make(map[*Object]bool)
	o.Arts = make(map[string]*Art)
	o.Callbacks = make(map[int]Callback)

	o.renderable = o.RI.Create(o.P)
	// o.Body = o.renderable.Body

	if o.PI != nil {
		o.Shape = o.PI.Create(o.P)
		o.Shape.Body.UserData = o
		o.renderable.Shape = o.createRenderableShape()
	}

	for _, a := range arts {
		a.Art = a.RI.CreateArt(a.P)
		o.renderable.AppendArt(a.Art)
	}

	if !o.P.Static {
		Objects[o] = true
	}
}

// //createPhys - set physics to object
// func (o *Object) createPhys(i *PhysInstruction) {
// 	// log.Printf("%++v\n", p)
// 	o.Shape = o.createShape(i)

// 	var body *phys.Body
// 	if i.Mass > 0 {
// 		o.Shape.SetElasticity(1.1)
// 		body = phys.NewBody(i.Mass, o.Shape.Moment(i.Mass))
// 		body.SetMass(i.Mass)
// 	} else {
// 		body = phys.NewBodyStatic()
// 	}

// 	body.SetPosition(o.Param.Pos.Vect())
// 	// if o.Node == nil {
// 	// 	body.SetPosition(vect.Vect{o.Param.Pos.X, o.Param.Pos.Y})
// 	// } else {
// 	// 	body.SetPosition(vect.Vect{o.Node.Location.X(), o.Node.Location.Y()})
// 	// }

// 	body.SetAngle(o.Param.Rot)
// 	body.AddShape(o.Shape)
// 	space.AddBody(body)

// 	o.Shape.Body.UserData = o
// }

// func (o *Object) createShape(i *PhysInstruction) (shape *phys.Shape) {
// 	switch i.ShapeType {
// 	case phys.ShapeType_Polygon:
// 		verts := phys.Vertices{
// 			vect.Vect{0, i.H / 2},
// 			vect.Vect{i.W / 2, -i.H / 2},
// 			vect.Vect{-i.W / 2, -i.H / 2},
// 		}
// 		shape = phys.NewPolygon(verts, vect.Vector_Zero)
// 	case phys.ShapeType_Box:
// 		shape = phys.NewBox(vect.Vector_Zero, i.W, i.H)
// 	case phys.ShapeType_Circle:
// 		shape = phys.NewCircle(vect.Vector_Zero, i.W)
// 	default:
// 		log.Fatalf("WARNING: shape type `%s` not yet!", i.ShapeType)
// 	}
// 	shape.Group = i.Group
// 	return
// }

// func (o *Object) Create(arts ...*Art) {
// 	// o.Body = render.NewBody(p.StaticType)
// 	// o.Body.Shadow = p.Shadow
// 	// o.Body.Transparent = p.Transparent
// 	// o.ArtRotate = make(map[string]*Art)
// 	// o.ArtStatic = make(map[string]*Art)
// 	o.Arts = make(map[string]*Art)
// 	o.Callbacks = make(map[int]Callback)
// 	o.Childs = make(map[*Object]bool)

// 	// if p.Phys != nil {
// 	// 	o.SetPhys(p.Phys)
// 	// }

// 	for _, a := range arts {
// 		o.Arts[a.Name] = a
// 	}

// 	if p.StaticType {
// 		Scene = append(Scene, o)
// 	} else {
// 		Objects[o] = true
// 	}
// }

// func (o *Object) createRenderable(p ObjectParam) {

// 	var r *fizzle.Renderable
// 	switch p.Node {
// 	case "plane":
// 		r = fizzle.CreatePlaneV(mgl32.Vec3{0, -p.Size.X / 2, 0}, mgl32.Vec3{p.Size.Y, p.Size.X / 2, 0})
// 	case "box":
// 		log.Println("warning: fixed size")
// 		r = fizzle.CreateCube(-2, -2, -2, 2, 2, 2)
// 	default:
// 		r = assets.GetModel(p.Node)
// 	}

// 	r.Material = NewMaterial(p.Material)
// 	r.Location = p.Pos.Vec3()
// 	r.LocalRotation = mgl32.AnglesToQuat(0, 0, p.Rot, 1)
// 	// o.SetRotation(p.Rot)

// 	body := render.NewBody(r, p.StaticType)
// 	body.Shadow = p.Shadow
// 	body.Transparent = p.Transparent

// 	if o.Param.Phys != nil {
// 		body.Shape = o.createRenderableShape()
// 	}

// 	o.Body = body.Body

// 	return
// }

func (o *Object) createRenderableShape() *fizzle.Renderable {
	if o.Shape == nil {
		return nil
	}

	var renderShape *fizzle.Renderable

	switch o.Shape.ShapeType() {
	case phys.ShapeType_Box:
		shape := o.Shape.GetAsBox()
		w := shape.Width
		h := shape.Height
		renderShape = fizzle.CreateWireframeCube(-h/2, -w/2, -0.1, h/2, w/2, 0.1)
	case phys.ShapeType_Circle:
		shape := o.Shape.GetAsCircle()
		renderShape = fizzle.CreateWireframeCircle(0, 0, 0, shape.Radius, 32, fizzle.X|fizzle.Y)
	case phys.ShapeType_Polygon:
		log.Println("shapetype polygon not yet ready")
		// renderShape = createTriangle(o.Param.Phys.W, o.Param.Phys.H, 1)
	default:

		log.Fatalf("WARNING: shape type `%s` not yet!", o.Shape.ShapeType())
	}

	renderShape.Material = (&materials.Instruction{Name: "shape", Shader: "color", DiffColor: mgl32.Vec4{1, 0.1, 0.1, 0.75}}).Create()

	// o.Body.Shape = renderShape
	return renderShape

	// o.ArtRotate["renderShape"] = &Art{
	// 	Name:       "shape",
	// 	Node:       renderShape,
	// 	RenderLine: true,
	// }
}
