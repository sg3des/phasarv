package engine

import (
	"log"
	"point"
	"render"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/fizzle"

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

//NewObject create new object
func NewObject(Name string, PI *phys.Instruction, RI *render.Instruction, P *point.Param) *Object {

	o := &Object{
		PI: PI,
		RI: RI,
		P:  P,
	}

	if NeedRender && RI != nil {
		o.renderable = o.RI.Create(o.P)
	}

	if PI != nil {
		o.shape = o.PI.Create(o.P)
		o.shape.UserData = o

		if NeedRender && o.renderable != nil {
			o.renderable.AddShape(o.PI)
		}
	}

	if !o.P.Static {
		Objects.Add(o)
	}

	return o
}

//Create new object by instructons
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

func (o *Object) AddTrail(offset mgl32.Vec3, count int, size point.P, scale float32) *render.Particle {
	return o.renderable.NewTrail(offset, count, size, scale)
}

func (o *Object) SetDestroyFunc(f func()) {
	o.destroyFunc = f
}

func (o *Object) Destroy() {
	if o.destroyFunc != nil {
		o.destroyFunc()
		return
	}

	o.Remove()
}

func (o *Object) Remove() {
	if o.shape != nil {
		o.shape.Body.Enabled = false
		// phys.RemoveShape(o.shape)
		phys.RemoveBody(o.shape.Body)
		o.shape = nil
	}

	o.renderable.Destroy()
	o.needDestroy = true

	for child := range o.Childs {
		child.Remove()
	}
}

func (o *Object) Material() (*fizzle.Material, bool) {
	if o.renderable == nil {
		log.Printf("ERROR: %s not renderable", o.Name)
		return nil, false
	}

	if o.renderable.Body == nil {
		log.Printf("ERROR: %s does not have a body", o.Name)
		return nil, false
	}

	return o.renderable.Body.Material, true
}
