package engine

import (
	"assets"
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"

	"param"
	"phys"
	"phys/vect"
)

var (
	Objects   = make(map[*Object]bool)
	Materials = make(map[string]*fizzle.Material)
)

type Bullet struct {
	Parent *Object
	Param  *param.Bullet
}

type Object struct {
	Name string

	Node  *fizzle.Renderable
	Shape *phys.Shape

	RollAngle float32

	Shadow      bool
	Transparent bool

	ArtStatic map[string]*Art
	ArtRotate map[string]*Art

	Callback func(*Object, float32)

	Param param.Object
}

type Art struct {
	Name          string
	Value         float32
	MaxValue      float32
	LocalPosition mgl32.Vec3
	Art           *fizzle.Renderable
	Line          bool
}

func Material(p param.Material) *fizzle.Material {
	// m, ok := Materials[p.Name]
	// if ok {
	// 	return m
	// }

	if p.Shader == "" {
		p.Shader = "color"
	}
	if p.Texture == "" {
		p.Texture = "gray"
	}

	m := fizzle.NewMaterial()
	m.Shader = assets.Shaders[p.Shader]

	m.DiffuseTex = assets.Textures[p.Texture].Diffuse
	m.NormalsTex = assets.Textures[p.Texture].Normals

	if p.DiffColor.Len() != 0 {
		m.DiffuseColor = p.DiffColor
	}

	m.SpecularColor = mgl32.Vec4{p.SpecLevel, p.SpecLevel, p.SpecLevel, 1}
	m.Shininess = p.SpecLevel

	Materials[p.Name] = m

	return m
}

//NewPlane create renderable plane
func NewPlane(p param.Object, w, h float32) *Object {
	e := &Object{
		Name: p.Name,
		// Node: fizzle.CreatePlaneXY(0, w/2, h, -w/2),
		Node:        fizzle.CreatePlaneV(mgl32.Vec3{0, -w / 2, 0}, mgl32.Vec3{h, w / 2, 0}),
		Transparent: p.Transparent,
	}
	e.Node.Material = Material(p.Material)

	Objects[e] = true

	return e
}

//NewPlanePoint create renderable plane from two points
func NewPlanePoint(p param.Object, p0, p1 mgl32.Vec3) *Object {
	e := &Object{
		Name:        p.Name,
		Node:        fizzle.CreatePlaneV(p0, p1),
		Transparent: p.Transparent,
	}
	e.Node.Material = Material(p.Material)

	Objects[e] = true

	return e
}

//SetPhys - set physics to object
func (e *Object) SetPhys(p *param.Phys) {
	if p == nil {
		return
	}
	// log.Printf("%++v\n", p)
	e.Shape = e.NewShape(p)

	var body *phys.Body
	if p.Mass > 0 {
		e.Shape.SetElasticity(2)
		body = phys.NewBody(p.Mass, e.Shape.Moment(p.Mass))
		body.SetMass(p.Mass)
	} else {
		body = phys.NewBodyStatic()
	}

	body.SetPosition(vect.Vect{e.Node.Location.X(), e.Node.Location.Y()})
	body.AddShape(e.Shape)
	space.AddBody(body)

	e.Shape.Body.UserData = e
}

func (e *Object) NewShape(p *param.Phys) (shape *phys.Shape) {
	switch p.Type {
	case phys.ShapeType_Box:
		shape = phys.NewBox(vect.Vector_Zero, p.W, p.H)
	case phys.ShapeType_Circle:
		shape = phys.NewCircle(vect.Vector_Zero, p.W)
	default:
		log.Fatalf("WARNING: shape type `%s` not yet!", p.Type)
	}
	shape.Group = p.Group
	return
}

//NewBox generated mesh box with shader diffuse_texbumped and TestCube texture
func NewBox(name string) *Object {
	e := &Object{
		Name: name,
		Node: fizzle.CreateCube(-0.5, -0.5, -0.5, 0.5, 0.5, 0.5),
	}
	e.Node.Material = Material(param.Material{Name: "box", Shader: "diffuse_texbumped_shadows", Texture: "TestCube"})

	Objects[e] = true
	return e
}

//NewObject create object
func NewObject(p param.Object, arts ...param.Art) *Object {
	e := &Object{
		Name:        p.Name,
		Node:        assets.GetModel(p.Mesh.Model),
		Shadow:      p.Mesh.Shadow,
		Transparent: p.Transparent,
		Param:       p,
	}

	e.ArtRotate = make(map[string]*Art)
	e.ArtStatic = make(map[string]*Art)

	e.Node.Material = Material(p.Material)

	e.Node.Location = mgl32.Vec3{p.Pos.X, p.Pos.Y, p.Pos.Z}

	if p.Phys != nil {
		e.SetPhys(p.Phys)
		e.AddRenderShape()
	}

	// e.Childs = make(map[string]*Bar)
	for _, art := range arts {
		e.NewArt(art)
	}

	Objects[e] = true

	return e
}

