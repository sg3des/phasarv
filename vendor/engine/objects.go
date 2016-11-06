package engine

import (
	"assets"
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/graphicsprovider"
	"github.com/tbogdala/fizzle/renderer/forward"

	"param"
	"phys"
	"phys/vect"
)

var (
	Objects     = make(map[*Object]bool)
	cNewObjects = make(chan paramObject)
	chanObjects = make(chan *Object)
	Materials   = make(map[string]*fizzle.Material)
	Scene       []*Object
)

type Bullet struct {
	Parent *Object
	Param  *param.Bullet
}

type Object struct {
	Name string

	Node  *fizzle.Renderable
	Shape *phys.Shape

	RollAngle    float32
	MaxRollAngle float32

	Shadow      bool
	Transparent bool

	ArtStatic map[string]*Art
	ArtRotate map[string]*Art

	Childs map[*Object]bool

	// Callback    func(*Object, float32)
	Callbacks   map[int]Callback
	DestroyFunc func()

	Param param.Object
}

type Callback func(float32)

type Art struct {
	Name          string
	Value         float32
	MaxValue      float32
	LocalPosition mgl32.Vec3
	Art           *fizzle.Renderable
	Line          bool
}

var basicShader *fizzle.RenderShader

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

	if basicShader == nil {
		var err error
		basicShader, err = forward.CreateBasicShader()
		if err != nil {
			log.Fatalln(err)
		}
	}

	m := fizzle.NewMaterial()

	if p.Shader == "basic" {
		m.Shader = basicShader
	} else {
		m.Shader = assets.GetShader(p.Shader)
	}

	// m.Shader = assets.GetShader(p.Shader)

	m.DiffuseTex = assets.GetTexture(p.Texture).Diffuse
	m.NormalsTex = assets.GetTexture(p.Texture).Normals

	if p.DiffColor.Len() != 0 {
		m.DiffuseColor = p.DiffColor
	}

	m.SpecularColor = mgl32.Vec4{p.SpecLevel, p.SpecLevel, p.SpecLevel, 1}
	m.Shininess = p.SpecLevel

	Materials[p.Name] = m

	return m
}

// //NewPlane create renderable plane
// func NewPlane(p param.Object, w, h float32) *Object {
// 	o := &Object{
// 		Name: p.Name,
// 		// Node: fizzle.CreatePlaneXY(0, w/2, h, -w/2),
// 		Node:        fizzle.CreatePlaneV(mgl32.Vec3{0, -w / 2, 0}, mgl32.Vec3{h, w / 2, 0}),
// 		Transparent: p.Transparent,
// 	}
// 	o.Node.Material = Material(p.Material)

// 	Objects[o] = true

// 	return o
// }

//SetPhys - set physics to object
func (o *Object) SetPhys(p *param.Phys) {
	if p == nil {
		return
	}
	// log.Printf("%++v\n", p)
	o.Shape = o.NewShape(p)

	var body *phys.Body
	if p.Mass > 0 {
		o.Shape.SetElasticity(1.1)
		body = phys.NewBody(p.Mass, o.Shape.Moment(p.Mass))
		body.SetMass(p.Mass)
	} else {
		body = phys.NewBodyStatic()
	}

	body.SetPosition(vect.Vect{o.Node.Location.X(), o.Node.Location.Y()})
	body.AddShape(o.Shape)
	space.AddBody(body)

	o.Shape.Body.UserData = o
}

