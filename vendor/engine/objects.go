package engine

import (
	"assets"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"
	"github.com/tbogdala/fizzle/graphicsprovider"

	"param"
	"phys"
	"phys/vect"
)

var (
	Objects = make(map[*Object]bool)
	Lines   = make(map[*fizzle.Renderable]bool)
)

type Bullet struct {
	Parent *Object
	Param  *param.Bullet
}

type Object struct {
	Name string

	Node        *fizzle.Renderable
	Shape       *phys.Shape
	Shadow      bool
	Transparent bool

	ArtStatic map[string]*Art
	ArtRotate map[string]*Art

	RollAngle float32

	// Player *param.Player
	// Bullet *Bullet

	Callback func(*Object, float32)
}

type Art struct {
	Name          string
	Value         float32
	MaxValue      float32
	LocalPosition mgl32.Vec3
	Art           *fizzle.Renderable
	Line          bool
}

//NewPlane create renderable plane
func NewPlane(p param.Object, w, h float32) *Object {
	e := &Object{
		Name: p.Name,
		// Node: fizzle.CreatePlaneXY(0, w/2, h, -w/2),
		Node:        fizzle.CreatePlaneV(mgl32.Vec3{0, -w / 2, 0}, mgl32.Vec3{h, w / 2, 0}),
		Transparent: p.Transparent,
	}

	e.Node.Core.Shader = assets.Shaders[p.Mesh.Shader]
	e.Node.Core.DiffuseColor = mgl32.Vec4{1, 1, 1, 1}
	e.Node.Core.SpecularColor = mgl32.Vec4{0.3, 0.3, 0.3, 1.0}

	e.Node.Core.Shininess = 6.0
	e.Node.Core.Tex[0] = assets.Textures[p.Mesh.Texture].Diffuse
	e.Node.Core.Tex[1] = assets.Textures[p.Mesh.Texture].Normals

	// e.Node.Location = mgl32.Vec3{p.Pos.X, p.Pos.Y, p.Pos.Z}

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

	e.Node.Core.Shader = assets.Shaders[p.Mesh.Shader]
	e.Node.Core.DiffuseColor = mgl32.Vec4{1, 1, 1, 1}
	e.Node.Core.SpecularColor = mgl32.Vec4{0.3, 0.3, 0.3, 1.0}

	e.Node.Core.Shininess = 6.0
	e.Node.Core.Tex[0] = assets.Textures[p.Mesh.Texture].Diffuse
	e.Node.Core.Tex[1] = assets.Textures[p.Mesh.Texture].Normals

	// e.Node.Location = mgl32.Vec3{p.Pos.X, p.Pos.Y, p.Pos.Z}

	Objects[e] = true

	return e
}

//SetPhys - set physics to object
func (e *Object) SetPhys(p param.Phys) {
	if p.Mass == 0 {
		p.Mass = 1
	}

	e.Shape = phys.NewBox(vect.Vector_Zero, p.W, p.H)
	e.Shape.Group = p.Group
	// e.Shape = phys.NewPolygon(phys.Vertices{
	// 	vect.Vect{-1, -1},
	// 	vect.Vect{-1, 1},
	// 	vect.Vect{1, 1},
	// 	vect.Vect{1, -1},
	// }, vect.Vector_Zero)

	body := phys.NewBody(p.Mass, e.Shape.Moment(p.Mass))
	body.SetMass(p.Mass)
	body.AddShape(e.Shape)

	space.AddBody(body)

	pos := e.Node.Location
	// log.Fatal(x, y, e.Name)
	e.Shape.Body.SetPosition(vect.Vect{pos.X(), pos.Y()})
	e.Shape.Body.UserData = e
}

//NewHideBox create not renderable box
func NewHideBox(p param.Object) *Object {
	e := &Object{
		Name: p.Name,
		Node: fizzle.CreateCube(0, 0, 0, 1, 1, 1),
	}

	e.SetPhys(p.PH)

	Objects[e] = false

	return e
}

//NewBox generated mesh box with shader diffuse_texbumped and TestCube texture
func NewBox(name string) *Object {
	e := &Object{
		Name: name,
		Node: fizzle.CreateCube(-0.5, -0.5, -0.5, 0.5, 0.5, 0.5),
	}

	e.Node.Core.Shader = assets.Shaders["diffuse_texbumped_shadows"]
	e.Node.Core.DiffuseColor = mgl32.Vec4{1.0, 1.0, 1.0, 1.0}
	e.Node.Core.SpecularColor = mgl32.Vec4{0.3, 0.3, 0.3, 1.0}

	e.Node.Core.Shininess = 6.0
	e.Node.Core.Tex[0] = assets.Textures["TestCube"].Diffuse
	e.Node.Core.Tex[1] = assets.Textures["TestCube"].Normals

	e.Node.Location = mgl32.Vec3{0, 0, 0}

	// fmt.Println(e.Node)
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
	}

	e.Node.Core.Shader = assets.Shaders[p.Mesh.Shader]
	e.Node.Core.DiffuseColor = mgl32.Vec4{1.0, 1.0, 1.0, 1.0}
	e.Node.Core.SpecularColor = mgl32.Vec4{0.3, 0.3, 0.3, 1.0}

	e.Node.Core.Shininess = 6.0
	e.Node.Core.Tex[0] = assets.Textures[p.Mesh.Texture].Diffuse
	e.Node.Core.Tex[1] = assets.Textures[p.Mesh.Texture].Normals

	// log.Println(p.Pos.X, p.Pos.Y, p.Name)
	// e.Node.Location.Add(mgl32.Vec3{p.Pos.X, p.Pos.Y, p.Pos.Z})
	e.Node.Location = mgl32.Vec3{p.Pos.X, p.Pos.Y, p.Pos.Z}

	if p.PH.Mass > 0 {
		e.SetPhys(p.PH)
		e.Shape.SetElasticity(0.95)
	}

	// e.Childs = make(map[string]*Bar)
	for _, art := range arts {
		e.NewArt(art)
	}

	Objects[e] = true

	return e
}

