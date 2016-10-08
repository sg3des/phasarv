package engine

import (
	"phys"
	"phys/vect"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	space                = phys.NewSpace()
	MaxRollAngle float32 = 1.5
)

func InitPhys(airResist float32) {
	// space = phys.NewSpace()
	space.LinearDamping = airResist
	space.AngularDamping = airResist
}

func physRender(dt float32) {
	space.Step(dt)

	for o, b := range Objects {
		if o.Shape == nil || o.Node == nil || !b {
			continue
		}

		// update position
		pos := o.Shape.Body.Position()
		o.Node.Location = mgl32.Vec3{pos.X, pos.Y, o.Node.Location.Z()}

		// update rotation
		ang := o.Shape.Body.Angle()
		if o.RollAngle != 0 {
			q := mgl32.AnglesToQuat(0, 0, ang, 1).Mul(mgl32.AnglesToQuat(o.RollAngle, 0, 0, 1))
			o.Node.LocalRotation = q

			shape := o.Shape.GetAsBox()
			shape.Width = o.Param.Phys.W - o.Param.Phys.W*o.ShapeWidthPercent()
			shape.UpdatePoly()
		} else {
			o.Node.LocalRotation = mgl32.AnglesToQuat(0, 0, ang, 1)
		}
	}
}

//Hit return object under point
func Hit(x, y float32) (object *Object) {
	shape := space.SpacePointQueryFirst(vect.Vect{x, y}, 0, 0, false)
	// for _, shape := range shapes {
	if shape != nil && shape.Body != nil && shape.Body.UserData != nil {
		// objects = append(objects, shape.Body.UserData.(*Object))
		object = shape.Body.UserData.(*Object)
	}
	// }
	return
}

//Raycast
func Raycast(x0, y0, x1, y1 float32, group int, ignoreBody *phys.Body) *phys.RayCastHit {
	// r := []phys.RayCastHit{phys.RayCastHit{}}
	hits := space.RayCastAll(vect.Vect{x0, y0}, vect.Vect{x1, y1}, group, ignoreBody)

	for _, hit := range hits {
		if hit.Body.UserData != nil {
			if hit.Body.UserData.(*Object).Name == "bullet" {
				continue
			}

			firstpos := vect.Vect{x0, y0}
			if firstpos == hit.Body.Position() {
				continue
			}

			return hit
		}
		// log.Println(hit.MinT, x0, y0, hit.Body.Position(), hit.Body.UserData)
	}

	return nil
}
