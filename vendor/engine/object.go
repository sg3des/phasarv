package engine

import (
	"log"
	"point"
	"render"

	"github.com/go-gl/mathgl/mgl32"

	"phys"
)

//Object is universal object type union 3d renderable object and 2d physic object
type Object struct {
	Name string

	RollAngle    float32
	MaxRollAngle float32

	shape       *phys.Shape
	renderable  *render.Renderable
	needDestroy bool

	Childs map[*Object]bool
	Arts   []*Art

	callbacks   []func(float32)
	destroyFunc func()

	P  *point.Param
	RI *render.Instruction
	PI *phys.Instruction

	UserData interface{}
}

//Create object by instructons
func (o *Object) Create(arts ...*Art) {
	if NeedRender {
		o.renderable = o.RI.Create(o.P)
		for _, a := range arts {
			if a != nil {
				o.AppendArt(a)
			}
		}
	}

	if o.PI != nil {
		o.shape = o.PI.Create(o.P)
		o.shape.UserData = o
		if NeedRender {
			// log.Println(o.Name, o.PI.ShapeType)
			o.renderable.AddShape(o.PI)
		}
	}

	if !o.P.Static {
		Objects.Add(o)
	}
}

func (o *Object) SetUserData(i interface{}) {
	// if o.shape == nil {
	// 	log.Fatalln("shape not created, name:", o.Name)
	// }
	// o.shape.UserData = i
	o.UserData = i
}

func (o *Object) AddCallback(fs ...func(float32)) {
	for _, f := range fs {
		o.callbacks = append(o.callbacks, f)
	}
}

func (o *Object) SetCallbackCollision(f func(arb *phys.Arbiter) bool) {
	if o.shape != nil && o.shape.Body != nil {
		o.shape.Body.CallBackCollision = f
	} else {
		log.Println("warning! shape or shape.body is nil of object, ", o.Name)
	}
}

func (o *Object) AddChild(child *Object) {
	if o.Childs == nil {
		o.Childs = make(map[*Object]bool)
	}
	o.Childs[child] = true
}

func (o *Object) AddTrail(offset mgl32.Vec3, count int, size point.P) {
	o.renderable.NewTrail(offset, count, size)
}

func (o *Object) SetDestroyFunc(f func()) {
	o.destroyFunc = f
}

func (o *Object) Destroy() {
	if o.destroyFunc != nil {
		o.destroyFunc()
		return
	}

	// o.needDestroy = true

	if o.shape != nil {
		o.shape.Body.Enabled = false
		phys.RemoveBody(o.shape.Body)
		// space.RemoveShape(o.Shape) - crash need TODO
	}

	o.renderable.Destroy()
	o.needDestroy = true

	// Objects[o] = false

	for child := range o.Childs {
		child.Destroy()
	}

	// o.renderable.Destroy()
	// delete(Objects, o)
	// o = nil
}

func (o *Object) Remove() {
	// o.needDestroy = true

	if o.shape != nil {
		o.shape.Body.Enabled = false
		phys.RemoveBody(o.shape.Body)
		// space.RemoveShape(o.Shape) - crash need TODO
	}

	o.renderable.Destroy()
	o.needDestroy = true

	// Objects[o] = false

	for child := range o.Childs {
		child.Remove()
	}

}

// func (o *Object) Clone() *Object {
// 	newObject := &Object{
// 		Name:         o.Name,
// 		Node:         o.Node.Clone(),
// 		MaxRollAngle: o.MaxRollAngle,

// 		Shadow:      o.Shadow,
// 		Transparent: o.Transparent,

// 		ArtStatic: o.ArtStatic,
// 		ArtRotate: o.ArtRotate,

// 		// Callback:    o.Callback,
// 		Callbacks:   o.Callbacks,
// 		DestroyFunc: o.DestroyFunc,

// 		Param: o.Param,
// 	}

// 	newObject.SetPhys(o.Param.Phys)
// 	newObject.Node.Material = NewMaterial(o.Param.Material)

// 	Objects[newObject] = true

// 	return newObject
// }
