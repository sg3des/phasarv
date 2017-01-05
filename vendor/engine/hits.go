package engine

import (
	"phys"
	"phys/vect"
)

// type PhysInstruction struct {
// 	W, H, Mass float32

// 	ShapeType phys.ShapeType
// 	Group     int
// }

// var space = phys.NewSpace()

// func init() {
// 	space.LinearDamping = 0.4
// 	space.AngularDamping = 0.2
// }

//Hit return object under point
func Hit(x, y float32) (object *Object) {
	shape := phys.Hit(x, y)
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
	hits := phys.Hits(x0, y0, x1, y1, group, ignoreBody)

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