func (o *Object) NewShape(p *param.Phys) (shape *phys.Shape) {
	switch p.Type {
	case phys.ShapeType_Polygon:
		verts := phys.Vertices{
			vect.Vect{0, p.H / 2},
			vect.Vect{p.W / 2, -p.H / 2},
			vect.Vect{-p.W / 2, -p.H / 2},
		}
		shape = phys.NewPolygon(verts, vect.Vector_Zero)
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

//NewObject create object
func NewObject(p param.Object, arts ...param.Art) *Object {
	c := make(chan *Object)
	cNewObjects <- paramObject{p, arts, c}

	return <-c
}

type paramObject struct {
	p    param.Object
	arts []param.Art
	c    chan *Object
}

func (n paramObject) create() {
	o := &Object{
		Name:        n.p.Name,
		Node:        createNode(n.p.Mesh),
		Shadow:      n.p.Shadow,
		Transparent: n.p.Transparent,
		Param:       n.p,
		ArtRotate:   make(map[string]*Art),
		ArtStatic:   make(map[string]*Art),
		Callbacks:   make(map[int]Callback),
		Childs:      make(map[*Object]bool),
	}

	o.Node.Material = Material(n.p.Material)
	o.Node.Location = mgl32.Vec3{n.p.Pos.X, n.p.Pos.Y, n.p.Pos.Z}

	if n.p.Phys != nil {
		o.SetPhys(n.p.Phys)
		o.AddRenderShape()
	}

	// o.Childs = make(map[string]*Bar)
	for _, art := range n.arts {
		o.NewArt(art)
	}

	if n.p.Static {
		Scene = append(Scene, o)
	} else {
		Objects[o] = true
	}

	n.c <- o
}

func createNode(mesh param.Mesh) (node *fizzle.Renderable) {
	switch mesh.Model {
	case "plane":
		node = fizzle.CreatePlaneV(mgl32.Vec3{0, -mesh.X / 2, 0}, mgl32.Vec3{mesh.Y, mesh.X / 2, 0})

	// case "ground":
	// 	log.Println(mesh)
	// 	node = fizzle.CreatePlaneXY(-mesh.X, -mesh.Y, mesh.X, mesh.Y)

	case "box":
		node = fizzle.CreateCube(-2, -2, -2, 2, 2, 2)

	default:
		node = assets.GetModel(mesh.Model)
	}

	return
}

func (o *Object) AddRenderShape() {
	if o.Shape == nil {
		return
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
		renderShape = CreateTriangle(o.Param.Phys.W, o.Param.Phys.H, 1)
	default:

		log.Fatalf("WARNING: shape type `%s` not yet!", o.Shape.ShapeType())
	}

	renderShape.Material = Material(param.Material{Name: "shape", Shader: "color", DiffColor: mgl32.Vec4{1, 0.1, 0.1, 0.75}})

	o.ArtRotate["renderShape"] = &Art{Art: renderShape, Line: true}
}

// CreateTriangle wireframe triangle,[not correct]
func CreateTriangle(w, h, z float32) *fizzle.Renderable {
	const floatSize = 4
	const uintSize = 4

	r := fizzle.NewRenderable()
	r.Core = fizzle.NewRenderableCore()
	r.FaceCount = 12

	verts := [...]float32{
		0, h / 2, z,
		w / 2, -h / 2, z,
		-w / 2, -h / 2, z,
		0, h / 2, z,
	}
	indexes := [...]uint32{
		0, 1, 2, 3, 0,
	}

	r.Core.VertVBO = e.gfx.GenBuffer()
	e.gfx.BindBuffer(graphicsprovider.ARRAY_BUFFER, r.Core.VertVBO)
	e.gfx.BufferData(graphicsprovider.ARRAY_BUFFER, floatSize*len(verts), e.gfx.Ptr(&verts[0]), graphicsprovider.STATIC_DRAW)

	// create a VBO to hold the face indexes
	r.Core.ElementsVBO = e.gfx.GenBuffer()
	e.gfx.BindBuffer(graphicsprovider.ELEMENT_ARRAY_BUFFER, r.Core.ElementsVBO)
	e.gfx.BufferData(graphicsprovider.ELEMENT_ARRAY_BUFFER, uintSize*len(indexes), e.gfx.Ptr(&indexes[0]), graphicsprovider.STATIC_DRAW)

	return r
}

//NewArt to object
func (o *Object) NewArt(p param.Art) *Art {
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

	return o.applyArt(art, p)
}

//NewLine to object
func (o *Object) NewLine(p param.Art) *Art {
	art := &Art{
		Name:     p.Name,
		Value:    p.Value,
		MaxValue: p.MaxValue,
		Art:      fizzle.CreateLine(p.LocalPos.X(), p.LocalPos.Y(), 1, p.W+p.LocalPos.X(), p.H+p.LocalPos.Y(), 1),
		// Art:           fizzle.CreatePlaneXY(p.Name, 5, 1, 4, -10),
		LocalPosition: p.LocalPos,
	}

	return o.applyArt(art, p)
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

func (o *Object) applyArt(art *Art, p param.Art) *Art {
	art.Line = p.Line
	art.Art.Material = Material(p.Material)

	if o.ArtStatic == nil {
		o.ArtStatic = make(map[string]*Art)
		o.ArtRotate = make(map[string]*Art)
	}

	switch p.Type {
	case param.ArtStatic:
		o.ArtStatic[p.Name] = art
	case param.ArtRotate:
		o.ArtRotate[p.Name] = art
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

func (o *Object) GetArt(name string) (*Art, bool) {
	if art, ok := o.ArtStatic[name]; ok {
		return art, true
	}

	if art, ok := o.ArtRotate[name]; ok {
		return art, true
	}

	return nil, false
}

//NewCamera create camera set it how main camera and return it
func NewCamera(eyePos mgl32.Vec3) *fizzle.YawPitchCamera {
	e.camera = fizzle.NewYawPitchCamera(eyePos)
	return e.camera
}