func (e *Object) AddRenderShape() {
	if e.Shape == nil {
		return
	}

	var renderShape *fizzle.Renderable

	switch e.Shape.ShapeType() {
	case phys.ShapeType_Box:
		shape := e.Shape.GetAsBox()
		w := shape.Width
		h := shape.Height
		renderShape = fizzle.CreateWireframeCube(-h/2, -w/2, -0.1, h/2, w/2, 0.1)
	case phys.ShapeType_Circle:
		shape := e.Shape.GetAsCircle()
		renderShape = fizzle.CreateWireframeCircle(0, 0, 0, shape.Radius, 32, fizzle.X|fizzle.Y)
	default:
		log.Fatalf("WARNING: shape type `%s` not yet!", e.Shape.ShapeType())
	}

	renderShape.Material = Material(param.Material{Name: "shape", Shader: "color", DiffColor: mgl32.Vec4{1, 0.1, 0.1, 0.75}})

	e.ArtRotate["renderShape"] = &Art{Art: renderShape, Line: true}
}

//NewArt to object
func (e *Object) NewArt(p param.Art) *Art {
	art := &Art{
		Name:     p.Name,
		Value:    p.Value,
		MaxValue: p.MaxValue,
		Art:      fizzle.CreatePlaneV(mgl32.Vec3{}, mgl32.Vec3{p.W, p.H}),
		// Art:      fizzle.CreatePlaneXY(p.LocalPos.X(), p.LocalPos.Y(), p.W+p.LocalPos.X(), p.H+p.LocalPos.Y()),
		// Art:           fizzle.CreatePlaneXY(p.Name, 5, 1, 4, -10),
		LocalPosition: p.LocalPos,
	}

	art.Art.Material = Material(p.Material)

	return e.applyArt(art, p)
}

//NewLine to object
func (e *Object) NewLine(p param.Art) *Art {
	art := &Art{
		Name:     p.Name,
		Value:    p.Value,
		MaxValue: p.MaxValue,
		Art:      fizzle.CreateLine(p.LocalPos.X(), p.LocalPos.Y(), 1, p.W+p.LocalPos.X(), p.H+p.LocalPos.Y(), 1),
		// Art:           fizzle.CreatePlaneXY(p.Name, 5, 1, 4, -10),
		LocalPosition: p.LocalPos,
	}

	return e.applyArt(art, p)
}

func (e *Object) NewCircle(p param.Art) *Art {
	art := &Art{
		Name:     p.Name,
		Value:    p.Value,
		MaxValue: p.MaxValue,
		Art: fizzle.CreateWireframeCircle(p.LocalPos.X(), p.LocalPos.Y(), p.LocalPos.Z(), 0.5, 64, fizzle.X|
			fizzle.Y),
		LocalPosition: p.LocalPos,
	}

	return e.applyArt(art, p)
}

func (e *Object) applyArt(art *Art, p param.Art) *Art {
	art.Line = p.Line
	art.Art.Material = Material(p.Material)

	if e.ArtStatic == nil {
		e.ArtStatic = make(map[string]*Art)
		e.ArtRotate = make(map[string]*Art)
	}

	switch p.Type {
	case param.ArtStatic:
		e.ArtStatic[p.Name] = art
	case param.ArtRotate:
		e.ArtRotate[p.Name] = art
	}

	return art
}

func NewHealthBar(value float32) param.Art {
	return param.Art{
		Name:     "health",
		Value:    value,
		MaxValue: value,
		W:        2,
		H:        0.2,
		LocalPos: mgl32.Vec3{0, 1, 1.1},
		Type:     param.ArtStatic,

		Material: param.Material{Name: "healthBar", DiffColor: mgl32.Vec4{0, 0.6, 0, 0.7}},
	}
}

//Resize bar
func (b *Art) Resize() {
	if b.Line {
		b.Art.FaceCount = uint32(b.MaxValue * b.Value)
	} else {
		percent := b.Value / b.MaxValue
		b.Art.Scale = mgl32.Vec3{percent, 1, 1}
	}
}

func (e *Object) GetArt(name string) (*Art, bool) {
	if art, ok := e.ArtStatic[name]; ok {
		return art, true
	}

	if art, ok := e.ArtRotate[name]; ok {
		return art, true
	}

	return nil, false
}

//NewCamera create camera set it how main camera and return it
func NewCamera(eyePos mgl32.Vec3) *fizzle.YawPitchCamera {
	engine.Camera = fizzle.NewYawPitchCamera(eyePos)
	return engine.Camera
}
