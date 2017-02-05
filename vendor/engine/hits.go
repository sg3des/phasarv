package engine

import "phys"

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

// group 2 is players group

//Hit return object under point
func Hit(x, y float32) interface{} {
	shape := phys.Hit(x, y, 2)

	if shape != nil && shape.UserData != nil && shape.Group == 2 {
		return shape.UserData
		// switch shape.UserData.(type) {
		// 	case
		// }
		// objects = append(objects, shape.Body.UserData.(*Object))
		// object = shape.Body.UserData.(*Object)
		// if
	}

	return nil
}

// //Raycast
// func Raycast(x0, y0, x1, y1 float32, ignoreBody *phys.Body) (players []interface{}) {
// 	// r := []phys.RayCastHit{phys.RayCastHit{}}
// 	hits := phys.Hits(x0, y0, x1, y1, 2, ignoreBody)

// 	for _, hit := range hits {
// 		if hit.Shape.Group != 2 || hit.Shape.UserData == nil {
// 			continue
// 		}

// 		players = append(players, hit.Shape.UserData)
// 		// if hit.Body.UserData.(*Object).Name == "bullet" {
// 		// 	continue
// 		// }

// 		// firstpos := vect.Vect{x0, y0}
// 		// if firstpos == hit.Body.Position() {
// 		// 	continue
// 		// }

// 		// 	return hit
// 		// }
// 		// log.Println(hit.MinT, x0, y0, hit.Body.Position(), hit.Body.UserData)
// 	}

// 	return
// }
