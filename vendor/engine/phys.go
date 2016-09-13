package engine

import (
	"phys/vect"
	"time"

	"phys"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	space *phys.Space

	Movement func(float32)
	Control  func()

	physT0 time.Time
	physT1 time.Time
	physDT float32
)

func InitPhys(airResist float32) {
	space = phys.NewSpace()
	space.LinearDamping = airResist
	space.AngularDamping = airResist

	physT0 = time.Now()
}

func physFrame() {
	physT1 = time.Now()
	physDT = float32(physT1.Sub(physT0).Seconds())
	physT0 = physT1

	// log.Println(physDT)

	if Control != nil {
		Control()
	}
	if Movement != nil {
		Movement(physDT)
	}

	for o := range Objects {
		if o.Callback != nil {
			o.Callback(o, physDT)
		}
	}

	space.Step(physDT)
}

// //milk emulates the air resistance
// func milk(o *Object, dt float32) {
// 	// vel := o.Shape.Body.Velocity()
// 	// vel2 := o.Shape.Body.Velocity()

// 	// vel2.Mult(dt * AirResist / o.Shape.Body.Mass())

// 	// log.Println(space.LinearDamping)

// 	// vel.Sub(vel2)
// 	// o.Shape.Body.SetVelocity(vel.X, vel.Y)

// 	// avel := o.Shape.Body.AngularVelocity()
// 	// o.Shape.Body.SetAngularVelocity(avel - avel*dt)
// }

func physRender() {
	for o, b := range Objects {
		if o.Shape == nil || o.Node == nil || !b {
			continue
		}

		// update position
		pos := o.Shape.Body.Position()
		o.Node.Location = mgl32.Vec3{pos.X, pos.Y, 0}

		// update rotation
		ang := o.Shape.Body.Angle()
		if o.Player != nil {
			q := mgl32.AnglesToQuat(0, 0, ang, 1).Mul(mgl32.AnglesToQuat(o.Player.SubAngle, 0, 0, 1))
			o.Node.LocalRotation = q
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
