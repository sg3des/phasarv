package phys

import (
	"log"
	"phys/vect"
	"point"
)

var space *Space

func Init() {
	space = NewSpace()
	space.LinearDamping = 0.4
	space.AngularDamping = 0.2
}

func NextFrame(dt float32) {
	space.Step(dt)
}

type Instruction struct {
	W, H, Mass float32

	ShapeType
	Group int
}

func (i *Instruction) Create(p point.Param) *Shape {
	shape := i.createShape()

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

func (i *Instruction) createShape() (shape *Shape) {
	switch i.ShapeType {
	case ShapeType_Polygon:
		verts := Vertices{
			vect.Vect{0, i.H / 2},
			vect.Vect{i.W / 2, -i.H / 2},
			vect.Vect{-i.W / 2, -i.H / 2},
		}
		shape = NewPolygon(verts, vect.Vector_Zero)
	case ShapeType_Box:
		shape = NewBox(vect.Vector_Zero, i.W, i.H)
	case ShapeType_Circle:
		shape = NewCircle(vect.Vector_Zero, i.W)
	default:
		log.Fatalf("WARNING: shape type `%s` not yet!", i.ShapeType)
	}
	shape.Group = i.Group
	// shape.SetElasticity(1.1)
	return
}

func Hit(x, y float32) *Shape {
	return space.SpacePointQueryFirst(vect.Vect{x, y}, 0, 0, false)
}

func Hits(x0, y0, x1, y1 float32, group int, ignoreBody *Body) []*RayCastHit {
	return space.RayCastAll(vect.Vect{x0, y0}, vect.Vect{x1, y1}, group, ignoreBody)
}

func RemoveBody(body *Body) {
	space.RemoveBody(body)
}
