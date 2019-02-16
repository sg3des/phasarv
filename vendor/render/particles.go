package render

import (
	"materials"
	"math"
	"phys/vect"
	"point"

	"github.com/chewxy/math32"
	"github.com/go-gl/mathgl/mgl32"
)

var Particles []*Particle

type calcFunc func(dt float32)

func drawParticles(dt float32) {
	for i := 0; i < len(Particles); i++ {
		p := Particles[i]

		if len(p.objects) == 0 {
			copy(Particles[i:], Particles[i+1:])
			Particles[len(Particles)-1] = nil // or the zero value of T
			Particles = Particles[:len(Particles)-1]
			i--
			continue
		}

		p.calculate(dt)
		p.render()
	}
}

type Particle struct {
	parent *Renderable
	Scale  float32

	i       int
	objects []*Renderable
	Offset  mgl32.Vec3

	dist            float32
	prevParentPoint mgl32.Vec3

	calculate calcFunc
}

func AngleBetweenPoints(a, b mgl32.Vec3) float32 {
	return float32(math.Atan2(float64(b.Y()-a.Y()), float64(b.X()-a.X())))
}

func (p *Particle) render() {
	for _, o := range p.objects {
		o.render()
	}
}

func (parent *Renderable) NewTrail(offset mgl32.Vec3, count int, size point.P, scale float32) *Particle {

	p := &Particle{
		parent: parent,
		Scale:  scale,
		Offset: offset,
		dist:   size.X * 0.9,
	}

	for i := 0; i < count; i++ {
		plane := &Renderable{
			Transparent: true,
			P:           &point.Param{Size: size},
			RI: &Instruction{
				MeshName:    "vector",
				Material:    &materials.Instruction{Name: "laser", Texture: "laser", Shader: "colorblend", DiffColor: mgl32.Vec4{1, 1, 1, 0.7}},
				Transparent: true,
			},
		}
		p.objects = append(p.objects, plane)
	}

	p.calculate = p.trailCalc

	Particles = append(Particles, p)
	parent.particles = append(Particles, p)

	return p
}

func (p *Particle) trailCalc(dt float32) {
	if p == nil || p.parent == nil || p.parent.Body == nil {
		return
	}

	// rot := p.parent.Body.LocalRotation
	// _, _, z := Euler(rot)

	// v2 := vect.FromAngle(z + p.Offset.Y())
	// v2.Mult(p.Offset.X())
	pos := p.parent.Body.Location
	pos = vect.FromVec3(pos).SubPoint(p.parent.Angle(), vect.FromVec3(p.Offset)).Vec3()

	// pos = pos.Add(v2.Vec3())

	angle := AngleBetweenPoints(pos, p.prevParentPoint)
	rot := mgl32.AnglesToQuat(0, 0, angle, 1)

	dist := pos.Sub(p.prevParentPoint).Len()

	if dist > p.dist {
		p.prevParentPoint = pos

		last := len(p.objects) - 1
		p.objects = append(p.objects[last:], p.objects[:last]...)

		p.objects[0].Body.Location = pos
		p.objects[0].Body.LocalRotation = rot
		p.objects[0].Body.Material.DiffuseColor[3] = 0.1
		p.objects[1].Body.Material.DiffuseColor[3] = 0.7
	}

	var summAlpha float32
	for _, o := range p.objects {
		o.Body.Material.DiffuseColor[3] -= (dt * 10) / float32(len(p.objects))
		if o.Body.Material.DiffuseColor[3] < 0 {
			o.Body.Material.DiffuseColor[3] = 0
		}
		summAlpha += o.Body.Material.DiffuseColor[3]
	}

	if summAlpha <= 1 && (p.parent == nil || p.parent.needDestroy) {
		p.objects = nil
	}
}

func (p *Particle) Destroy() {
	if p == nil {
		return
	}
	for _, o := range p.objects {
		o.needDestroy = true
	}
	// p.objects = nil
}

// Norm returns the L1-Norm of a Quaternion (W,X,Y,Z) -> Sqrt(W*W+X*X+Y*Y+Z*Z)
func Norm(q mgl32.Quat) float32 {
	x, y, z := q.V.Elem()
	return math32.Sqrt(q.W*q.W + x*x + y*y + z*z)
}

// Unit returns the Quaternion rescaled to unit-L1-norm
func Unit(q mgl32.Quat) mgl32.Quat {
	// x, y, z := q.V.Elem()
	k := Norm(q)
	q.W = q.W / k
	q.V[0] = q.V[0] / k
	q.V[1] = q.V[1] / k
	q.V[2] = q.V[2] / k
	return q
	// return mgl32.Quat{
	// 	W: q.W,
	// 	V: mgl32.Vec3[0]
	// }
	// return mgl32.Quat{q.W / k, x / k, y / k, z / k}
}

func Euler(q mgl32.Quat) (float32, float32, float32) {
	q = Unit(q)
	x, y, z := q.V.Elem()
	phi := math32.Atan2(2*(q.W*x+y*z), 1-2*(x*x+y*y))
	theta := math32.Asin(2 * (q.W*y - z*x))
	psi := math32.Atan2(2*(x*y+q.W*z), 1-2*(y*y+z*z))
	return phi, theta, psi
}

// func QuatAxis(q mgl32.Quat) (float32,float32,float32) {
// 	angle := 2 * math32.Acos(q.W)
// 	x = qx / sqrt(1-qw*qw)
// 	y = qy / sqrt(1-qw*qw)
// 	z = qz / sqrt(1-qw*qw)
// }
