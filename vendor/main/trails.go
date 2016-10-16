package main

import (
	"engine"
	"param"
	"phys/vect"

	"github.com/go-gl/mathgl/mgl32"
)

func createTrail(p *engine.Object, piecelength float32, count int, offset mgl32.Vec2) {
	x, y := p.Position()
	t := &Trail{
		parent:    p,
		maxLength: piecelength,
		offset:    offset,
		x:         x,
		y:         y,
	}

	for i := 0; i < count; i++ {
		t.objects = append(t.objects, engine.NewPlane(param.Object{
			Name:        "trail",
			Material:    param.Material{Name: "laser", Texture: "laser", Shader: "colorblend", DiffColor: mgl32.Vec4{1, 1, 1, 1}},
			Transparent: true,
		}, 0.2, piecelength))
		t.points = append(t.points, trialPoints{})
	}

	engine.AddCallback(t.trailCallback)
}

//Trail from airplanes and rockets
type Trail struct {
	parent *engine.Object
	// prototype *engine.Object

	objects []*engine.Object
	points  []trialPoints

	maxLength float32
	offset    mgl32.Vec2

	x float32
	y float32
}

type trialPoints struct {
	X, Y, Angle, Alpha float32
}

func (p *trialPoints) Vect() vect.Vect {
	return vect.Vect{p.X, p.Y}
}

func (p *trialPoints) Vec2() mgl32.Vec2 {
	return mgl32.Vec2{p.X, p.Y}
}

func (t *Trail) trailCallback(dt float32) bool {

	//calculate alpha channel for trail pieces
	var sumAlpha float32
	for i, o := range t.objects {
		t.points[i].Alpha = t.points[i].Alpha - 1/float32(len(t.points)) - dt/2
		o.Node.Material.DiffuseColor[3] = t.points[i].Alpha
		if i == 0 {
			o.Node.Material.DiffuseColor[0] = 1
			o.Node.Material.DiffuseColor[1] = 0.3
			o.Node.Material.DiffuseColor[2] = 0
			o.Node.Scale = mgl32.Vec3{1.1, 2, 1}
		}
		if i == 1 {
			o.Node.Material.DiffuseColor[0] = 1
			o.Node.Material.DiffuseColor[1] = 0.6
			o.Node.Material.DiffuseColor[2] = 0.1
			o.Node.Scale = mgl32.Vec3{1.1, 1.5, 1}
		}
		sumAlpha += t.points[i].Alpha
	}

	//destroy trail if parent is nil
	if t.parent.Shape.Body == nil || t.parent == nil || t.parent.Node == nil || t.parent.Shape == nil {
		if sumAlpha <= 0.1 {
			for _, o := range t.objects {
				o.Destroy()
			}

			t = nil
			return false
		}
		return true
	}

	//calculate offset
	px, py := t.parent.Position()
	off := t.offset.Mul(vect.FAbs(t.parent.RollAngle / 2)).Sub(t.offset)
	vx, vy := t.parent.VectorSide(-off.X(), t.offset.Y())
	px += vx
	py += vy

	//if distance more then trail length, renew\shift trails
	dist := vect.Dist(vect.Vect{px, py}, t.points[0].Vect())
	if dist > t.maxLength {
		point := trialPoints{px, py, t.parent.Rotation(), 1}
		t.points = append([]trialPoints{point}, t.points[:len(t.points)]...)

		for i, o := range t.objects {
			o.SetPosition(t.points[i].X, t.points[i].Y)
			if i == 0 {
				o.SetRotation(t.points[i].Angle)
			} else {
				o.SetRotation(AngleObjectPoint(o, t.points[i-1].Vec2()))
			}
		}
	}

	return true
}