// func CreateCurve(axis int) *fizzle.Renderable {
// 	// e := &Object{Name: "wireframe"}
// 	line := fizzle.CreateWireframeCircle(0, 0, 0, 3, 32, axis)

// 	line.Core.Shader = assets.Shaders["color"]
// 	line.Core.DiffuseColor = mgl32.Vec4{1.0, 1.0, 1.0, 1.0}
// 	line.Core.SpecularColor = mgl32.Vec4{0.3, 0.3, 0.3, 1.0}

// 	// e.Node.Core.Shininess = 6.0
// 	// e.Node.Core.Tex0 = assets.Textures["TestCube"].Diffuse
// 	// e.Node.Core.Tex1 = assets.Textures["TestCube"].Normals

// 	// line.Location = mgl32.Vec3{3, 13, 1}

// 	Lines[line] = true
// 	return line
// }

//NewArt to object
func (e *Object) NewArt(p param.Art) *Art {
	// p0 := p.LocalPos.Add()
	// p1 := p.LocalPos.Add(mgl32.Vec3{p.W / 2, p.H / 2, 0})

	art := &Art{
		Name:     p.Name,
		Value:    p.Value,
		MaxValue: p.MaxValue,
		Art:      fizzle.CreatePlaneV(mgl32.Vec3{}, mgl32.Vec3{p.W, p.H}),
		// Art:      fizzle.CreatePlaneXY(p.LocalPos.X(), p.LocalPos.Y(), p.W+p.LocalPos.X(), p.H+p.LocalPos.Y()),
		// Art:           fizzle.CreatePlaneXY(p.Name, 5, 1, 4, -10),
		LocalPosition: p.LocalPos,
	}

	if p.Texture == "" {
		art.Art.Core.Tex[0] = assets.Textures["gray"].Diffuse
		art.Art.Core.Tex[1] = assets.Textures["gray"].Normals
	} else {
		art.Art.Core.Tex[0] = assets.Textures[p.Texture].Diffuse
		art.Art.Core.Tex[1] = assets.Textures[p.Texture].Normals
	}

	if p.Shader != "" {
		art.Art.Core.Shader = assets.Shaders[p.Shader]
	}

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

//NewCurve to object
func (e *Object) NewCurve(p param.Art) *Art {
	curve := fizzle.NewRenderable()
	curve.Core = fizzle.NewRenderableCore()

	verts := []float32{}
	indexes := []uint32{}

	radsPerSeg := math.Pi * 2.0 / float64(p.Seg)
	for i := 0; i <= int(p.MaxValue); i++ {
		verts = append(verts, p.W*float32(math.Cos(radsPerSeg*float64(i))))
		verts = append(verts, p.H*(p.W*float32(math.Sin(radsPerSeg*float64(i)))))
		verts = append(verts, 1)

		indexes = append(indexes, uint32(i))
		if i != int(p.MaxValue)-1 {
			indexes = append(indexes, uint32(i)+1)
		}
	}
	curve.FaceCount = uint32(p.MaxValue)

	// calculate the memory size of floats used to calculate total memory size of float arrays
	const floatSize = 4
	const uintSize = 4

	// create a VBO to hold the vertex data
	curve.Core.VertVBO = gfx.GenBuffer()
	gfx.BindBuffer(graphicsprovider.ARRAY_BUFFER, curve.Core.VertVBO)
	gfx.BufferData(graphicsprovider.ARRAY_BUFFER, floatSize*len(verts), gfx.Ptr(&verts[0]), graphicsprovider.STATIC_DRAW)

	// create a VBO to hold the face indexes
	curve.Core.ElementsVBO = gfx.GenBuffer()
	gfx.BindBuffer(graphicsprovider.ELEMENT_ARRAY_BUFFER, curve.Core.ElementsVBO)
	gfx.BufferData(graphicsprovider.ELEMENT_ARRAY_BUFFER, uintSize*len(indexes), gfx.Ptr(&indexes[0]), graphicsprovider.STATIC_DRAW)

	art := &Art{
		Name:          p.Name,
		Value:         p.Value,
		MaxValue:      p.MaxValue,
		Art:           curve,
		LocalPosition: p.LocalPos,
	}

	// art.Art.FaceCount = uint32(p.Seg)

	// art.Art.LocalRotation = mgl32.AnglesToQuat(1.5704, 0, 0, 1)
	// art.Art.FaceCount = 16

	return e.applyArt(art, p)
}

func (e *Object) applyArt(art *Art, p param.Art) *Art {
	art.Line = p.Line

	if art.Art.Core.Shader == nil {
		art.Art.Core.Shader = assets.Shaders["color"]
	}
	art.Art.Core.DiffuseColor = p.Color

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
		Color:    mgl32.Vec4{0, 0.6, 0, 0.7},
		LocalPos: mgl32.Vec3{0, 1, 1.1},
		Type:     param.ArtStatic,
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
	Camera = fizzle.NewYawPitchCamera(eyePos)
	return Camera
}
