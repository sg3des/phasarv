package engine

import (
	"github.com/go-gl/mathgl/mgl32"
)

var Objects objects

type objects []*Object

func (objects) Add(o *Object) {
	Objects = append(Objects, o)
}

func (objects) Del(i int) {
	Objects[i] = Objects[len(Objects)-1]
	Objects[len(Objects)-1] = nil
	Objects = Objects[:len(Objects)-1]
}

func (objects) loopPhysToRender() {
	for i, o := range Objects {
		if i >= len(Objects) {
			break
		}
		if o.needDestroy {
			o.renderable.Destroy()
			Objects.Del(i)
		}
	}

	for _, o := range Objects {
		if o == nil || o.needDestroy {
			continue
		}
		if o.renderable.Body == nil {
			continue
		}

		// update position
		o.renderable.Body.Location = o.PositionVec3()
		// if o.renderable.Shape != nil {

		// }

		// update rotation
		ang := o.Rotation()
		// log.Println(o.Name, ang)

		if o.renderable.Shape != nil {
			o.renderable.Shape.LocalRotation = mgl32.AnglesToQuat(0, 0, ang, mgl32.XYZ)
		}

		//if rollAngle exist then need roll renderable object
		if o.RollAngle != 0 {
			q := mgl32.AnglesToQuat(0, 0, ang, mgl32.XYZ).Mul(mgl32.AnglesToQuat(o.RollAngle, 0, 0, mgl32.XYZ))
			o.renderable.Body.LocalRotation = q

			shape := o.shape.GetAsBox()
			// log.Println(shape.Width, shape.Height)
			shape.Width = o.P.Size.X - o.P.Size.X*o.ShapeWidthPercent()
			// log.Println(shape.Width)
			if o.renderable.Shape != nil {

				o.renderable.Shape.Scale = mgl32.Vec3{o.P.Size.Y, shape.Width, 1}
				// log.Println(o.renderable.Shape.Scale)

			}
			shape.UpdatePoly()
		} else {
			o.renderable.Body.LocalRotation = mgl32.AnglesToQuat(0, 0, ang, mgl32.XYZ)
		}

		// for _, art := range o.Arts {
		// 	art.Art.Angle = 1
		// }
	}
}

func (objects) loopCallbacks(dt float32) {
	for _, o := range Objects {
		if o == nil || o.needDestroy {
			// Objects.Del(i)
			// o.renderable
			continue
		}
		for _, f := range o.callbacks {
			f(dt)
		}
	}
}
