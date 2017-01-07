package engine

import (
	"log"
	"point"
	"render"

	"phys"
)

var (
	//dynamic objects
	Objects = make(map[*Object]bool)
)

//Callback function type
type Callback func(float32)

//Object is universal object type union 3d renderable object and 2d physic object
type Object struct {
	Name string

	RollAngle    float32
	MaxRollAngle float32

	Shape       *phys.Shape
	renderable  *render.Renderable
	needDestroy bool

	Childs map[*Object]bool
	Arts   map[string]*Art

	Callbacks   map[int]Callback
	DestroyFunc func()

	P  point.Param
	RI *render.Instruction
	PI *phys.Instruction
}

//Create object by instructons
func (o *Object) Create(arts ...*Art) {
	o.renderable = o.RI.Create(o.P)

	if o.PI != nil {
		o.Shape = o.PI.Create(o.P)
		// o.Shape.Body.UserData = o
		o.renderable.AddShape(o.Shape)
	}

	for _, a := range arts {
		o.AppendArt(a)
	}

	if !o.P.Static {
		Objects[o] = true
	}
}

func (o *Object) SetUserData(i interface{}) {
	if o.Shape == nil {
		log.Fatalln("shape not created, name:", o.Name)
	}
	o.Shape.UserData = i
}

func (o *Object) AddCallback(f Callback) {
	if o.Callbacks == nil {
		o.Callbacks = make(map[int]Callback)
	}
	o.Callbacks[len(o.Callbacks)] = f
}

func (o *Object) AddChild(child *Object) {
	if o.Childs == nil {
		o.Childs = make(map[*Object]bool)
	}
	o.Childs[child] = true
}

func (o *Object) Destroy() {
	if o.DestroyFunc != nil {
		o.DestroyFunc()
		return
	}

	o.needDestroy = true

	if o.Shape != nil {
		o.Shape.Body.Enabled = false
		phys.RemoveBody(o.Shape.Body)
		// space.RemoveShape(o.Shape) - crash need TODO
	}

	// Objects[o] = false

	for child := range o.Childs {
		child.Destroy()
	}

	// o.renderable.Destroy()
	// delete(Objects, o)
	// o = nil
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
