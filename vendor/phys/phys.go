package phys

import (
	"log"
	"phys/vect"
	"point"
)

const GROUP_PLAYER = 0
const GROUP_STATIC = 1
const GROUP_BULLET = 2

var space *Space

func Init() {
	space = NewSpace()
	space.LinearDamping = 0.4
	space.AngularDamping = 0.2
}

func NextFrame(dt float32) {
	if space == nil {
		return
	}
	space.Step(dt)
}

type Instruction struct {
	Mass      float32
	ShapeType ShapeType
	Group     int
}

func (i *Instruction) Create(p *point.Param) *Shape {
	shape := i.createShape(p.Size)

	var body *Body
	if i.Mass > 0 {
		body = NewBody(i.Mass, shape.Moment(i.Mass))
		body.SetMass(i.Mass)
	} else {
		body = NewBodyStatic()
	}

	// log.Println(p.Pos)
	body.SetPosition(p.Pos.Vect())
	body.SetAngle(p.Angle)

	body.AddShape(shape)
	space.AddBody(body)

	return shape
}

func (i *Instruction) createShape(p point.P) (shape *Shape) {
	switch i.ShapeType {
	case ShapeType_Polygon:
		verts := Vertices{
			vect.Vect{0, p.X / 2},
			vect.Vect{p.X / 2, -p.Y / 2},
			vect.Vect{-p.X / 2, -p.Y / 2},
		}
		shape = NewPolygon(verts, vect.Vector_Zero)
	case ShapeType_Box:
		shape = NewBox(vect.Vector_Zero, p.X, p.Y)
	case ShapeType_Circle:
		shape = NewCircle(vect.Vector_Zero, p.X)
	default:
		log.Fatalf("WARNING: shape type `%s` not yet!", i.ShapeType.ToString())
	}
	shape.Group = i.Group
	// shape.SetElasticity(1.1)
	return
}

func Hit(x, y float32, group int) *Shape {
	return space.SpacePointQueryFirst(vect.Vect{x, y}, 0, group, false)
}

func Hits(x0, y0, x1, y1 float32, group int, ignoreBody *Body) []*RayCastHit {
	return space.RayCastAll(vect.Vect{x0, y0}, vect.Vect{x1, y1}, group, ignoreBody)
}

func RemoveBody(body *Body) {
	space.RemoveBody(body)
}

func RemoveShape(shape *Shape) {
	space.RemoveShape(shape)
}
